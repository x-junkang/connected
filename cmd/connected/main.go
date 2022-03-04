package main

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/x-junkang/connected/internal/clog"
	"github.com/x-junkang/connected/internal/config"
	"github.com/x-junkang/connected/internal/connect"
	"github.com/x-junkang/connected/internal/httphandler"
	"github.com/x-junkang/connected/pkg/ciface"
)

func init() {
	logConf := clog.Config{
		ConsoleLoggingEnabled: true,
		EncodeLogsAsJson:      true,
		FileLoggingEnabled:    true,
		Directory:             config.GlobalObject.LogDir,
		Filename:              config.GlobalObject.LogFile,
		MaxSize:               5,
		MaxBackups:            10,
		MaxAge:                7,
		Level:                 config.GlobalObject.LogLevel,
	}
	clog.Configure(logConf)
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
	log.Info().Msg("server start!")
	select {}
}

func main() {
	// server := connect.NewServer()
	// server.Start()
	service := NewService()
	service.Start()
}
