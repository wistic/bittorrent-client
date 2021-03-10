package peer

import (
	"bittorrent-go/peer/message"
	"bittorrent-go/util"
	"net"
	"time"
)

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
	handshake := message.Handshake{Protocol: "BitTorrent protocol", Extension: util.Extension{}, InfoHash: *infoHash, PeerID: *peerID}
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
