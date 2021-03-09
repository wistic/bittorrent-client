package message

import (
	"encoding/binary"
	"errors"
)

type Request struct {
	Index  uint32
	Begin  uint32
	Length uint32
}

func (request *Request) Encode() ([]byte, error) {
	if request == nil {
		return nil, errors.New("request is empty")
	}
	length := 13
	buffer := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buffer[0:4], uint32(length))
	buffer[4] = byte(MsgRequest)
	binary.BigEndian.PutUint32(buffer[5:9], request.Index)
	binary.BigEndian.PutUint32(buffer[9:13], request.Begin)
	binary.BigEndian.PutUint32(buffer[13:], request.Length)
	return buffer, nil
}

func (request *Request) Decode(data []byte) error {
	if request == nil {
		return errors.New("request is empty")
	} else if data == nil {
		return errors.New("empty data buffer")
	}
	length := binary.BigEndian.Uint32(data[0:4])
	if int(length) != len(data)-4 {
		return errors.New("mismatched length")
	}
	if messageID(data[4]) != MsgRequest {
		return errors.New("not a request message")
	}
	request.Index = binary.BigEndian.Uint32(data[5:9])
	request.Begin = binary.BigEndian.Uint32(data[9:13])
	request.Length = binary.BigEndian.Uint32(data[13:])
	return nil
}
