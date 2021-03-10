package message

import (
	"encoding/binary"
)

type Have struct {
	Index uint32
}

func (have *Have) Encode() []byte {
	length := 5
	buffer := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buffer[0:4], uint32(length))
	buffer[4] = byte(MsgHave)
	binary.BigEndian.PutUint32(buffer[5:9], have.Index)
	return buffer
}

func (have *Have) Decode(data []byte) {
	have.Index = binary.BigEndian.Uint32(data[5:9])
}
