package main

import (
	"fmt"
	"net/http"

	"github.com/x-junkang/connected/internal/clog"
	"github.com/x-junkang/connected/internal/connect"
	"github.com/x-junkang/connected/internal/httphandler"
	"github.com/x-junkang/connected/pkg/ciface"
)

func init() {
	clog.InitLogger("./connect.log", "debug")
}

type Service struct {
	server     ciface.IServer
	httpServer *http.ServeMux
}

func NewService() *Service {
	return &Service{
		server:     connect.NewServer(),
		httpServer: &http.ServeMux{},
	}
}

func (s *Service) Start() {
	s.server.AddRouter(1, &connect.HelloRouter{})
	s.server.AddRouter(2, &connect.SendMsgRouter{})
	go func() {
		s.server.Start()
	}()
	go func() {
		http.ListenAndServe(":8080", httphandler.NewHttpServer(s.server))
	}()
	fmt.Println("server start!")
	select {}
}

func main() {
	// server := connect.NewServer()
	// server.Start()
	service := NewService()
	service.Start()
}
