package kiface

type IMsgHandler interface {
	DoMsgHandler(req IRequest)
	AddRouter(msgID uint32, router IRouter)
	StartWorkerPool()
	SendMsgToTaskQueue(req IRequest)
}
