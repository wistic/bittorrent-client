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
	Connection net.Conn
	PeerID     util.PeerID
	Address    string
}

func AttemptConnection(address string, peerID *util.PeerID, infoHash *util.Hash) (*Peer, error) {
	connection, err := net.DialTimeout("tcp", address, 6*time.Second)
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
		Connection: connection,
		PeerID:     receivedHandshake.PeerID,
		Address:    address,
	}
	return &peer, nil
}
