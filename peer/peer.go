package peer

import (
	"net"
)

// PCInfo stores connection details about each peer
type PCInfo struct {
	IP   net.IP
	Port uint16
}
