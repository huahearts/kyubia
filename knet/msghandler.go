package knet

import (
	"fmt"
	"strconv"

	"github.com/huahearts/kyubia/kiface"
)

type MsgHandler struct {
	Apis           map[uint32]kiface.IRouter
	WorkerPoolSize uint32
	TaskQueue      []chan kiface.IRequest
}

func NewMsgHandler() kiface.IMsgHandler { //需要配置
	return &MsgHandler{
		Apis:           make(map[uint32]kiface.IRouter),
		WorkerPoolSize: (uint32(4)),
		TaskQueue:      make([]chan kiface.IRequest, 4),
	}
}

func (mh *MsgHandler) SendMsgToTaskQueue(req kiface.IRequest) {
	//负载均衡
	workerID := req.GetConnection().GetConnID() % mh.WorkerPoolSize
	mh.TaskQueue[workerID] <- req
}

func (mh *MsgHandler) DoMsgHandler(req kiface.IRequest) {
	handle, ok := mh.Apis[req.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", req.GetMsgID(), "is unexist")
		return
	}

	handle.PreCallback(req)
	handle.Callback(req)
	handle.PostCallback(req)
}

func (mh *MsgHandler) AddRouter(msgId uint32, router kiface.IRouter) {
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api,msgID=" + strconv.Itoa(int(msgId)))
	}
	mh.Apis[msgId] = router
	fmt.Println("Api add msgID,", msgId)
}

func (mh *MsgHandler) StartOneWorker(workerID int, taskqueue chan kiface.IRequest) {
	fmt.Println("[Worker ID:]", workerID, "is started")
	for {
		select {
		case req := <-taskqueue:
			mh.DoMsgHandler(req)
		}
	}
}

func (mh *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan kiface.IRequest, 0) //配置
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}
