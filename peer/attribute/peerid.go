package attribute

import (
	"math/rand"
	"strconv"
	"time"
)

const max = 999999999999
const min = 100000000000

type PeerID [20]byte

// GeneratePeerID generates a random PeerID for our client
func GeneratePeerID() PeerID {
	rand.Seed(time.Now().UnixNano())
	peerID := [20]byte{}
	copy(peerID[:], "-BG0001-"+strconv.Itoa(rand.Intn(max-min)+min))
	return peerID
}
