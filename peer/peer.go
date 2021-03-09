package peer

import (
	"bittorrent-go/util"
	"net"
)

type Peer struct {
	Connection     net.Conn
	AmChoking      bool
	AmInterested   bool
	PeerChoking    bool
	PeerInterested bool
	BitField       util.BitField
	PeerID         util.PeerID
	ConnectionInfo util.ConnectionInfo
}
