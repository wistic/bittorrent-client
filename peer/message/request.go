package message

import (
	"encoding/binary"
)

type Request struct {
	Index  uint32
	Begin  uint32
	Length uint32
}

func (request *Request) GetMessageID() MsgID {
	return MsgRequest
}

func (request *Request) GetPayload() []byte {
	buffer := make([]byte, 12)
	binary.BigEndian.PutUint32(buffer[0:4], request.Index)
	binary.BigEndian.PutUint32(buffer[4:8], request.Begin)
	binary.BigEndian.PutUint32(buffer[8:], request.Length)
	return buffer
}

func (request *Request) Decode(data []byte) {
	request.Index = binary.BigEndian.Uint32(data[5:9])
	request.Begin = binary.BigEndian.Uint32(data[9:13])
	request.Length = binary.BigEndian.Uint32(data[13:])
}
