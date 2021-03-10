package message

import (
	"encoding/binary"
)

type Have struct {
	Index uint32
}

func (have *Have) GetMessageID() MsgID {
	return MsgHave
}

func (have *Have) GetPayload() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer[:], have.Index)
	return buffer
}

func (have *Have) Decode(data []byte) {
	have.Index = binary.BigEndian.Uint32(data[5:9])
}
