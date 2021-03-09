package attribute

import (
	"net"
)

// ConnectionInfo stores connection details about each peer
type ConnectionInfo struct {
	IP   net.IP
	Port uint16
}
