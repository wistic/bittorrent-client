package peer

import (
	"bittorrent-go/message"
	"bittorrent-go/util"
	"context"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"time"
)

func ReceiverRoutine(ctx context.Context, address *util.Address, connection net.Conn, messageChannel chan message.Message) {
	logrus.Debugln("receiver", address, "started")
	defer logrus.Debugln("receiver", address, "ended")

	defer close(messageChannel)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		err := connection.SetReadDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			logrus.Debugln("receiver", address, "deadline error:", err)
			return
		}

		packet, err := message.ReceiveMessage(connection)
		if err != nil {
			logrus.Debugln("receiver", address, "parsing error:", err)
			if os.IsTimeout(err) {
				continue
			}
			return
		}

		logrus.Traceln("receiver", address, "packet received")
		messageChannel <- packet
	}
}
