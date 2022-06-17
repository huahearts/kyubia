package knet

import (
	"bytes"
	"encoding/binary"

	"github.com/huahearts/kyubia/kiface"
)

var defaultHeaderLen uint32 = 8

type DataPack struct {
}

func NewDataPacket() kiface.IPacket {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	return defaultHeaderLen
}

func (dp *DataPack) Pack(msg kiface.IMessage) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	if err := binary.Write(buffer, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.LittleEndian, msg.GetID()); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (dp *DataPack) Unpack(data []byte) (kiface.IMessage, error) {
	dataBuff := bytes.NewReader(data)
	msg := &Message{}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}

	//是否超出最大包长
	return msg, nil
}
