package main

import (
	"fmt"
	"github.com/huahearts/kyubia/kiface"
	"github.com/huahearts/kyubia/knet"
	"github.com/huahearts/kyubia/mmo/core"
)

func OnConnectionAdd(conn kiface.IConnection) {
	user := core.NewUser(conn)
	user.SyncPid()
	user.BroadCastStartPosition()
	fmt.Println("=====> Player pidId = ", user.Pid, " arrived ====")
}

func main() {
	s := knet.NewServer()
	s.SetOnConnStart(OnConnecionAdd)
	s.Serve()
}
