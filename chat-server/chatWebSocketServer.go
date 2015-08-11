package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"runtime"
)

type ChatServer struct {
}

//web socket message
// type: 1 - 广播， 2 - 点对点
// DstId: type为1时，存topicId；type为2时，存recId
type Request struct {
	SessionId string
	Command   string
	Type      int
	SrcId     int
	SrcName   string
	DstId     int
	DstName   string
	Msg       string
}

// type:
// type: 1 - 广播， 2 - 点对点
// SrcId: type为1，存topicId；type为2时，存recId
type Response struct {
	Command   string
	Type      int
	Msg       string
	DstId   int
	DstName string
	SrcId     int
	SrcName   string
	Time string
	Status    string
}

var (
	port *int = flag.Int("p", 7777, "Port to listen.")
	//config.ini
	configFile = flag.String("cfg", "config.ini", "General configuration file")

	topics     map[int]*hub = make(map[int]*hub, 100)
	topicsLock sync.RWMutex

	wholeHub *hub = newHub()
	
	clientPath string
)

func initwsConnection(ws *websocket.Conn) {
	//Here we are creating list of clients that gets connected
	defer func() {
		defer ws.Close()
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	socketClientIP := ws.Request().RemoteAddr
	fmt.Printf("Connected remote IP: %s", socketClientIP)

	con := &connection{send: make(chan []byte, 256), ws: ws, syncChan: make(chan int)}
	conNode := &ConNode{con: con}
	wholeHub.register <- conNode
	/*panic(45)*/

	c := <-con.syncChan
	if c == 1 {
		go con.writer()
		con.reader()
	}
}

//init single http handle
func initHttpHandle(pattern string, handler websocket.Handler) {
	http.HandleFunc(pattern,
		func(w http.ResponseWriter, req *http.Request) {
			s := websocket.Server{Handler: websocket.Handler(handler)}
			s.ServeHTTP(w, req)
		})
}

func initHttpHandles() {
	/*http.HandleFunc("/listTopics", listTopics)*/
	initHttpHandle("/chat", initwsConnection)
	http.Handle("/", http.FileServer(http.Dir(clientPath)))
}

func initTopics() {
	for k, tName := range TOPIC {
		id, err := strconv.Atoi(k)
		if err != nil {
			fmt.Printf("Error topic define: [%s] - [%s]", k, tName)
			continue
		}
		topic := newHub()
		topic.topicId = id
		topic.tag = tName

		topics[id] = topic
	}

	for _, v := range topics {
		go v.run()
	}
}

func (this *ChatServer) StartServer() {
	flag.Parse()
	configHelper := &ConfigHelper{}
	configHelper.LoadConfig(configFile)
	
	if runtime.GOOS == WINDOWS{
		clientPath = "web-client"
	}else if runtime.GOOS == LINUX{
		clientPath = "../web-client"
	}

	go wholeHub.run()
	initTopics()
	initHttpHandles()

	fmt.Printf("Server started.Listen at port:%d/\n", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		panic("ChatServer ListenANdServe: " + err.Error())
	}
}
