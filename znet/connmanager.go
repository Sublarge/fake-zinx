package znet

import (
	"fmt"
	"github.com/kataras/iris/core/errors"
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (self *ConnManager) Add(conn ziface.IConnection) {
	self.connLock.Lock()
	defer self.connLock.Unlock()

	self.connections[conn.GetConnId()] = conn
	fmt.Println("connection add to ConnManager successfully: conn num = ", self.Len())
}
func (self *ConnManager) Remove(conn ziface.IConnection) {
	self.connLock.Lock()
	defer self.connLock.Unlock()

	delete(self.connections, conn.GetConnId())
	fmt.Println(conn.GetConnId(), "connection remove to ConnManager successfully: conn num = ", self.Len())
}
func (self *ConnManager) Get(connId uint32) (ziface.IConnection, error) {
	self.connLock.RLock()
	defer self.connLock.RUnlock()

	connection := self.connections[connId]
	if connection == nil {
		return nil, errors.New("connection not found!")
	} else {
		return connection, nil
	}
}
func (self *ConnManager) Len() int {
	return len(self.connections)
}
func (self *ConnManager) ClearConn() {
	self.connLock.Lock()
	defer self.connLock.Unlock()

	for connId, conn := range self.connections {
		conn.Stop()
		delete(self.connections, connId)
	}
	fmt.Println("clear all connection successfully!")
}
