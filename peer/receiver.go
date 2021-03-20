package peer

import (
	"bittorrent-go/message"
	"bittorrent-go/util"
	"fmt"
	"net"
	"time"
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

func ReceiverCoroutine(address *util.Address, connection net.Conn, messageChannel chan message.Message, errorChannel chan error) {
	fmt.Println("[receiver ", address.String(), "] ", "routine started")
	defer fmt.Println("[receiver ", address.String(), "] ", "routine finished")
	defer close(messageChannel)
	defer close(errorChannel)
	for {
		err := connection.SetDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			fmt.Println("[receiver ", address.String(), "] ", "deadline error: ", err)
			errorChannel <- err
			return
		}

		packet, err := message.ReceiveMessage(connection)
		if err != nil {
			fmt.Println("[receiver ", address.String(), "] ", "parsing error: ", err)
			errorChannel <- err
			return
		}
		messageChannel <- packet
	}
}
