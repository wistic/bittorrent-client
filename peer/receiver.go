package peer

import (
	"bittorrent-go/message"
	"bittorrent-go/util"
	"context"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

func ReceiverRoutine(ctx context.Context, wg *sync.WaitGroup, address *util.Address, connection net.Conn, messageChannel chan message.Message) {
	wg.Add(1)
	defer wg.Done()

	fmt.Println("[receiver ", address.String(), "] ", "routine started")
	defer fmt.Println("[receiver ", address.String(), "] ", "routine finished")

	defer close(messageChannel)

	for {
		err := connection.SetReadDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			fmt.Println("[receiver ", address.String(), "] ", "deadline error: ", err)
			return
		}

		packet, err := message.ReceiveMessage(connection)
		if err != nil {
			fmt.Println("[receiver ", address.String(), "] ", "parsing error: ", err)
			if !os.IsTimeout(err) {
				continue
			}
			return
		}

		fmt.Println("[receiver ", address.String(), "] ", "packet received")
		messageChannel <- packet
	}
}
