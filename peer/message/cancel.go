package message

import (
	"encoding/binary"
	"errors"
)

type Cancel struct {
	Index  uint32
	Begin  uint32
	Length uint32
}

func (cancel *Cancel) Encode() ([]byte, error) {
	if cancel == nil {
		return nil, errors.New("cancel is empty")
	}
	length := 13
	buffer := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buffer[0:4], uint32(length))
	buffer[4] = byte(MsgCancel)
	binary.BigEndian.PutUint32(buffer[5:9], cancel.Index)
	binary.BigEndian.PutUint32(buffer[9:13], cancel.Begin)
	binary.BigEndian.PutUint32(buffer[13:], cancel.Length)
	return buffer, nil
}

func (cancel *Cancel) Decode(data []byte) error {
	if cancel == nil {
		return errors.New("cancel is empty")
	} else if data == nil {
		return errors.New("empty data buffer")
	}
	length := binary.BigEndian.Uint32(data[0:4])
	if int(length) != len(data)-4 {
		return errors.New("mismatched length")
	}
	if messageID(data[4]) != MsgCancel {
		return errors.New("not a cancel message")
	}
	cancel.Index = binary.BigEndian.Uint32(data[5:9])
	cancel.Begin = binary.BigEndian.Uint32(data[9:13])
	cancel.Length = binary.BigEndian.Uint32(data[13:])
	return nil
}
