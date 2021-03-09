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
func Encode(handshake Handshake) []byte {
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
	return buffer
}

// Decode the Handshake struct
func Decode(data []byte) (Handshake, error) {
	if len(data) != handshakeLength {
		return Handshake{}, errors.New("Handshake length  mismatched")
	}
	if data[0] != byte(len(protocolIdentifier)) || bytes.Equal(data[1:len(protocolIdentifier)], []byte(protocolIdentifier)) {
		return Handshake{}, errors.New("Protocol Id mismatch")
	}
	extension := [extensionLength]byte{}
	infoHash := [infoHashLength]byte{}
	peerID := [peerIDLength]byte{}
	copy(extension[0:extensionLength], data[1+len(protocolIdentifier):])
	copy(infoHash[0:infoHashLength], data[1+len(protocolIdentifier)+extensionLength:])
	copy(peerID[0:peerIDLength], data[1+len(protocolIdentifier)+extensionLength+infoHashLength:])
	return Handshake{Extension: extension, InfoHash: infoHash, PeerID: peerID}, nil
}
