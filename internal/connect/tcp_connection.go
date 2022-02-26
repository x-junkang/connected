package connect

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/x-junkang/connected/internal/clog"
	"github.com/x-junkang/connected/internal/configure"
	"github.com/x-junkang/connected/internal/protocol"
	"go.uber.org/zap"
)

type ConnectionTCP struct {
	sync.RWMutex
	TCPServer *Server
	Conn      *net.TCPConn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnID uint32

	ExitChan chan struct{}

	msgBuffChan       chan []byte
	HeartbeatInterval time.Duration
	MaxBodySize       int

	property map[string]interface{}
	////保护当前property的锁
	propertyLock sync.Mutex
	//当前连接的关闭状态
	isClosed bool
}

func NewConnectionTcp(server *Server, conn *net.TCPConn, connID uint32) *ConnectionTCP {
	c := &ConnectionTCP{
		TCPServer: server,
		Conn:      conn,
		ConnID:    connID,
		isClosed:  false,
		ExitChan:  make(chan struct{}),
		// MsgHandler:  msgHandler,
		msgBuffChan: make(chan []byte, configure.GlobalObject.MaxMsgChanLen),
		MaxBodySize: 500,
		property:    nil,
	}
	return c
}

func (ct *ConnectionTCP) Start() {
	go ct.startReader()
	go ct.startWriter()
	//todo 将该连接加入全局map管理
	// ct.TCPServer.ConnMgr.Add(ct)
}

func (ct *ConnectionTCP) Stop() {
	if ct.isClosed {
		return
	}
	ct.isClosed = true
	close(ct.ExitChan)
}

func (ct *ConnectionTCP) startReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println(ct.Conn.RemoteAddr().String(), "[conn Reader exit!]")
	defer ct.Stop()
	for {
		select {
		case <-ct.ExitChan:
			return
		default:
		}
		header, err := ct.readHeader()
		if err != nil {
			clog.Error("read header fail", zap.String("err", err.Error()))
			return
		}
		data, err := ct.readBody(header.BodyLength)
		if err != nil {
			clog.Error("read body fail", zap.String("err", err.Error()))
			return
		}
		//todo handle data
		fmt.Println(string(data))
		ct.msgBuffChan <- data
	}
}

func (ct *ConnectionTCP) startWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(ct.Conn.RemoteAddr().String(), "[conn Writer exit!]")
	for {
		select {
		case data, ok := <-ct.msgBuffChan:
			if ok {
				//有数据要写给客户端
				header := &protocol.MarsHeader{
					BodyLength: int32(len(data)),
				}
				if _, err := ct.write(header, data); err != nil {
					fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				return
			}
		case <-ct.ExitChan:
			return
		}
	}
}

func (ct *ConnectionTCP) readHeader() (*protocol.MarsHeader, error) {
	var header protocol.MarsHeader

	// ct.Conn.SetReadDeadline(time.Now().Add(ct.HeartbeatInterval * 2))

	err := binary.Read(ct.Conn, binary.BigEndian, &header)
	if err != nil {
		return nil, err
	}

	if header.HeaderLength != protocol.MarsHeaderLength {
		return nil, errors.New("length is error")
	}

	bodyLen := int(header.BodyLength)
	if bodyLen > ct.MaxBodySize {
		return nil, errors.New("body is too large")
	}
	if bodyLen < 0 {
		return nil, errors.New("bodyLen must be larger than 0")
	}

	return &header, nil
}

func (ct *ConnectionTCP) readBody(len int32) ([]byte, error) {
	body := make([]byte, len)
	_, err := io.ReadFull(ct.Conn, body)
	if err != nil {
		clog.Error("read msg data fail", zap.String("err", err.Error()))
	}
	return body, err
}

func (ct *ConnectionTCP) write(header *protocol.MarsHeader, data []byte) (int, error) {
	err := binary.Write(ct.Conn, binary.BigEndian, header)
	if err != nil {
		return protocol.MarsHeaderLength, err
	}
	n, err := ct.Conn.Write(data)
	return n + protocol.MarsHeaderLength, err
}
