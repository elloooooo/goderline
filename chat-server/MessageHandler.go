package main

import (
	"encoding/json"
	"fmt"
	"time"
)

//message handler, do massage parser business
type MessageHandler struct {
	con *connection
}

//user list group by topic
type topicInfo struct {
	Id       int        `json:"id"`
	Text     string     `json:"text"`
	Type     int        `json:"type"`
	Children []userInfo `json:"children"`
}

type userInfo struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
	Type int    `json:"type"`
}

func (this *MessageHandler) handleMsg(msg Request) {
	var resp Response
	resp.Command = msg.Command
	resp.Time = time.Now().Format("2006-01-02 15:04:05")
	if msg.Command == "CONNECT" {
		isExisted := false

		//save user
		usersLock.Lock()
		for _, v := range users {
			if v.UserName == msg.SrcName {
				isExisted = true
				break
			}
		}

		if isExisted {
			this.con.ws.Close()
			close(this.con.send)
			this.con.isClose = true
			wholeHub.close <- this.con
		} else {
			onlineUsers[msg.SessionId] = this.con
			this.con.sessionId = msg.SessionId
			this.con.userId = msg.SrcId
			users[msg.SessionId] = &User{UserId: msg.SrcId, UserName: msg.SrcName}

			wholeHub.connections[this.con] = msg.SessionId

			resp.Status = "OK"
			b := encodeMsg(resp)
			tmpMsg := p2pMsg{sessionId: msg.SessionId, msg: b}
			wholeHub.p2pSendMsg <- tmpMsg
		}
		usersLock.Unlock()
	} else if msg.Command == "CHANGENAME" {
		usersLock.Lock()
		user := users[msg.SessionId]
		user.UserName = msg.SrcName
		users[msg.SessionId] = user
		usersLock.Unlock()

		//refresh user list
		notifyRefreshUserList(resp, msg)
	} else if msg.Command == "ADDTOPIC" {
		topicsLock.Lock()
		topicId := generateTopicId()
		topics[topicId] = newHub()
		topics[topicId].topicId = topicId
		topics[topicId].tag = msg.DstName
		topicsLock.Unlock()
		resp.DstId = topicId
		resp.DstName = msg.DstName

		//refresh user list
		notifyRefreshUserList(resp, msg)
	} else if msg.Command == "LISTTOPICS" {
		notifyRefreshUserList(resp, msg)
	} else if msg.Command == "JOINTOPIC" {
		tmpTopic := joinTopic(msg.DstId, msg.SessionId, this.con)

		c := <-this.con.syncChan
		if c == 1 {
			//refresh user list
			notifyRefreshUserList(resp, msg)
		}

		sendNotifyMsg(resp,msg,tmpTopic,fmt.Sprintf("User [%s] Join in.", msg.SrcName))
	} else if msg.Command == "QUITTOPIC" {
		tmpTopic := quitTopic(msg.DstName, msg.SessionId, this.con)

		if tmpTopic != nil {
			c := <-this.con.syncChan
			if c == 1 {
				//refresh user list
				notifyRefreshUserList(resp, msg)
			}

			sendNotifyMsg(resp,msg,tmpTopic,fmt.Sprintf("User [%s] Quit.", msg.SrcName))
		}
	} else if msg.Command == "MSG" {
		resp.Msg = msg.Msg
		resp.Type = msg.Type
		resp.DstId = msg.DstId
		resp.DstName = msg.DstName

		//broadcast msg
		if msg.Type == 1 {
			resp.SrcId = msg.SrcId
			resp.SrcName = msg.SrcName
			b := encodeMsg(resp)

			tmpMsg := p2pMsg{sessionId: msg.SessionId, msg: b}

			// Write send a message to the client.
			if msg.DstId == -1 {
				wholeHub.broadcast <- tmpMsg
			} else {
				topicsLock.Lock()
				if topics[msg.DstId] != nil {
					topics[msg.DstId].broadcast <- tmpMsg
				} else {
					fmt.Printf("Topic gourp id [%d] not existed.", msg.DstId)
				}
				topicsLock.Unlock()
			}
		} else if msg.Type == 2 {
			//send to a single person
			resp.SrcId = msg.SrcId
			resp.SrcName = msg.SrcName

			var dstSessionId string
			usersLock.Lock()
			for k, v := range users {
				if v.UserId == msg.DstId {
					dstSessionId = k
					break
				}
			}
			usersLock.Unlock()

			b := encodeMsg(resp)

			tmpMsg := p2pMsg{sessionId: dstSessionId, msg: b}
			wholeHub.p2pSendMsg <- tmpMsg
		}
	}
}

