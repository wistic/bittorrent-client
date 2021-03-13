package peer

import (
	"bittorrent-go/util"
	"fmt"
)

func StartWorker(address *util.Address, peerID *util.PeerID, infoHash *util.Hash) {
	peer, err := AttemptConnection(address, peerID, infoHash)
	if err != nil {
		fmt.Println("Could not establish a connection with", address.String(), ". Reason:", err)
		return
	}
	defer peer.Connection.Close()
	fmt.Println("Connected to", address.String())
}
