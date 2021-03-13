package peer

import (
	"bittorrent-go/message"
	"bittorrent-go/util"
	"fmt"
)

type Channel struct {
	Sender   chan *message.Directive
	Receiver chan *message.Directive
	Error    chan error
}

func StartSender(address string, peerID util.PeerID, infoHash util.Hash) {
	peer, err := AttemptConnection(address, &peerID, &infoHash)
	if err != nil {
		fmt.Println("Could not establish a connection with", address, ". Reason:", err)
		return
	}
	defer peer.Connection.Close()
	fmt.Println("Connected to", address)
}
