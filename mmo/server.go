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
	fmt.Println("=====> Player pidId = ", user.Pid, " arrived ====")
	conn.SetProperty("pid", user.Pid)

}

func main() {
	s := knet.NewServer()
	s.SetOnConnStart(OnConnectionAdd)

	s.AddRouter(2, &api.WorldChatAPI{})
	s.Serve()
}
