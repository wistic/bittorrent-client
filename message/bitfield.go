package message

import (
	"bittorrent-go/util"
)

type BitField struct {
	Field util.BitField
}

func (bitfield *BitField) GetMessageID() MsgID {
	return MsgBitfield
}

func (bitfield *BitField) GetPayload() []byte {
	buffer := make([]byte, len(bitfield.Field.Value))
	copy(buffer[:], bitfield.Field.Value)
	return buffer
}

func (bitfield *BitField) Deserialize(data []byte) {
	bitfield.Field.Value = data
}
