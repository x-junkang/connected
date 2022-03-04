package connect

import (
	"errors"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/x-junkang/connected/pkg/ciface"
)

//ConnManager 连接管理模块
type ConnManager struct {
	connections map[uint64]ciface.IConnection
	connLock    sync.RWMutex
}

//NewConnManager 创建一个链接管理
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint64]ciface.IConnection),
	}
}

//Add 添加链接
func (connMgr *ConnManager) Add(conn ciface.IConnection) {

	connMgr.connLock.Lock()
	//将conn连接添加到ConnMananger中
	connMgr.connections[conn.GetConnID()] = conn
	connMgr.connLock.Unlock()

	log.Info().Int("conns num", connMgr.Len()).Msg("connection add to ConnManager successfully")
}

//Remove 删除连接
func (connMgr *ConnManager) Remove(conn ciface.IConnection) {

	connMgr.connLock.Lock()
	//删除连接信息
	delete(connMgr.connections, conn.GetConnID())
	connMgr.connLock.Unlock()
	log.Info().Uint64("connID", conn.GetConnID()).Int("conns num", connMgr.Len()).Msg("remove connID successfully")
}

//Get 利用ConnID获取链接
func (connMgr *ConnManager) Get(connID uint64) (ciface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	}

	return nil, errors.New("connection not found")

}

//Len 获取当前连接
func (connMgr *ConnManager) Len() int {
	connMgr.connLock.RLock()
	length := len(connMgr.connections)
	connMgr.connLock.RUnlock()
	return length
}

//ClearConn 清除并停止所有连接
func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()

	//停止并删除全部的连接信息
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		//删除
		delete(connMgr.connections, connID)
	}
	connMgr.connLock.Unlock()
	log.Info().Int("conns num", connMgr.Len()).Msg("clear all connections successfully")
}

//ClearOneConn  利用ConnID获取一个链接 并且删除
func (connMgr *ConnManager) ClearOneConn(connID uint64) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	connections := connMgr.connections
	if conn, ok := connections[connID]; ok {
		//停止
		conn.Stop()
		//删除
		delete(connections, connID)
		log.Info().Uint64("connID", connID).Msg("clear connection succeed")
		return
	}
	log.Error().Uint64("connID", connID).Msg("clear connection fail")
}
