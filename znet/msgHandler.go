package znet

import (
	"fmt"
	"zinx/config"
	"zinx/ziface"
)

type MsgHandle struct {
	Apis           map[uint32]ziface.IRouter
	TaskQueue      [] chan ziface.IRequest
	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),

		WorkerPoolSize: config.GlobalConfig.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, config.GlobalConfig.WorkerPoolSize),
	}
}

func (self *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	router := self.Apis[request.GetMsgId()]
	if router == nil {
		fmt.Println("api msgId ,", request.GetMsgId(), " is not found")
		return
	}
	router.PreHandle(request)
	router.Handle(request)
	router.Handle(request)
}
func (self *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	self.Apis[msgId] = router
}
func (self *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	workerId := request.GetConnection().GetConnId() % self.WorkerPoolSize
	fmt.Println("add connId: ", request.GetConnection().GetConnId(),
		"request MsgId: ", request.GetMsgId(),
		"to WorkerId: ", workerId)
	self.TaskQueue[workerId] <- request
}

func (self *MsgHandle) StartWorkPool() {
	for i := 0; i < int(self.WorkerPoolSize); i++ {
		self.TaskQueue[i] = make(chan ziface.IRequest, config.GlobalConfig.MaxWorkerTaskLen)
		go self.startOneWorker(i, self.TaskQueue[i])
	}
}
func (self *MsgHandle) startOneWorker(workerId int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker Id : ", workerId, " is started ...")
	for {
		select {
		case request := <-taskQueue:
			self.DoMsgHandler(request)
		}
	}
}
