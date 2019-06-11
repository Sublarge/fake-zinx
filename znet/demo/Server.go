package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Ping Handling ...")

	fmt.Println("recv from client: msgId = ", request.GetMsgId(), ",data: = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(0, []byte("ping... ping... ping... "))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloRouter struct {
	znet.BaseRouter
}

func (this *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Hello Handling ...")

	fmt.Println("recv from client: msgId = ", request.GetMsgId(), ",data: = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("hello... hello... hello... "))
	if err != nil {
		fmt.Println(err)
	}
}
func main() {
	server := znet.NewServer()

	server.AddRouter(0, &PingRouter{})
	server.AddRouter(1, &HelloRouter{})
	server.SetOnConnStart(func(connection ziface.IConnection) {
		fmt.Println("connection: ", connection.GetConnId(), " created!")
		connection.SetProperty("Name", "mld")
		property, _ := connection.GetProperty("Name")
		fmt.Println("name is : ", property)
	})

	server.SetOnConnStop(func(connection ziface.IConnection) {
		fmt.Println("connection: ", connection.GetConnId(), " destroyed!")
	})

	server.Serve()
}
