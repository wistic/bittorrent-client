package peer

import (
	"bittorrent-go/message"
	"bittorrent-go/util"
	"bytes"
	"errors"
	"net"
	"time"
)

const protocolIdentifier = "BitTorrent protocol"

type Peer struct {
	Connection     net.Conn
	PeerID         util.PeerID
	ConnectionInfo *util.Address
}

func AttemptConnection(address *util.Address, peerID *util.PeerID, infoHash *util.Hash) (*Peer, error) {
	connection, err := net.DialTimeout("tcp", address.String(), 6*time.Second)
	if err != nil {
		return nil, err
	}

	connection.SetDeadline(time.Now().Add(5 * time.Second))
	defer connection.SetDeadline(time.Time{})

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
	peer := Peer{
		Connection:     connection,
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
