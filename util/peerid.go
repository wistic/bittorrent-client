package util

import (
	"math/rand"
	"strconv"
	"time"
)

const max = 999999999999
const min = 100000000000

type PeerID struct {
	Value [20]byte
}

// NewPeerID constructs PeerID
func NewPeerID(value [20]byte) *PeerID {
	return &PeerID{Value: value}
}

func DefaultPeerID() *PeerID {
	return NewPeerID([20]byte{})
}

// GeneratePeerID generates a random PeerID for our client
func GeneratePeerID() *PeerID {
	rand.Seed(time.Now().UnixNano())
	peerID := [20]byte{}
	copy(peerID[:], "-BG0001-"+strconv.Itoa(rand.Intn(max-min)+min))
	return NewPeerID(peerID)
}

// String
func (peerID *PeerID) String() string {
	return string(peerID.Value[:])
}
