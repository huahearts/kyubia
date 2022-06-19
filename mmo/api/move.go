package api

import (
	"fmt"

	"github.com/huahearts/kyubia/kiface"
	"github.com/huahearts/kyubia/knet"
	"github.com/huahearts/kyubia/mmo/core"
	"github.com/huahearts/kyubia/mmo/pb"
	"google.golang.org/protobuf/proto"
)

type MoveApi struct {
	knet.BaseRouter
}

func (*MoveApi) Callback(req kiface.IRequest) {
	msg := &pb.Position{}

	err := proto.Unmarshal(req.GetData(), msg)
	if err != nil {
		fmt.Println("Move: Position Unmarshal error ", err)
		return
	}

	pid, err := req.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid error", err)
		req.GetConnection().Stop()
		return
	}
	fmt.Printf("user pid = %d , move(%f,%f,%f,%f)", pid, msg.X, msg.Y, msg.Z, msg.V)
	//3. 根据pid得到player对象
	user := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//4. 让player对象发起移动位置信息广播
	user.UpdatePos(msg.X, msg.Y, msg.Z, msg.V)
}
