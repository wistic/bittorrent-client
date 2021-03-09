package util

import (
	"fmt"
	"net"
)

// Address stores connection details about each peer
type Address struct {
	IP   net.IP
	Port uint16
}

func (address *Address) String() string {
	return fmt.Sprint(address.IP.String(), ":", address.Port)
}
