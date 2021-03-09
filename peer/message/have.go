package message

import (
	"encoding/binary"
	"errors"
)

type Have struct {
	Index uint32
}

func (have *Have) Encode() ([]byte, error) {
	if have == nil {
		return nil, errors.New("have is empty")
	}
	length := 5
	buffer := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buffer[0:4], uint32(length))
	buffer[4] = byte(MsgHave)
	binary.BigEndian.PutUint32(buffer[5:9], have.Index)
	return buffer, nil
}

func (have *Have) Decode(data []byte) error {
	if have == nil {
		return errors.New("have is empty")
	} else if data == nil {
		return errors.New("empty data buffer")
	}
	length := binary.BigEndian.Uint32(data[0:4])
	if int(length) != len(data)-4 {
		return errors.New("mismatched length")
	}
	if messageID(data[4]) != MsgHave {
		return errors.New("not a have message")
	}
	have.Index = binary.BigEndian.Uint32(data[5:9])
	return nil
}
