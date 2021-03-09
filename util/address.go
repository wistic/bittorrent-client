package util

import (
	"net"
)

// Address stores connection details about each peer
type Address struct {
	IP   net.IP
	Port uint16
}

// NewAddress constructs Address
func NewAddress(ip net.IP, port uint16) *Address {
	return &Address{IP: ip, Port: port}
}
