package message

import (
	"bittorrent-go/util"
	"errors"
	"io"
)

// Handshake contains data used in protocol handshake
type Handshake struct {
	Protocol  string
	Extension util.Extension
	InfoHash  util.Hash
	PeerID    util.PeerID
}

func WriteHandshake(handshake *Handshake, writer io.Writer) error {
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

func ReadHandshake(reader io.Reader) (*Handshake, error) {
	lengthBuff := make([]byte, 1)
	n, err := reader.Read(lengthBuff)
	if err != nil {
		return nil, err
	} else if n != 1 {
		return nil, errors.New("length buffer empty")
	}
	extension := util.Extension{}
	infoHash := util.Hash{}
	peerID := util.PeerID{}
	protocolLength := int(lengthBuff[0])
	if protocolLength == 0 {
		return nil, errors.New("protocol identifier is empty")
	}
	payloadLength := protocolLength + len(extension.Slice()) + len(infoHash.Slice()) + len(peerID.Slice())
	buff := make([]byte, payloadLength)
	n, err = reader.Read(buff)
	if err != nil {
		return nil, err
	} else if n != payloadLength {
		return nil, errors.New("handshake payload is corrupt")
	}
	protocol := string(buff[0:protocolLength])
	copy(extension.Slice(), buff[protocolLength:])
	copy(infoHash.Slice(), buff[protocolLength+len(extension.Value):])
	copy(peerID.Slice(), buff[protocolLength+len(extension.Value)+len(infoHash.Value):])
	return &Handshake{protocol, extension, infoHash, peerID}, nil
}
