package util

// Hash is the common hash used in bittorrent protocol
type Hash struct {
	Value [20]byte
}

// Slice returns the hash value as slice
func (hash Hash) Slice() []byte {
	return hash.Value[:]
}
