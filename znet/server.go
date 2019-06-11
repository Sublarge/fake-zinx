package znet

import (
	"fmt"
	"net"
	"zinx/config"
	"zinx/ziface"
)

type Server struct {
	Name      string
	IpVersion string
	Ip        string
	Port      uint16
	MsgHandle ziface.IMsgHandle
	ConnMgr   ziface.IConnManager


	OnConnStart func(conn ziface.IConnection)
	OnConnStop  func(conn ziface.IConnection)
}

func (s *Server) Start() {
	fmt.Printf("[Server Name] %s (version: %s)", config.GlobalConfig.Name, config.GlobalConfig.Version)
	fmt.Printf("[Listen Addr] %s:%d\n", config.GlobalConfig.Host, config.GlobalConfig.TcpPort)

	go func() {
		s.MsgHandle.StartWorkPool()
		//	1 get a tcp address
		addr, err := net.ResolveTCPAddr(s.IpVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}
		// 2 listen the address
		listener, err := net.ListenTCP(s.IpVersion, addr)
		if err != nil {
			fmt.Println("listening error: ", err)
			return
		}
		fmt.Printf("[Listen] Server Listener at IP: %s, Port: %d Success!\n", s.Ip, s.Port)
		// 3. block for connecting
		var cid uint32 = 0
		for {
			conn, err := listener.AcceptTCP()
			if s.ConnMgr.Len() >= config.GlobalConfig.MaxConn {
				conn.Close()
				fmt.Println("Connection Pool is full!.........")
				continue
			}
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			dealConn := NewConnection(s, conn, cid, s.MsgHandle)
			cid++
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[Stop] Zinx server name: ", s.Name)
	s.ConnMgr.ClearConn()
}
func (s *Server) Serve() {
	s.Start()
	select {}
}
func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandle.AddRouter(msgId, router)
	fmt.Println("Add Router Success!")
}
func NewServer() ziface.IServer {
	s := &Server{
		Name:        config.GlobalConfig.Name,
		IpVersion:   "tcp4",
		Ip:          config.GlobalConfig.Host,
		Port:        config.GlobalConfig.TcpPort,
		MsgHandle:   NewMsgHandle(),
		ConnMgr:     NewConnManager(),
		OnConnStart: func(conn ziface.IConnection) {},
		OnConnStop:  func(conn ziface.IConnection) {},
	}
	return s
}
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

func (s *Server) SetOnConnStart(f func(connection ziface.IConnection)) {
	s.OnConnStart = f
}
func (s *Server) SetOnConnStop(f func(connection ziface.IConnection)) {
	s.OnConnStop = f
}
func (s *Server) CallOnConnStart(connection ziface.IConnection) {
	s.OnConnStart(connection)
}
func (s *Server) CallOnConnStop(connection ziface.IConnection) {
	s.OnConnStop(connection)
}
