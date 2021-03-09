package peer

import (
	"net"
)

// Address stores address details about each peer
type Address struct {
	IP   net.IP
	Port uint16
}
