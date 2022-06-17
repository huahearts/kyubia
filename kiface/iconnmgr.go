package kiface

// 链接管理抽象
type IConnMgr interface {
	Add(conn IConnection)
	Remove(conn IConnection)
	Get(connId uint32) (IConnection, error)
	Len() int
	ClearConn()
}
