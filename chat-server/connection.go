package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
)

type connection struct {
	sessionId string

	userId int

	//The websocket connection
	ws *websocket.Conn

	//Buffered channel for send messages
	send chan []byte

	// topic, not used yet
	/*	h *hub*/

	//If connection closed
	isClose bool

	//sync channel
	syncChan chan int
}

func decodeMsg(buf []byte) Request {
	var msg Request
	err := json.Unmarshal(buf, &msg)

	if err != nil {
		fmt.Println(err)
	}

	return msg
}

func encodeMsg(rsp Response) []byte {
	b, err := json.Marshal(rsp)

	if err != nil {
		fmt.Println(err)
	}

	return b
}

func (c *connection) reader() {
	/*panic("test error")*/
	for {
		// Read at most 1024 bytes.
		buf := make([]byte, 1024)
		n, err := c.ws.Read(buf)
		if err != nil {
			/*			m := fmt.Sprintf("Client [%s] logout.", c.sessionId)
						panic(m)*/
			fmt.Println(err)
			break
		}
		fmt.Printf("recv:%q\n", buf[:n])

		msg := decodeMsg(buf[:n])

		msgHandle := &MessageHandler{con: c}
		msgHandle.handleMsg(msg)
	}
	defer func() {
		c.ws.Close()
		close(c.send)
		c.isClose = true
		wholeHub.close <- c
		for _, topic := range topics {
			topic.close <- c
		}

		usersLock.Lock()
		delete(users, c.sessionId)
		usersLock.Unlock()
	}()
}

func (c *connection) writer() {
	for message := range c.send {
		fmt.Printf("send:%q\n", message)
		_, err := c.ws.Write(message)
		if err != nil {
			/*	m := fmt.Sprintf("Client [%s] logout.", c.sessionId)
				panic(m)*/
			fmt.Println(err)
			break
		}
	}
	c.ws.Close()
}
