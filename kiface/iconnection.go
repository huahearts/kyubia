package kiface

import (
	"context"
	"net"
)

type IConnection interface {
	Start()
	Stop()
	Context() context.Context
	GetTCPCOnnection() *net.TCPConn
	GetConnID() uint32
	RemoteAddr() net.Addr

	SendMsg(msgId uint32, data []byte) error
	SendBuffMsg(msgId uint32, data []byte) error

	//connection property
	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, error)
	RemoveProperty(key string)
}
