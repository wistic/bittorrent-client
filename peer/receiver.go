package peer

import (
	"bittorrent-go/message"
	"net"
)

func StartReceiver(connection net.Conn, address string, channel ReceiverChannel) {
	for {
		packet, err := message.ReceiveMessage(connection)
		if err != nil {
			channel.Error <- &message.ErrorMessage{Address: address, Value: err}
			return
		}
		channel.Data <- message.NewDirective(packet, address)
	}
}
