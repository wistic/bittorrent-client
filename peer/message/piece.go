package message

import (
	"encoding/binary"
	"errors"
)

type Piece struct {
	Index uint32
	Begin uint32
	Block []byte
}

func (piece *Piece) Encode() ([]byte, error) {
	if piece == nil {
		return nil, errors.New("piece is empty")
	}
	length := 9 + len(piece.Block)
	buffer := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buffer[0:4], uint32(length))
	buffer[4] = byte(MsgPiece)
	binary.BigEndian.PutUint32(buffer[5:9], piece.Index)
	binary.BigEndian.PutUint32(buffer[9:13], piece.Begin)
	copy(buffer[13:], piece.Block)
	return buffer, nil
}

func (piece *Piece) Decode(data []byte) error {
	if piece == nil {
		return errors.New("piece is empty")
	} else if data == nil {
		return errors.New("empty data buffer")
	}
	length := binary.BigEndian.Uint32(data[0:4])
	if int(length) != len(data)-4 {
		return errors.New("mismatched length")
	}
	if messageID(data[4]) != MsgPiece {
		return errors.New("not a piece message")
	}
	piece.Index = binary.BigEndian.Uint32(data[5:9])
	piece.Begin = binary.BigEndian.Uint32(data[9:13])
	piece.Block = data[13:]
	return nil
}
