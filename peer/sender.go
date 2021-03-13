package peer

import (
	"bittorrent-go/message"
	"bittorrent-go/util"
	"fmt"
	"github.com/kr/pretty"
)

type Channel struct {
	Sender   SenderChannel
	Receiver ReceiverChannel
}

type SenderChannel struct {
	Data  chan *message.Directive
	Error chan *message.ErrorMessage
}

type ReceiverChannel struct {
	Data  chan *message.Directive
	Error chan *message.ErrorMessage
}

func StartSender(address string, peerID util.PeerID, infoHash util.Hash, channel Channel) {
	peer, err := AttemptConnection(address, &peerID, &infoHash)
	if err != nil {
		channel.Sender.Error <- &message.ErrorMessage{Address: address, Value: err}
		return
	}
	defer peer.Connection.Close()
	fmt.Println("Connected to", address)

	go StartReceiver(peer.Connection, peer.Address, channel.Receiver)

	for {
		packet, ok := <-channel.Sender.Data
		pretty.Println(packet)
		if !ok {
			peer.Connection.Close()
			return
		}
		err := message.SendMessage(packet.Message, peer.Connection)
		if err != nil {
			channel.Sender.Error <- &message.ErrorMessage{Address: address, Value: err}
			return
		}
	}
}
