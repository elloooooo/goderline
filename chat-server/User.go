package main

import (
	"sync"
)

//type: 1 - 管理员， 2 - 普通用户
type User struct {
	UserId   int
	UserName string
	Type     int
}

var (
	//[sessionId:User]
	users     map[string]*User = make(map[string]*User, 10000)
	usersLock sync.RWMutex

	//[sessionId:connection]暂时没用到
	onlineUsers map[string]*connection = make(map[string]*connection, 10000)
	userLock    sync.RWMutex
)
