package peer

import (
	"bittorrent-go/message"
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

func ReceiverCoroutine(connection net.Conn, messageChannel chan message.Message, errorChannel chan error) {
	fmt.Println("start receive coroutine")
	defer fmt.Println("stop receive coroutine")
	defer close(messageChannel)
	defer close(errorChannel)
	for {
		err := connection.SetDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			fmt.Println("receiver coroutine error (deadline): ", err)
			errorChannel <- err
			return
		}

		packet, err := message.ReceiveMessage(connection)
		if err != nil {
			fmt.Println("receiver coroutine error (reading) ", err)
			errorChannel <- err
			return
		}
		messageChannel <- packet
	}
}
