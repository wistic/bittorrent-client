package message

import (
	"bittorrent-go/util"
	"encoding/binary"
	"errors"
)

type BitField struct {
	Field util.BitField
}

func (bitfield *BitField) Encode() ([]byte, error) {
	if bitfield == nil {
		return nil, errors.New("bitfield is empty")
	}
	length := 1 + len(bitfield.Field)
	buffer := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buffer[0:4], uint32(length))
	buffer[4] = byte(MsgBitfield)
	copy(buffer[5:], bitfield.Field)
	return buffer, nil
}

func (bitfield *BitField) Decode(data []byte) error {
	if bitfield == nil {
		return errors.New("bitfield is empty")
	} else if data == nil {
		return errors.New("empty data buffer")
	}
	length := binary.BigEndian.Uint32(data[0:4])
	if int(length) != len(data)-4 {
		return errors.New("mismatched length")
	}
	if messageID(data[4]) != MsgBitfield {
		return errors.New("not a bitfield message")
	}
	bitfield.Field = data[5:]
	return nil
}
