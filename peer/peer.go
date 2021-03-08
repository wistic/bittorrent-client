package peer

import (
	"fmt"
	"net"
)

// Peer stores details about each peer
type Peer struct {
	IP   net.IP
	Port uint16
}

// String implements Stringer interface to properly print Peer struct
func (p Peer) String() string {
	return fmt.Sprintf("IP: %v\tPort: %v\n", p.IP, p.Port)
}
