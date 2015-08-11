# Goder Line

Goder Line is a [Go](http://golang.org/) implementation of the chat software

Features:

* Public chat & topic chat
* Topic group created & joined
* Show prompting message when new user joined
* Show all users of every topic
* Login with exclusive name
* Manage list of topics by configuration file
* User can join any topics at the same time

### Documentation 
  
* [Guide for users](http://code.xiaojukeji.com/wiki/create/team6/User-guide-for-Team6)
* [Total design document](http://code.xiaojukeji.com/wiki/team6/Design-Document-for-Team6) 
* [Detail design document](http://code.xiaojukeji.com/wiki/create/team6/Detail-document-for-Team6)
* [Project schedule](http://code.xiaojukeji.com/wiki/create/team6/Project-schedule-for-Team6)

### Team Introduce

* Team Name   : GODER

* Team Slogan : Go run,Coders!

* Team Mmbers : Qiao Zhaohai(captain), Liu Boyu, Mo Yunfeng, He Wei

### Installation

####  PlatForm: Windows

* First, you should have download and install the Go compilers, if you not, you can refer to the [Course](http://jingyan.baidu.com/article/8cdccae965595c315413cda4.html)

* Second, you should have downlaod and install two third-party packages. 

	** one is a [websocket](code.google.com/p/go.net/websocket), You can get it from the [website](http://www.golangtc.com/download/package), then install the package under your Go environment variable folder, such as C:\Go\src\code.google.com.

	** other one is a [config](github.com/larspensjo/config), you should install the package under the Go environment variable folder, such as C:\Go\src\github.com\larspensjo\config.

* Third, download, build  all files

	git clone git@git.xiaojukeji.com:student/team6.git

	cd /team6/

	go build ./chat-server...

	chat-server.exe

* Fourth, run Goder Line
	
	Please input the ip address "localhost:7777" on your browser

####  PlatForm: Linux

* First, you should hava download and install the Go compilers, if you not, you can refer the [Course](http://www.cnblogs.com/huligong1234/p/golang.html)

* Second, you should have downlaod and install two third-party packages. 

	** one is a [websocket](code.google.com/p/go.net/websocket), You can get it from the [website](http://www.golangtc.com/download/package), then install the package under your Go environment variable folder, such as /usr/lib/golang/src/code.google.com.

	** other one is a [config](github.com/larspensjo/config), you should install the package under the Go environment variable folder, such as /usr/lib/golang/src/github.com/larspensjo/config.

* Third, download, build  all files and run Goder Line

	git clone git@git.xiaojukeji.com:student/team6.git

	cd /team6/chat-server

	go build 
	
	mv ./../config.ini ./

	./chat-server

* Fourth, run Goder Line
	
	Please input the ip address "localhost:7777" on your browser