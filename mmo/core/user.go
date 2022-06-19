package core

import (
	"fmt"
	"math/rand"
	"sync"

	//"time"

	"github.com/golang/protobuf/proto"
	"github.com/huahearts/kyubia/kiface"
	"github.com/huahearts/kyubia/mmo/pb"
)

type User struct {
	Pid  int32
	Conn kiface.IConnection
	X    float32
	Y    float32
	Z    float32
	V    float32
}

var PidGen int32 = 1
var IdLock sync.Mutex

func NewUser(conn kiface.IConnection) *User {
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()
	return &User{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)),
		Y:    0,
		Z:    float32(134 + rand.Intn(17)),
		V:    0,
	}
}

func (u *User) SendMsg(msgId uint32, data proto.Message) {
	fmt.Printf("before Marshal data=%+v\n", data)
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err :", err)
		return
	}

	fmt.Printf("after marshal data=%v\n", msg)
	if u.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}

	if err := u.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("usr sendmsg err!")
		return
	}
	return
}

func (u *User) SyncPid() {
	data := &pb.SyncPid{
		Pid: u.Pid,
	}

	u.SendMsg(1, data)
}

func (u *User) BroadCastStartPosition() {

	msg := &pb.BroadCast{
		PID: u.Pid,
		Tp:  2, //TP2 代表广播坐标
		Data: &pb.BroadCast_P{
			&pb.Position{
				X: u.X,
				Y: u.Y,
				Z: u.Z,
				V: u.V,
			},
		},
	}

	u.SendMsg(200, msg)
}

func (u *User) Talk(content string) {
	msg := &pb.BroadCast{
		PID: u.Pid,
		Tp:  1,
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	users := WorldMgrObj.GetAllPlayers()
	for _, user := range users {
		user.SendMsg(200, msg)
	}
}

func (u *User) SyncSurrounding() {
	pids := WorldMgrObj.AoiMgr.GetPIDsByPos(u.X, u.Z)
	users := make([]*User, 0, len(pids))
	for _, pid := range pids {
		users = append(users, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}

	msg := &pb.BroadCast{
		PID: u.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: u.X,
				Y: u.Y,
				Z: u.Z,
				V: u.V,
			},
		},
	}

	for _, user := range users {
		user.SendMsg(200, msg)
	}
	playersData := make([]*pb.Player, 0, len(users))
	for _, player := range users {
		p := &pb.Player{
			PID: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		playersData = append(playersData, p)
	}

	SyncPlayersMsg := &pb.SyncPlayers{
		Ps: playersData[:],
	}

	u.SendMsg(202, SyncPlayersMsg)
}

func (u *User) UpdatePos(x, y, z, v float32) {
	u.X = x
	u.Y = y
	u.Z = z
	u.V = v
	msg := &pb.BroadCast{
		PID: u.Pid,
		Tp:  4,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: u.X,
				Y: u.Y,
				Z: u.Z,
				V: u.V,
			},
		},
	}

	users := u.GetSurroundingPlayers()
	//向周边的每个玩家发送MsgID:200消息，移动位置更新消息
	for _, user := range users {
		user.SendMsg(200, msg)
	}
}

func (u *User) GetSurroundingPlayers() []*User {
	pids := WorldMgrObj.AoiMgr.GetPIDsByPos(u.X, u.Z)
	//将所有pid对应的Player放到Player切片中
	users := make([]*User, 0, len(pids))
	for _, pid := range pids {
		users = append(users, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}

	return users
}
