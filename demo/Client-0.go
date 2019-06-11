package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

func main() {
	fmt.Printf("client-0 start ... ")
	time.Sleep(1 * time.Second)

	//	1. connect to server
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err ,exit!")
		return
	}
	//  2. send data
	for {
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMessage(0, []byte("ping server")))
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = conn.Write(binaryMsg)
		if err != nil {
			fmt.Println(err)
			return
		}

		binaryHead := make([]byte, dp.GetHeaderLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error: ", err)
		}
		msg, err := dp.UnpackHeader(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHead error: ", err)
			break
		}
		if msg.GetMsgLen() > 0 {
			data := make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, data); err != nil {
				fmt.Println("read body error: ", err)
				return
			}
			msg.SetData(data)
			fmt.Println("recv server msg : ")
			fmt.Println("msg id   : ", msg.GetMsgId())
			fmt.Println("msg len  : ", msg.GetMsgLen())
			fmt.Println("msg body : ", string(msg.GetData()))
			fmt.Println()
		}

		time.Sleep(1 * time.Second)
	}

}
