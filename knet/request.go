package knet

import "github.com/huahearts/kyubigo/kiface"

type Request struct {
	conn kiface.IConnection
	msg  kiface.IMessage
}

func (r *Request) GetConnection() kiface.IConnection {
	return r.conn
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetID()
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}
