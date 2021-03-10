package message

import (
	"encoding/binary"
)

type Cancel struct {
	Index  uint32
	Begin  uint32
	Length uint32
}

func (cancel *Cancel) Encode() []byte {
	length := 13
	buffer := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buffer[0:4], uint32(length))
	buffer[4] = byte(MsgCancel)
	binary.BigEndian.PutUint32(buffer[5:9], cancel.Index)
	binary.BigEndian.PutUint32(buffer[9:13], cancel.Begin)
	binary.BigEndian.PutUint32(buffer[13:], cancel.Length)
	return buffer
}

func (cancel *Cancel) Decode(data []byte) {
	cancel.Index = binary.BigEndian.Uint32(data[5:9])
	cancel.Begin = binary.BigEndian.Uint32(data[9:13])
	cancel.Length = binary.BigEndian.Uint32(data[13:])
}
