package connect

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/x-junkang/connected/internal/config"
	"github.com/x-junkang/connected/internal/protocol"
	"github.com/x-junkang/connected/pkg/ciface"
)

type ConnectionTCP struct {
	sync.RWMutex
	TCPServer ciface.IServer
	Conn      *net.TCPConn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnID uint64

	MsgHandler ciface.IMsgHandle
	ExitChan   chan struct{}

	msgBuffChan       chan []byte
	HeartbeatInterval time.Duration
	MaxBodySize       int

	property map[string]interface{}
	////保护当前property的锁
	propertyLock sync.Mutex
	//当前连接的关闭状态
	isClosed bool
}

func NewConnectionTcp(server ciface.IServer, conn *net.TCPConn, connID uint64, msgHandler ciface.IMsgHandle) *ConnectionTCP {
	c := &ConnectionTCP{
		TCPServer:         server,
		Conn:              conn,
		ConnID:            connID,
		isClosed:          false,
		ExitChan:          make(chan struct{}),
		MsgHandler:        msgHandler,
		msgBuffChan:       make(chan []byte, config.GlobalObject.MaxMsgChanLen),
		MaxBodySize:       500,
		HeartbeatInterval: 3 * time.Second,
		property:          nil,
	}
	return c
}

func (ct *ConnectionTCP) Start() {
	go ct.startReader()
	go ct.startWriter()
	//todo 将该连接加入全局map管理
	ct.TCPServer.GetConnMgr().Add(ct)
}

func (ct *ConnectionTCP) Stop() {
	if ct.isClosed {
		return
	}
	ct.isClosed = true
	ct.TCPServer.GetConnMgr().Remove(ct)
	close(ct.ExitChan)
}

func (ct *ConnectionTCP) GetTCPConnection() *net.TCPConn {
	return ct.Conn
}

func (ct *ConnectionTCP) GetConnID() uint64 {
	return ct.ConnID
}
func (ct *ConnectionTCP) RemoteAddr() net.Addr {
	return ct.Conn.RemoteAddr()
}
func (ct *ConnectionTCP) SendMsg(msgID uint32, data []byte) error {
	// 需要完善
	ct.msgBuffChan <- data
	return nil
}
func (ct *ConnectionTCP) SendBuffMsg(msgID uint32, data []byte) error {
	return nil
}
func (ct *ConnectionTCP) SetProperty(key string, value interface{}) {
	ct.propertyLock.Lock()
	defer ct.propertyLock.Unlock()
	ct.property[key] = value
}
func (ct *ConnectionTCP) GetProperty(key string) (interface{}, error) {
	ct.propertyLock.Lock()
	defer ct.propertyLock.Unlock()
	if v, ok := ct.property[key]; ok {
		return v, nil
	}
	return nil, errors.New("key does not exit")
}

func (ct *ConnectionTCP) RemoveProperty(key string) {
	ct.propertyLock.Lock()
	defer ct.propertyLock.Unlock()
	delete(ct.property, key)
}

func (ct *ConnectionTCP) startReader() {
	log.Info().Msg("[reader goroutine is running]")
	defer log.Info().Uint64("connID", ct.ConnID).Msg("[conn reader exit!]")
	defer ct.Stop()

	var err error

	for {
		select {
		case <-ct.ExitChan:
			return
		default:
		}
		msg := protocol.NewMarsMsg()
		msg.MarsHeader, err = ct.readHeader()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Err(err).Msg("read header fail")
			return
		}
		if msg.GetHeaderLen() > 20 {
			opt, err := ct.readBytes(msg.GetHeaderLen() - 20)
			if err != nil {
				return
			}
			msg.Opt = opt
		}
		data, err := ct.readBytes(msg.GetDataLen())
		if err != nil {
			log.Err(err).Msg("read body fail")
			return
		}
		msg.SetData(data)
		//todo handle data
		log.Info().Msgf("msg = %s", string(msg.Data))

		req := &Request{
			conn: ct,
			msg:  msg,
		}
		ct.MsgHandler.DoMsgHandler(req)
		// ct.msgBuffChan <- data

	}
}

func (ct *ConnectionTCP) startWriter() {
	log.Info().Msg("[writer goroutine is running]")
	defer log.Info().Uint64("connID", ct.ConnID).Msg("[conn writer exit!]")
	for {
		select {
		case data, ok := <-ct.msgBuffChan:
			if ok {
				//有数据要写给客户端
				header := &protocol.MarsHeader{
					HeaderLength: 20,
					BodyLength:   uint32(len(data)),
				}
				if _, err := ct.write(header, data); err != nil {
					log.Warn().Err(err).Msg("conn writer exit")
					return
				}
			} else {
				log.Info().Msg("msgBuffChan is closed")
				return
			}
		case <-ct.ExitChan:
			return
		}
	}
}

func (ct *ConnectionTCP) readHeader() (*protocol.MarsHeader, error) {
	var header protocol.MarsHeader

	ct.Conn.SetReadDeadline(time.Now().Add(ct.HeartbeatInterval * 2))

	err := binary.Read(ct.Conn, binary.LittleEndian, &header)
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

func (ct *ConnectionTCP) readBytes(len uint32) ([]byte, error) {
	body := make([]byte, len)
	_, err := io.ReadFull(ct.Conn, body)
	if err != nil {
		log.Err(err).Msg("read msg data fail")
	}
	return body, err
}

func (ct *ConnectionTCP) write(header *protocol.MarsHeader, data []byte) (int, error) {
	err := binary.Write(ct.Conn, binary.LittleEndian, header)
	if err != nil {
		return protocol.MarsHeaderLength, err
	}
	n, err := ct.Conn.Write(data)
	return n + protocol.MarsHeaderLength, err
}
