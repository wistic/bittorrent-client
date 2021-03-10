package message

import (
	"bittorrent-go/util"
	"encoding/binary"
)

type BitField struct {
	Field util.BitField
}

func (bitfield *BitField) Encode() []byte {
	length := 1 + len(bitfield.Field.Value)
	buffer := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buffer[0:4], uint32(length))
	buffer[4] = byte(MsgBitfield)
	copy(buffer[5:], bitfield.Field.Value)
	return buffer
}

func (bitfield *BitField) Decode(data []byte) {
	bitfield.Field.Value = data[5:]
}
