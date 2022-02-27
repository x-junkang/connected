package connect

import (
	"fmt"
	"net"

	"github.com/x-junkang/connected/internal/clog"
	"github.com/x-junkang/connected/internal/configure"
	"github.com/x-junkang/connected/pkg/ciface"
	"go.uber.org/zap"
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

	CID uint64
	// packet ziface.Packet
}

type Option func(s interface{})

//NewServer 创建一个服务器句柄
func NewServer(opts ...Option) *Server {
	s := &Server{
		Name:      configure.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        configure.GlobalObject.Host,
		Port:      configure.GlobalObject.TCPPort,
		// msgHandler: NewMsgHandle(),
		// ConnMgr:    NewConnManager(),
		// packet:     NewDataPack(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Start() {
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		clog.Logger.Fatal("addr is error", zap.String("errsmg", err.Error()))
	}
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		clog.Logger.Fatal("bind port fail", zap.String("errmsg", err.Error()))
	}
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			clog.Logger.Error("create new conn fail", zap.String("errmsg", err.Error()))
		}
		go s.handler(conn)
	}
}

func (s *Server) handler(tcpConn *net.TCPConn) {
	fmt.Println("hello client")
	conn := NewConnectionTcp(s, tcpConn, s.CID)
	conn.Start()
	s.CID++
}

func (s *Server) Stop() {
	s.ConnMgr.ClearConn()
	clog.Logger.Info("server stopped")
}

func (s *Server) AddRouter(msgID uint32, router ciface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
}

func (s *Server) GetConnMgr() ciface.IConnManager {
	return nil
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
	return nil
}