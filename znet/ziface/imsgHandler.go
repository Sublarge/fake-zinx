package ziface

type IMsgHandle interface {
	DoMsgHandler(request IRequest)
	AddRouter(msgId uint32, router IRouter)
	StartWorkPool()
	SendMsgToTaskQueue(request IRequest)
}
