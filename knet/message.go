package knet

import "github.com/huahearts/kyubia/kiface"

type Message struct {
	DataLen uint32
	ID      uint32
	Data    []byte
}

func NewMsgPacket(ID uint32, data []byte) kiface.IMessage {
	return &Message{
		DataLen: uint32(len(data)),
		ID:      ID,
		Data:    data,
	}
}

func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}
func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) GetID() uint32 {
	return m.ID
}

func (m *Message) SetMsgID(msgID uint32) {
	m.ID = msgID
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}
