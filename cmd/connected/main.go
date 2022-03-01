package main

import (
	"github.com/x-junkang/connected/internal/clog"
	"github.com/x-junkang/connected/internal/connect"
)

func init() {
	clog.InitLogger("./connect.log", "debug")
}

func main() {
	server := connect.NewServer()
	server.AddRouter(1, &connect.HelloRouter{})
	server.AddRouter(2, &connect.SendMsgRouter{})
	server.Start()
}
