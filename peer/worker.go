package peer

import (
	"bittorrent-go/peer/message"
	"bittorrent-go/util"
	"fmt"
	"github.com/kr/pretty"
	"net"
	"time"
)

func Worker(address *util.Address, peerID *util.PeerID, infoHash *util.Hash) {
	connection, err := net.DialTimeout("tcp", address.String(), 6*time.Second)
	if err != nil {
		fmt.Println("connection error:", err)
		return
	}
	defer connection.Close()
	handshake := message.Handshake{Protocol: "BitTorrent protocol", Extension: util.Extension{}, InfoHash: *infoHash, PeerID: *peerID}
	err = message.WriteHandshake(&handshake, connection)
	if err != nil {
		fmt.Println("handshake send error:", err)
		return
	}
	receivedHandshake, err := message.ReadHandshake(connection)
	if err != nil {
		fmt.Println("handshake receive error:", err)
		return
	}
	pretty.Println(handshake)
	pretty.Println(receivedHandshake)
	pretty.Diff(handshake, receivedHandshake)
}
