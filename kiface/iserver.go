package kiface

type IServer interface {
	Start()
	Stop()
	Serve()
	AddRouter(msgID uint32, router IRouter)
	GetConnMgr() IConnMgr
	SetOnConnStart(func(IConnection))
	SetOnConnStop(func(IConnection))
	OnConnStartCallback(conn IConnection)
	OnConnStopCallback(conn IConnection)
	Packet() Packet
}
