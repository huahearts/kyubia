package main

import (
	"fmt"

	"github.com/huahearts/kyubia/kiface"
	"github.com/huahearts/kyubia/knet"
	"github.com/huahearts/kyubia/mmo/api"
	"github.com/huahearts/kyubia/mmo/core"
)

func OnConnectionAdd(conn kiface.IConnection) {
	user := core.NewUser(conn)
	user.SyncPid()
	user.BroadCastStartPosition()
	core.WorldMgrObj.AddUser(user)
	conn.SetProperty("pid", user.Pid)
	user.SyncSurrounding()
	fmt.Println("=====> Player pidId = ", user.Pid, " arrived ====")

}

func main() {
	s := knet.NewServer()
	s.SetOnConnStart(OnConnectionAdd)
	s.AddRouter(2, &api.WorldChatAPI{})
	s.AddRouter(3, &api.MoveApi{}) //移动
	s.Serve()
}
