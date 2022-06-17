package kiface

type IMessage interface {
	GetDataLen() uint32
	GetID() uint32
	GetData() []byte

	SetMsgID(uint32)
	SetData([]byte)
	SetDataLen(uint32)
}
