package peer

import (
	"bittorrent-go/peer/message"
	"bittorrent-go/util"
	"encoding/binary"
	"errors"
	"net"
	"time"
)

const protocolIdentifier = "BitTorrent protocol"

type Peer struct {
	Connection     net.Conn
	AmChoking      bool
	AmInterested   bool
	PeerChoking    bool
	PeerInterested bool
	BitField       util.BitField
	PeerID         util.PeerID
	ConnectionInfo *util.Address
}

func AttemptConnection(address *util.Address, peerID *util.PeerID, infoHash *util.Hash) (*Peer, error) {
	connection, err := net.DialTimeout("tcp", address.String(), 6*time.Second)
	if err != nil {
		return nil, err
	}
	handshake := message.Handshake{Protocol: protocolIdentifier, Extension: util.Extension{}, InfoHash: *infoHash, PeerID: *peerID}
	err = message.WriteHandshake(&handshake, connection)
	if err != nil {
		_ = connection.Close()
		return nil, err
	}
	receivedHandshake, err := message.ReadHandshake(connection)
	if err != nil {
		_ = connection.Close()
		return nil, err
	}
	peer := Peer{
		Connection:     connection,
		AmChoking:      true,
		AmInterested:   false,
		PeerChoking:    true,
		PeerInterested: false,
		BitField:       util.BitField{},
		PeerID:         receivedHandshake.PeerID,
		ConnectionInfo: address,
	}
	return &peer, nil
}

func (peer *Peer) Send(message message.Message) error {
	if message == nil {
		packet := make([]byte, 4)
		_, err := peer.Connection.Write(packet) // keep-alive
		return err
	}
	messageID := message.GetMessageID()
	payload := message.GetPayload()
	if payload == nil {
		length := uint32(1)
		packet := make([]byte, 4+length)
		binary.BigEndian.PutUint32(packet[0:4], length)
		packet[4] = byte(messageID)
		_, err := peer.Connection.Write(packet)
		return err
	}
	length := uint32(1 + len(payload))
	packet := make([]byte, 4+length)
	binary.BigEndian.PutUint32(packet[0:4], length)
	packet[4] = byte(messageID)
	copy(packet[5:], payload)
	_, err := peer.Connection.Write(packet)
	return err
}

func (peer *Peer) Receive() (message.Message, error) {
	lengthBuff := make([]byte, 4)
	n, err := peer.Connection.Read(lengthBuff)
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
	n, err = peer.Connection.Read(packet)
	if err != nil {
		return nil, err
	} else if n != int(length) {
		return nil, errors.New("packet payload corrupt")
	}
	switch message.MsgID(packet[0]) {
	case message.MsgChoke:
		return &message.Choke{}, nil
	case message.MsgUnchoke:
		return &message.Unchoke{}, nil
	case message.MsgInterested:
		return &message.Interested{}, nil
	case message.MsgNotInterested:
		return &message.NotInterested{}, nil
	case message.MsgBitfield:
		bitfield := message.BitField{}
		bitfield.Deserialize(packet[1:])
		return &bitfield, nil
	case message.MsgRequest:
		request := message.Request{}
		request.Deserialize(packet[1:])
		return &request, nil
	case message.MsgPiece:
		piece := message.Piece{}
		piece.Deserialize(packet[1:])
		return &piece, nil
	case message.MsgCancel:
		cancel := message.Cancel{}
		cancel.Deserialize(packet[1:])
		return &cancel, nil
	default:
		return nil, errors.New("unexpected message type")
	}
}
