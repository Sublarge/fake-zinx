package znet

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net"
	"sync"
	"zinx/config"
	"zinx/ziface"
)

type Connection struct {
	TcpServer ziface.IServer
	Conn      *net.TCPConn
	ConnID    uint32
	isClosed  bool
	ExitChan  chan bool

	msgChan chan []byte

	MsgHandle ziface.IMsgHandle

	extendProperty     map[string]interface{}
	extendPropertyLock sync.RWMutex
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, handle ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:      server,
		Conn:           conn,
		ConnID:         connID,
		isClosed:       false,
		ExitChan:       make(chan bool, 1),
		msgChan:        make(chan []byte),
		MsgHandle:      handle,
		extendProperty: make(map[string]interface{}),
	}
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

func (c *Connection) Start() {
	fmt.Println("Connection Start... ConnID = ", c.ConnID)
	go c.startReader()
	go c.startWriter()
	c.TcpServer.CallOnConnStart(c)
}
func (c *Connection) Stop() {
	fmt.Println("Connection Stop...", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true
	c.TcpServer.CallOnConnStop(c)

	c.Conn.Close()
	c.ExitChan <- true
	c.TcpServer.GetConnMgr().Remove(c)

	close(c.ExitChan)
	close(c.msgChan)
}
func (c *Connection) GetConnId() uint32 {
	return c.ConnID
}
func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed")
	}
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("pack msg error: ", msgId)
		return errors.New("pack message error")
	}
	c.msgChan <- binaryMsg

	return nil
}

func (c *Connection) startReader() {
	fmt.Println("Reader GoRoutine is running")
	defer fmt.Println("ConnID: ", c.ConnID, "\n[Reader is exited] \nremote address is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeaderLen())
		_, err := io.ReadFull(c.GetTcpConnection(), headData)
		if err != nil {
			fmt.Println("read msg head error", err)
			return
		}
		msg, err := dp.UnpackHeader(headData)
		if err != nil {
			fmt.Println("unpack head error", err)
			return
		}
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTcpConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
			}
		}
		msg.SetData(data)

		req := Request{
			c,
			msg,
		}
		if config.GlobalConfig.WorkerPoolSize > 0 {
			c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandle.DoMsgHandler(&req)
		}

	}
}
func (c *Connection) startWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error: ", err)
				continue
			}
		case <-c.ExitChan:
			return
		}
	}
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.extendPropertyLock.Lock()
	defer c.extendPropertyLock.Unlock()
	c.extendProperty[key] = value
}
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.extendPropertyLock.RLock()
	defer c.extendPropertyLock.RUnlock()
	value := c.extendProperty[key]
	if value == nil {
		return nil, errors.New("No Property [: " + key + " :]")
	}
	return value, nil
}
func (c *Connection) RemoveProperty(key string) {
	c.extendPropertyLock.Lock()
	defer c.extendPropertyLock.Unlock()
	delete(c.extendProperty, key)
}
