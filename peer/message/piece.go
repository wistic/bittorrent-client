package message

import (
	"encoding/binary"
	"errors"
)

type Piece struct {
	ID    messageID
	Index uint32
	Begin uint32
	Block []byte
}

func (piece *Piece) Encode() ([]byte, error) {
	if piece == nil {
		return nil, errors.New("piece is empty")
	}
	length := 13 + len(piece.Block)
	buffer := make([]byte, length)
	binary.BigEndian.PutUint32(buffer[0:4], uint32(length))
	buffer[4] = byte(MsgPiece)
	binary.BigEndian.PutUint32(buffer[5:9], piece.Index)
	binary.BigEndian.PutUint32(buffer[9:13], piece.Begin)
	copy(buffer[13:], piece.Block)
	return buffer, nil
}
