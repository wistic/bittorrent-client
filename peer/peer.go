package peer

import (
	"net"
)

// Peer stores details about each peer
type Peer struct {
	IP   net.IP
	Port uint16
}
