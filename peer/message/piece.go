package message

import (
	"encoding/binary"
)

type Piece struct {
	Index uint32
	Begin uint32
	Block []byte
}

func (piece *Piece) GetMessageID() MsgID {
	return MsgPiece
}

func (piece *Piece) GetPayload() []byte {
	length := 8 + len(piece.Block)
	buffer := make([]byte, length)
	binary.BigEndian.PutUint32(buffer[0:4], piece.Index)
	binary.BigEndian.PutUint32(buffer[4:8], piece.Begin)
	copy(buffer[8:], piece.Block)
	return buffer
}

func (piece *Piece) Deserialize(data []byte) {
	piece.Index = binary.BigEndian.Uint32(data[0:4])
	piece.Begin = binary.BigEndian.Uint32(data[4:8])
	piece.Block = data[8:]
}
