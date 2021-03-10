package message

import (
	"encoding/binary"
)

type Request struct {
	Index  uint32
	Begin  uint32
	Length uint32
}

func (request *Request) Encode() []byte {
	length := 13
	buffer := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buffer[0:4], uint32(length))
	buffer[4] = byte(MsgRequest)
	binary.BigEndian.PutUint32(buffer[5:9], request.Index)
	binary.BigEndian.PutUint32(buffer[9:13], request.Begin)
	binary.BigEndian.PutUint32(buffer[13:], request.Length)
	return buffer
}

func (request *Request) Decode(data []byte) {
	request.Index = binary.BigEndian.Uint32(data[5:9])
	request.Begin = binary.BigEndian.Uint32(data[9:13])
	request.Length = binary.BigEndian.Uint32(data[13:])
}
