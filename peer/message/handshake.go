package message

import (
	"bittorrent-go/util"
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
	Extension util.Extension
	InfoHash  util.Hash
	PeerID    util.PeerID
}

func NewHandshake(extension util.Extension, infoHash util.Hash, peerID util.PeerID) *Handshake {
	return &Handshake{Extension: extension, InfoHash: infoHash, PeerID: peerID}
}

// EncodeHandshake the Handshake struct
func EncodeHandshake(handshake *Handshake) ([]byte, error) {
	if handshake == nil {
		return nil, errors.New("handshake is empty")
	}
	buffer := make([]byte, handshakeLength)
	index := 0
	buffer[index] = byte(len(protocolIdentifier))
	index++
	copy(buffer[index:], protocolIdentifier)
	index += len(protocolIdentifier)
	copy(buffer[index:], handshake.Extension.Value[:])
	index += len(handshake.Extension.Value)
	copy(buffer[index:], handshake.InfoHash.Value[:])
	index += len(handshake.InfoHash.Value)
	copy(buffer[index:], handshake.PeerID.Value[:])
	return buffer, nil
}

// DecodeHandshake the Handshake struct
func DecodeHandshake(data []byte) (*Handshake, error) {
	if len(data) != handshakeLength {
		return nil, errors.New("Handshake length  mismatched")
	}
	if data[0] != byte(len(protocolIdentifier)) || bytes.Equal(data[1:len(protocolIdentifier)], []byte(protocolIdentifier)) {
		return nil, errors.New("protocol id mismatch")
	}
	extension := util.DefaultExtension()
	infoHash := util.DefaultHash()
	peerID := util.DefaultPeerID()
	copy(extension.Value[:], data[1+len(protocolIdentifier):])
	copy(infoHash.Value[:], data[1+len(protocolIdentifier)+len(extension.Value):])
	copy(peerID.Value[:], data[1+len(protocolIdentifier)+len(extension.Value)+len(infoHash.Value):])
	return NewHandshake(*extension, *infoHash, *peerID), nil
}
