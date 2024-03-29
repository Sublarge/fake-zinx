package ziface

import (
	"net"
)

type IConnection interface {
	Start()
	Stop()
	GetConnId() uint32
	GetTcpConnection() *net.TCPConn
	RemoteAddr() net.Addr
	SendMsg(msgId uint32, data []byte) error

	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{},error)
	RemoveProperty(key string)
}

type HandleFunc func(*net.TCPConn, []byte, int) error
