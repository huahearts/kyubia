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

//当客户端断开连接的时候的hook函数
func OnConnectionLost(conn kiface.IConnection) {
	//获取当前连接的Pid属性
	fmt.Println("====>OnConnectionLost =====")
	pid, _ := conn.GetProperty("pid")

	//根据pid获取对应的玩家对象
	user := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//触发玩家下线业务
	if pid != nil {
		user.LostConnection()
	}

	fmt.Println("====> Player ", pid, " left =====")

}

func main() {
	s := knet.NewServer()
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)
	s.AddRouter(2, &api.WorldChatAPI{})
	s.AddRouter(3, &api.MoveApi{}) //移动
	s.Serve()
}
