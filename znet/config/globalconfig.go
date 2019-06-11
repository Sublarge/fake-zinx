package config

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

type GlobalObj struct {
	TcpServer        ziface.IServer
	Name             string
	Host             string
	TcpPort          uint16
	Version          string
	MaxConn          int
	MaxPackageSize   uint32
	WorkerPoolSize   uint32
	MaxWorkerTaskLen uint32
}

func (g *GlobalObj) LoadConfig() {
	data, err := ioutil.ReadFile("demo/conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalConfig)
	if err != nil {
		panic(err)
	}
}

var GlobalConfig *GlobalObj

func init() {
	GlobalConfig = &GlobalObj{
		Host:             "0.0.0.0",
		Name:             "ZinxDemo",
		Version:          "V0.4",
		TcpPort:          8899,
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   16,
		MaxWorkerTaskLen: 1024,
	}

	GlobalConfig.LoadConfig()
}
