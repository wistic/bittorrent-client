package message

import (
	"bittorrent-go/util"
	"bufio"
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

func NewHandshake(extension *util.Extension, infoHash *util.Hash, peerID *util.PeerID) *Handshake {
	return &Handshake{Extension: *extension, InfoHash: *infoHash, PeerID: *peerID}
}

func WriteHandshake(handshake *Handshake, writer *bufio.Writer) error {
	if handshake == nil {
		return errors.New("handshake is empty")
	}
	identifierLength := [1]byte{byte(len(protocolIdentifier))}
	_, err := writer.Write(identifierLength[:])
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(protocolIdentifier))
	if err != nil {
		return err
	}
	_, err = writer.Write(handshake.Extension.Value[:])
	if err != nil {
		return err
	}
	_, err = writer.Write(handshake.InfoHash.Value[:])
	if err != nil {
		return err
	}
	_, err = writer.Write(handshake.PeerID.Value[:])
	if err != nil {
		return err
	}
	return nil
}

func ReadHandshake(reader *bufio.Reader) (*Handshake, error) {
	identifierLength := [1]byte{}
	_, err := reader.Read(identifierLength[:])
	if err != nil {
		return nil, err
	}
	if identifierLength[0] != byte(len(protocolIdentifier)) {
		return nil, errors.New("protocol identifier length  mismatched")
	}
	protocolIdentifierReceived := [len(protocolIdentifier)]byte{}
	_, err = reader.Read(protocolIdentifierReceived[:])
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(protocolIdentifierReceived[:], []byte(protocolIdentifier)) {
		return nil, errors.New("protocol identifier mismatched")
	}
	extension := util.DefaultExtension()
	_, err = reader.Read(extension.Value[:])
	if err != nil {
		return nil, err
	}
	infoHash := util.DefaultHash()
	_, err = reader.Read(infoHash.Value[:])
	if err != nil {
		return nil, err
	}
	peerID := util.DefaultPeerID()
	_, err = reader.Read(peerID.Value[:])
	if err != nil {
		return nil, err
	}
	return NewHandshake(extension, infoHash, peerID), nil
}
