package util

import "bytes"

// Hash is the common hash used in bittorrent protocol
type Hash struct {
	Value [20]byte
}

// Slice returns the hash value as slice
func (hash *Hash) Slice() []byte {
	return hash.Value[:]
}

func (hash *Hash) Match(other *Hash) bool {
	return bytes.Equal(hash.Slice(), other.Slice())
}