func sendNotifyMsg(resp Response, msg Request, topic *hub, message string){
	resp.Command = "MSG"
	resp.SrcId = -1
	resp.SrcName = "系统提示："
	resp.Type = 1
	resp.DstId = msg.DstId
	resp.DstName = msg.DstName
	resp.Msg = message
	resp.Status = "OK"
	b := encodeMsg(resp)
	tmpMsg := p2pMsg{sessionId: msg.SessionId, msg: b, dstId: -1}

	topic.broadcast <- tmpMsg
}

//refresh user list
func notifyRefreshUserList(resp Response, msg Request) {
	resp.Command = "LISTTOPICS"
	resp.SrcId = msg.SrcId
	resp.Msg = string(listTopics())
	resp.Status = "OK"
	b := encodeMsg(resp)
	tmpMsg := p2pMsg{sessionId: msg.SessionId, msg: b, dstId: -1}

	wholeHub.broadcast <- tmpMsg
}

//quit topic
func quitTopic(topicName string, sessionId string, con *connection) *hub{
	topicsLock.Lock()
	isTopicExist := false
	var tmpTopic *hub = nil

	for _, topic := range topics {
		if topic.tag == topicName {
			topic.unregister <- con
			tmpTopic = topic
			isTopicExist = true
		}
	}

	if !isTopicExist {
		fmt.Printf("Topic gourp [%s] not existed.", topicName)
	}
	topicsLock.Unlock()

	return tmpTopic
}

//join topic
func joinTopic(topicId int, sessionId string, con *connection) *hub{
	topicsLock.Lock()
	var tmpTopic *hub
	if topics[topicId] != nil {
		tmpTopic = topics[topicId]

		conNode := &ConNode{userId: sessionId, con: con}
		tmpTopic.register <- conNode

	} else {
		fmt.Printf("Topic gourp id [%d] not existed.", topicId)
	}
	topicsLock.Unlock()

	return tmpTopic
}

//list topics
func listTopics() []byte {
	var topicInfos []topicInfo
	fmt.Println("ListTopic")

	topicsLock.Lock()
	usersLock.Lock()
	//Universe
	topicinfo := topicInfo{Id: wholeHub.topicId, Text: wholeHub.tag}

	var userInfos []userInfo
	for con, _ := range wholeHub.connections {
		userinfo := userInfo{Id: con.userId, Text: users[con.sessionId].UserName}
		userinfo.Type = 2
		userInfos = append(userInfos, userinfo)
	}
	topicinfo.Type = 1
	topicinfo.Children = userInfos
	topicInfos = append(topicInfos, topicinfo)

	fmt.Println(topicinfo)

	for topicId, topic := range topics {
		topicinfo := topicInfo{Id: topicId, Text: topic.tag}

		var userInfos []userInfo
		for con, _ := range topic.connections {
			userinfo := userInfo{Id: con.userId, Text: users[con.sessionId].UserName}
			userinfo.Type = 2
			userInfos = append(userInfos, userinfo)
		}
		topicinfo.Type = 1
		topicinfo.Children = userInfos
		topicInfos = append(topicInfos, topicinfo)
	}
	usersLock.Unlock()
	topicsLock.Unlock()
	/*fmt.Println(topicInfos)*/
	b, err := json.Marshal(topicInfos)
	if err != nil {
		fmt.Println(err.Error())
	}
	/*fmt.Printf("%q", b)*/
	if err != nil {
		return nil
	}

	return b
}
