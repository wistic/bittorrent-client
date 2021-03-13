package message

import (
	"encoding/binary"
)

type Cancel struct {
	Index  uint32
	Begin  uint32
	Length uint32
}

func (cancel *Cancel) GetMessageID() MsgID {
	return MsgCancel
}

func (cancel *Cancel) GetPayload() []byte {
	buffer := make([]byte, 12)
	binary.BigEndian.PutUint32(buffer[0:4], cancel.Index)
	binary.BigEndian.PutUint32(buffer[4:8], cancel.Begin)
	binary.BigEndian.PutUint32(buffer[8:], cancel.Length)
	return buffer
}

func (cancel *Cancel) Deserialize(data []byte) {
	cancel.Index = binary.BigEndian.Uint32(data[0:4])
	cancel.Begin = binary.BigEndian.Uint32(data[4:8])
	cancel.Length = binary.BigEndian.Uint32(data[8:])
}
