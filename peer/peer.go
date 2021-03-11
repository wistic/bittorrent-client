package peer

import (
	"bittorrent-go/peer/message"
	"bittorrent-go/util"
	"bytes"
	"errors"
	"fmt"
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

func exchangeBitFields(connection net.Conn, field util.BitField) (util.BitField, error) {
	go func() {
		connection.SetDeadline(time.Now().Add(5 * time.Second))
		bitfield := message.BitField{Field: field}
		_ = message.SendMessage(&bitfield, connection)
		connection.SetDeadline(time.Time{})
	}()
	connection.SetDeadline(time.Now().Add(5 * time.Second))
	bitfield, err := message.ReceiveMessage(connection)
	connection.SetDeadline(time.Time{})
	if err != nil {
		return util.BitField{}, err
	}
	id := bitfield.GetMessageID()
	if id != message.MsgBitfield {
		return util.BitField{}, errors.New("not a bitfield message")
	}
	return bitfield.(*message.BitField).Field, err
}

func AttemptConnection(address *util.Address, peerID *util.PeerID, infoHash *util.Hash) (*Peer, error) {
	connection, err := net.DialTimeout("tcp", address.String(), 6*time.Second)
	if err != nil {
		return nil, err
	}
	_ = connection.SetDeadline(time.Now().Add(5 * time.Second))
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
	if receivedHandshake.Protocol != protocolIdentifier || !bytes.Equal(receivedHandshake.InfoHash.Slice(), infoHash.Slice()) {
		connection.Close()
		return nil, errors.New("bad handshake")
	}
	connection.SetDeadline(time.Time{})
	bitfield, err := exchangeBitFields(connection, util.BitField{})
	if err != nil {
		fmt.Println(err)
	}
	peer := Peer{
		Connection:     connection,
		AmChoking:      true,
		AmInterested:   false,
		PeerChoking:    true,
		PeerInterested: false,
		BitField:       bitfield,
		PeerID:         receivedHandshake.PeerID,
		ConnectionInfo: address,
	}
	return &peer, nil
}

func (peer *Peer) Send(data message.Message) error {
	err := message.SendMessage(data, peer.Connection)
	return err
}

func (peer *Peer) Receive() (message.Message, error) {
	data, err := message.ReceiveMessage(peer.Connection)
	return data, err
}
