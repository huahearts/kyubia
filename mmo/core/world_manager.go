package core

import (
	"sync"
)

type WorldManager struct {
	AoiMgr *AOIManager
	Users  map[int32]*User
	pLock  sync.RWMutex
}

var WorldMgrObj *WorldManager

func init() {
	WorldMgrObj = &WorldManager{
		Users:  make(map[int32]*User),
		AoiMgr: NewAOIMgr(AOI_MIN_X, AOI_MAX_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_X, AOI_CNTS_Y),
	}
}

func (wm *WorldManager) AddUser(user *User) {
	wm.pLock.Lock()
	wm.Users[user.Pid] = user
	wm.pLock.Unlock()

	wm.AoiMgr.AddToGridByPos(int(user.Pid), user.X, user.Z)
}

func (wm *WorldManager) RemovePlayerByPid(pid int32) {
	wm.pLock.Lock()
	delete(wm.Users, pid)
	wm.pLock.Unlock()
}

//通过玩家ID 获取对应玩家信息
func (wm *WorldManager) GetPlayerByPid(pid int32) *User {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	return wm.Users[pid]
}

func (wm *WorldManager) GetAllPlayers() []*User {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	//创建返回的player集合切片
	users := make([]*User, 0)

	//添加切片
	for _, v := range wm.Users {
		users = append(users, v)
	}

	//返回
	return users
}
