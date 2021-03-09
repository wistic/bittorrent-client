package util

// Hash is the common hash used in bittorrent protocol
type Hash struct {
	Value [20]byte
}

// NewHash constructs Hash
func NewHash(value [20]byte) *Hash {
	return &Hash{Value: value}
}

func DefaultHash() *Hash {
	return NewHash([20]byte{})
}

// Slice returns the hash value as slice
func (hash Hash) Slice() []byte {
	return hash.Value[:]
}
