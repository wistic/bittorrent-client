package peer

import (
	"bittorrent-go/peer/message"
	"bittorrent-go/util"
	"bufio"
	"fmt"
	"github.com/kr/pretty"
	"net"
	"time"
)

func Worker(address *util.Address, peerID *util.PeerID, infoHash *util.Hash) {
	fmt.Print(address.String())
	connection, err := net.DialTimeout("tcp", address.String(), 6*time.Second)
	if err != nil {
		fmt.Print("Temp Log:", err)
		return
	}
	defer connection.Close()
	writer := bufio.NewWriter(connection)
	handshake := message.NewHandshake(util.DefaultExtension(), infoHash, peerID)
	err = message.WriteHandshake(handshake, writer)
	if err != nil {
		fmt.Print("Temp Log:", err)
		return
	}
	err = writer.Flush()
	if err != nil {
		fmt.Print("Temp Log:", err)
		return
	}
	reader := bufio.NewReader(connection)
	receivedHandshake, err := message.ReadHandshake(reader)
	if err != nil {
		fmt.Print("Temp Log:", err)
		return
	}
	pretty.Print(handshake)
	pretty.Print(receivedHandshake)
	pretty.Diff(handshake, receivedHandshake)
}
