package knet

import (
	"errors"
	"fmt"
	"sync"

	"github.com/huahearts/kyubia/kiface"
)

type ConnMgr struct {
	conns   map[uint32]kiface.IConnection
	connMtx sync.RWMutex
}

func NewConnMgr() *ConnMgr {
	return &ConnMgr{
		conns: make(map[uint32]kiface.IConnection),
	}
}

func (cm *ConnMgr) Add(conn kiface.IConnection) {
	{
		cm.connMtx.Lock()
		defer cm.connMtx.Unlock()
		cm.conns[conn.GetConnID()] = conn
	}
	fmt.Println("connection add to connmgr succ,conn num:", cm.Len())
}

func (cm *ConnMgr) Remove(conn kiface.IConnection) {
	{
		cm.connMtx.Lock()
		defer cm.connMtx.Unlock()
		delete(cm.conns, conn.GetConnID())
	}
	fmt.Println("connection del to connmgr succ,conn num:", cm.Len())
}

func (cm *ConnMgr) Get(connID uint32) (kiface.IConnection, error) {
	cm.connMtx.RLock()
	defer cm.connMtx.RUnlock()
	if v, ok := cm.conns[connID]; ok {
		return v, nil
	}
	return nil, errors.New("conn is not exist")
}

func (cm *ConnMgr) Len() int {
	cm.connMtx.RLock()
	defer cm.connMtx.RUnlock()
	return len(cm.conns)
}

func (cm *ConnMgr) ClearConn() {
	{
		cm.connMtx.Lock()
		defer cm.connMtx.Unlock()
		for connID, conn := range cm.conns {
			conn.Stop()
			delete(cm.conns, connID)
		}
	}

	fmt.Println("Clear All Connections successfully: conn num = ", cm.Len())
}

func (cm *ConnMgr) ClearOneConn(connID uint32) {
	{
		cm.connMtx.Lock()
		defer cm.connMtx.Unlock()
		if conn, ok := cm.conns[connID]; ok {
			conn.Stop()
			delete(cm.conns, connID)
			fmt.Println("Clear Connections ID:  ", connID, "succeed")
		}
		return
	}
	fmt.Println("Clear Connections ID:  ", connID, "err")
}
