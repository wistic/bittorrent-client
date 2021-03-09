package message

import (
	"bytes"
	"errors"
)

const protocolIdentifier = "BitTorrent protocol"
const extensionLength = 8
const infoHashLength = 20
const peerIDLength = 20
const handshakeLength = 1 + len(protocolIdentifier) + extensionLength + infoHashLength + peerIDLength

// Handshake contains data used in protocol handshake
type Handshake struct {
	Extension [extensionLength]byte
	InfoHash  [infoHashLength]byte
	PeerID    [peerIDLength]byte
}

// Encode the Handshake struct
func (handshake *Handshake) Encode() ([]byte, error) {
	if handshake == nil {
		return nil, errors.New("handshake is empty")
	}
	buffer := make([]byte, handshakeLength)
	index := 0
	buffer[index] = byte(len(protocolIdentifier))
	index++
	copy(buffer[index:], protocolIdentifier)
	index += len(protocolIdentifier)
	copy(buffer[index:], handshake.Extension[:])
	index += extensionLength
	copy(buffer[index:], handshake.InfoHash[:])
	index += len(handshake.InfoHash)
	copy(buffer[index:], handshake.PeerID[:])
	return buffer, nil
}

// Decode the Handshake struct
func (handshake *Handshake) Decode(data []byte) error {
	if len(data) != handshakeLength {
		return errors.New("Handshake length  mismatched")
	}
	if data[0] != byte(len(protocolIdentifier)) || bytes.Equal(data[1:len(protocolIdentifier)], []byte(protocolIdentifier)) {
		return errors.New("protocol id mismatch")
	}
	extension := [extensionLength]byte{}
	infoHash := [infoHashLength]byte{}
	peerID := [peerIDLength]byte{}
	copy(extension[:], data[1+len(protocolIdentifier):])
	copy(infoHash[:], data[1+len(protocolIdentifier)+extensionLength:])
	copy(peerID[:], data[1+len(protocolIdentifier)+extensionLength+infoHashLength:])
	handshake.Extension, handshake.InfoHash, handshake.PeerID = extension, infoHash, peerID
	return nil
}
