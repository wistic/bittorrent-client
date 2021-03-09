package util

import "errors"

// BitField represents which pieces the peer has
type BitField struct {
	Value []byte
}

// CheckPiece tells if a peer has the piece at the given index
func (bitfield *BitField) CheckPiece(index int) (bool, error) {
	byteIndex := index / 8
	byteOffset := index % 8
	if byteIndex < 0 || byteIndex >= len(bitfield.Value) {
		return false, errors.New("piece index out of bounds")
	}
	check := (bitfield.Value[byteIndex]>>(7-byteOffset))&1 == 1
	return check, nil
}

// SetPiece sets a bit to 1 if the peer has the piece at that index
func (bitfield *BitField) SetPiece(index int) error {
	byteIndex := index / 8
	byteOffset := index % 8

	if byteIndex < 0 || byteIndex >= len(bitfield.Value) {
		return errors.New("piece index out of bounds")
	}
	bitfield.Value[byteIndex] = bitfield.Value[byteIndex] | (1 << (7 - byteOffset))
	return nil
}

func (bitfield *BitField) Slice() []byte {
	return bitfield.Value[:]
}
