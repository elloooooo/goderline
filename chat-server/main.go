package main

import (
	"fmt"
	"runtime"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("--------------------------------------")
			fmt.Println(err)
		}
	}()

	runtime.GOMAXPROCS(runtime.NumCPU())

	chatServer := &ChatServer{}
	chatServer.StartServer()
}
