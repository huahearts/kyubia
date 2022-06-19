package api

import (
	"fmt"

	"github.com/huahearts/kyubia/kiface"
	"github.com/huahearts/kyubia/knet"
	"github.com/huahearts/kyubia/mmo/core"
	"github.com/huahearts/kyubia/mmo/pb"
	"google.golang.org/protobuf/proto"
)

type WorldChatAPI struct {
	knet.BaseRouter
}

func (*WorldChatAPI) Callback(req kiface.IRequest) {
	fmt.Println("TALK HANDLE")
	msg := &pb.Talk{}
	err := proto.Unmarshal(req.GetData(), msg)
	if err != nil {
		fmt.Println("talk Unmarshal error", err)
		return
	}

	pid, err := req.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid error", err)
		req.GetConnection().Stop()
		return
	}

	user := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	user.Talk(msg.Content)
}
