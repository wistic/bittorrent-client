package message

import (
	"encoding/binary"
)

type Piece struct {
	Index uint32
	Begin uint32
	Block []byte
}

func (piece *Piece) Encode() []byte {
	length := 9 + len(piece.Block)
	buffer := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buffer[0:4], uint32(length))
	buffer[4] = byte(MsgPiece)
	binary.BigEndian.PutUint32(buffer[5:9], piece.Index)
	binary.BigEndian.PutUint32(buffer[9:13], piece.Begin)
	copy(buffer[13:], piece.Block)
	return buffer
}

func (piece *Piece) Decode(data []byte) {
	piece.Index = binary.BigEndian.Uint32(data[5:9])
	piece.Begin = binary.BigEndian.Uint32(data[9:13])
	piece.Block = data[13:]
}
