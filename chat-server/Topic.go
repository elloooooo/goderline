package main

import (
	"fmt"
	"sync"
)

var topicIdGen int = 0
var topicIdGenLock sync.Mutex

type hub struct {
	topicId int

	tag string

	// Registered connections.[con:sessionId]
	connections map[*connection]string

	// Inbound messages from the connections.
	broadcast chan p2pMsg

	// P2P send message
	p2pSendMsg chan p2pMsg

	//fetch data message
	/*	fetchDataMsg chan p2pMsg*/

	// Register requests to topic.
	register chan *ConNode

	// Unregister requests from topic.
	unregister chan *connection

	// close connection
	close chan *connection
}

type ConNode struct {
	userId string
	con    *connection
}

type p2pMsg struct {
	sessionId string
	dstId     int
	msg       []byte
}

func newHub() *hub {
	return &hub{
		topicId:-1,
		tag: "Universe",
		broadcast:  make(chan p2pMsg),
		p2pSendMsg: make(chan p2pMsg),
		/*fetchDataMsg: make(chan p2pMsg),*/
		register:    make(chan *ConNode),
		unregister:  make(chan *connection),
		connections: make(map[*connection]string),
		close:       make(chan *connection),
	}
}

func (h *hub) run() {
	for {
		/*fmt.Println(h.tag, " ", len(h.connections))*/
		select {
		case c := <-h.register:
			h.connections[c.con] = c.userId
			c.con.syncChan <- 1
			fmt.Println("Topic register: ", h.connections[c.con])
		case c := <-h.unregister:
			for con, v := range h.connections {
				if con == c {
					delete(h.connections, c)
					fmt.Println("Topic Unregister: ", v)
				}
			}
			c.syncChan <- 1
		case c := <-h.close:
			for con, v := range h.connections {
				if con == c {
					delete(h.connections, c)
					fmt.Println("Topic Unregister: ", v)
				}
			}
		case m := <-h.broadcast:
			for c, id := range h.connections {
				if m.dstId == -1 || (id != m.sessionId && c != nil) {
					select {
					case c.send <- m.msg:
					default:
						delete(h.connections, c)
						close(c.send)
					}
				}
			}
		case p2pmsg := <-h.p2pSendMsg:
			for c, id := range h.connections {
				if id == p2pmsg.sessionId {
					select {
					case c.send <- p2pmsg.msg:
					default:
						delete(h.connections, c)
						close(c.send)
					}
				}
			}
		}
	}
}

func generateTopicId() int {
	topicIdGenLock.Lock()
	curTopicId := topicIdGen
	topicIdGen++
	topicIdGenLock.Unlock()

	return curTopicId
}
