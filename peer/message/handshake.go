package message

import (
	"bittorrent-go/util"
	"errors"
	"io"
)

const protocolIdentifier = "BitTorrent protocol"
const extensionLength = 8
const infoHashLength = 20
const peerIDLength = 20
const handshakeLength = 1 + len(protocolIdentifier) + extensionLength + infoHashLength + peerIDLength

// Handshake contains data used in protocol handshake
type Handshake struct {
	Protocol  string
	Extension util.Extension
	InfoHash  util.Hash
	PeerID    util.PeerID
}

func NewHandshake(protocol *string, extension *util.Extension, infoHash *util.Hash, peerID *util.PeerID) *Handshake {
	return &Handshake{Protocol: *protocol, Extension: *extension, InfoHash: *infoHash, PeerID: *peerID}
}

func WriteHS(handshake *Handshake, writer io.Writer) error {
	if handshake == nil {
		return errors.New("handshake is empty")
	}
	length := 1 + len(handshake.Protocol) + len(handshake.Extension.Value) + len(handshake.InfoHash.Value) + len(handshake.PeerID.Value)
	buff := make([]byte, length)
	index := 0
	buff[index] = byte(len(handshake.Protocol))
	index++
	copy(buff[index:], handshake.Protocol)
	index += len(handshake.Protocol)
	copy(buff[index:], handshake.Extension.Value[:])
	index += len(handshake.Extension.Value)
	copy(buff[index:], handshake.InfoHash.Value[:])
	index += len(handshake.InfoHash.Value)
	copy(buff[index:], handshake.PeerID.Value[:])
	_, err := writer.Write(buff)
	return err
}

func ReadHS(reader io.Reader) (*Handshake, error) {
	lengthBuff := make([]byte, 1)
	_, err := reader.Read(lengthBuff[:])
	if err != nil {
		return nil, err
	}
	extension := util.DefaultExtension()
	infoHash := util.DefaultHash()
	peerID := util.DefaultPeerID()
	protocolLength := int(lengthBuff[0])
	if protocolLength == 0 {
		return nil, errors.New("protocol identifier is empty")
	}
	buff := make([]byte, protocolLength+len(extension.Value)+len(infoHash.Value)+len(peerID.Value))
	_, err = reader.Read(buff[:])
	if err != nil {
		return nil, err
	}
	protocol := string(buff[0:protocolLength])
	copy(extension.Value[:], buff[protocolLength:])
	copy(infoHash.Value[:], buff[protocolLength+len(extension.Value):])
	copy(peerID.Value[:], buff[protocolLength+len(extension.Value)+len(infoHash.Value):])
	return NewHandshake(&protocol, extension, infoHash, peerID), nil
}
