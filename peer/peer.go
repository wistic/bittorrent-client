package peer

import (
	"bittorrent-go/peer/attribute"
	"net"
)

type Peer struct {
	Connection     net.Conn
	AmChoking      bool
	AmInterested   bool
	PeerChoking    bool
	PeerInterested bool
	BitField       attribute.BitField
	PeerID         attribute.PeerID
	ConnectionInfo attribute.ConnectionInfo
}
