package main

import (
	"fmt"
	"net"

	"github.com/x-junkang/connected/internal/clog"
	"go.uber.org/zap"
)

func init() {
	clog.InitLogger("./connect.log", "debug")
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:8090")
	if err != nil {
		clog.Fatal("bind port fail", zap.String("err", err.Error()))
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			clog.Error("create new conn fail", zap.String("err", err.Error()))
		}
		go handler(conn)
	}
}

func handler(conn net.Conn) {
	fmt.Println("hello client")
}
