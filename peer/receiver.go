package peer

import (
	"bittorrent-go/message"
	"bittorrent-go/util"
	"context"
	"fmt"
	"net"
	"os"
	"time"
)

func ReceiverRoutine(ctx context.Context, address *util.Address, connection net.Conn, messageChannel chan message.Message) {
	fmt.Println("[receiver ", address.String(), "] ", "routine started")
	defer fmt.Println("[receiver ", address.String(), "] ", "routine finished")

	defer close(messageChannel)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		err := connection.SetReadDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			fmt.Println("[receiver ", address.String(), "] ", "deadline error: ", err)
			return
		}

		packet, err := message.ReceiveMessage(connection)
		if err != nil {
			fmt.Println("[receiver ", address.String(), "] ", "parsing error: ", err)
			if os.IsTimeout(err) {
				continue
			}
			return
		}

		fmt.Println("[receiver ", address.String(), "] ", "packet received")
		messageChannel <- packet
	}
}
