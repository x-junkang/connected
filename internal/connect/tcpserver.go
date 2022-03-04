package connect

import (
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"github.com/x-junkang/connected/internal/config"
	"github.com/x-junkang/connected/internal/protocol"
	"github.com/x-junkang/connected/pkg/ciface"
)

type Server struct {
	//服务器的名称
	Name string
	//tcp4 or other
	IPVersion string
	//服务绑定的IP地址
	IP string
	//服务绑定的端口
	Port int
	//当前Server的消息管理模块，用来绑定MsgID和对应的处理方法
	msgHandler ciface.IMsgHandle
	//当前Server的链接管理器
	ConnMgr ciface.IConnManager
	//该Server的连接创建时Hook函数
	OnConnStart func(conn ciface.IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn ciface.IConnection)

	CID    uint64
	packet ciface.Packet
}

type Option func(s interface{})

//NewServer 创建一个服务器句柄
func NewServer(opts ...Option) *Server {
	s := &Server{
		Name:       config.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         config.GlobalObject.Host,
		Port:       config.GlobalObject.TCPPort,
		msgHandler: NewTcpHandler(),
		ConnMgr:    NewConnManager(),
		packet:     protocol.NewDataPack(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Start() {

	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		log.Err(err).Msg("addr is error")
	}
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		log.Err(err).Msg("bind port fail")
	}
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Err(err).Msg("create new conn fail")
		}
		go s.handler(conn)
	}
}

func (s *Server) handler(tcpConn *net.TCPConn) {
	fmt.Println("hello client")
	conn := NewConnectionTcp(s, tcpConn, s.CID, s.msgHandler)
	conn.Start()
	s.CID++
}

func (s *Server) Stop() {
	s.ConnMgr.ClearConn()
	log.Info().Str("name", s.Name).Msg("server stopped")
}

func (s *Server) AddRouter(msgID uint32, router ciface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
}

func (s *Server) GetConnMgr() ciface.IConnManager {
	return s.ConnMgr
}

func (s *Server) SetOnConnStart(fn func(ciface.IConnection)) {

}
func (s *Server) SetOnConnStop(fn func(ciface.IConnection)) {

}

func (s *Server) CallOnConnStart(conn ciface.IConnection) {

}
func (s *Server) CallOnConnStop(conn ciface.IConnection) {

}
func (s *Server) Packet() ciface.Packet {
	return s.packet
}
