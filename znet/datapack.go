package znet

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
	"zinx/config"
	"zinx/ziface"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeaderLen() uint32 {
	// DataLen uint32 + Id uint32 == 64 bits == 8 bytes
	return 8
}
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	if err := binary.Write(buffer, binary.BigEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	binary.Write(buffer, binary.BigEndian, msg.GetData())
	return buffer.Bytes(), nil
}
func (dp *DataPack) UnpackHeader(data []byte) (ziface.IMessage, error) {
	reader := bytes.NewReader(data)
	msg := &Message{}
	if err := binary.Read(reader, binary.BigEndian, &msg.Id); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.BigEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	if config.GlobalConfig.MaxPackageSize > 0 && msg.DataLen > config.GlobalConfig.MaxPackageSize {
		return nil, errors.New("too large msg data recv!")
	}

	return msg, nil

}
