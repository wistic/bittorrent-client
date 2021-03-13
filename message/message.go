package message

import (
	"encoding/binary"
	"errors"
	"io"
)

type MsgID uint8

type Message interface {
	GetMessageID() MsgID
	GetPayload() []byte
}

const (
	// MsgChoke chokes the receiver
	MsgChoke MsgID = 0
	// MsgUnchoke unchokes the receiver
	MsgUnchoke MsgID = 1
	// MsgInterested expresses interest in receiving data
	MsgInterested MsgID = 2
	// MsgNotInterested expresses disinterest in receiving data
	MsgNotInterested MsgID = 3
	// MsgHave alerts the receiver that the sender has downloaded a piece
	MsgHave MsgID = 4
	// MsgBitfield encodes which pieces that the sender has downloaded
	MsgBitfield MsgID = 5
	// MsgRequest requests a block of data from the receiver
	MsgRequest MsgID = 6
	// MsgPiece delivers a block of data to fulfill a request
	MsgPiece MsgID = 7
	// MsgCancel cancels a request
	MsgCancel MsgID = 8
)

func SendMessage(message Message, writer io.Writer) error {
	if message == nil {
		packet := make([]byte, 4)
		_, err := writer.Write(packet) // keep-alive
		return err
	}
	messageID := message.GetMessageID()
	payload := message.GetPayload()
	if payload == nil {
		length := uint32(1)
		packet := make([]byte, 4+length)
		binary.BigEndian.PutUint32(packet[0:4], length)
		packet[4] = byte(messageID)
		_, err := writer.Write(packet)
		return err
	}
	length := uint32(1 + len(payload))
	packet := make([]byte, 4+length)
	binary.BigEndian.PutUint32(packet[0:4], length)
	packet[4] = byte(messageID)
	copy(packet[5:], payload)
	_, err := writer.Write(packet)
	return err
}

func ReceiveMessage(reader io.Reader) (Message, error) {
	lengthBuff := make([]byte, 4)
	n, err := reader.Read(lengthBuff)
	if err != nil {
		return nil, err
	} else if n != 4 {
		return nil, errors.New("length buffer corrupt")
	}
	length := binary.BigEndian.Uint32(lengthBuff)
	if length == 0 {
		return nil, nil // keep-alive
	}
	packet := make([]byte, length)
	n, err = reader.Read(packet)
	if err != nil {
		return nil, err
	} else if n != int(length) {
		return nil, errors.New("packet payload corrupt")
	}
	switch MsgID(packet[0]) {
	case MsgChoke:
		return &Choke{}, nil
	case MsgUnchoke:
		return &Unchoke{}, nil
	case MsgInterested:
		return &Interested{}, nil
	case MsgNotInterested:
		return &NotInterested{}, nil
	case MsgBitfield:
		bitfield := BitField{}
		bitfield.Deserialize(packet[1:])
		return &bitfield, nil
	case MsgRequest:
		request := Request{}
		request.Deserialize(packet[1:])
		return &request, nil
	case MsgPiece:
		piece := Piece{}
		piece.Deserialize(packet[1:])
		return &piece, nil
	case MsgCancel:
		cancel := Cancel{}
		cancel.Deserialize(packet[1:])
		return &cancel, nil
	default:
		return nil, errors.New("unexpected message type")
	}
}
