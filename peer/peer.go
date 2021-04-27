package peer

import (
	"bittorrent-go/job"
	"bittorrent-go/message"
	"bittorrent-go/util"
	"bytes"
	"errors"
	"fmt"
	"net"
	"time"
)

const protocolIdentifier = "BitTorrent protocol"

type Peer struct {
	Connection net.Conn
	PeerID     util.PeerID
	Address    string
}

func HandshakeRoutine(connection net.Conn, handshake *message.Handshake) error {
	err := connection.SetDeadline(time.Now().Add(5 * time.Second))

	if err != nil {
		_ = connection.Close()
		return err
	}
	err = message.WriteHandshake(handshake, connection)
	if err != nil {
		_ = connection.Close()
		return err
	}
	rec, err := message.ReadHandshake(connection)
	if err != nil {
		_ = connection.Close()
		return err
	}
	if rec.Protocol != handshake.Protocol || !bytes.Equal(rec.InfoHash.Slice(), handshake.InfoHash.Slice()) {
		_ = connection.Close()
		return errors.New("bad handshake")
	}
	return nil
}

func ReceiverRoutine(address *util.Address, connection net.Conn, messageChannel chan message.Message) {
	fmt.Println("[receiver ", address.String(), "] ", "routine started")
	defer fmt.Println("[receiver ", address.String(), "] ", "routine finished")
	defer close(messageChannel)
	for {
		err := connection.SetDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			fmt.Println("[receiver ", address.String(), "] ", "deadline error: ", err)
			return
		}

		packet, err := message.ReceiveMessage(connection)
		if err != nil {
			fmt.Println("[receiver ", address.String(), "] ", "parsing error: ", err)
			return
		}
		messageChannel <- packet
	}
}

func WorkerRoutine(address *util.Address, peerID *util.PeerID, infoHash *util.Hash, queue chan *job.Job) {
	fmt.Println("[worker ", address.String(), "] ", "routine started")
	defer fmt.Println("[worker ", address.String(), "] ", "routine finished")

	connection, err := net.DialTimeout("tcp", address.String(), 6*time.Second)
	if err != nil {
		fmt.Println("[worker ", address.String(), "] ", "connection error: ", err)
		return
	}
	fmt.Println("[worker ", address.String(), "] ", "connection established")

	err = HandshakeRoutine(connection, &message.Handshake{Protocol: protocolIdentifier, InfoHash: *infoHash, PeerID: *peerID, Extension: util.Extension{}})
	if err != nil {
		fmt.Println("[worker ", address.String(), "] ", "handshake error: ", err)
		return
	}
	fmt.Println("[worker ", address.String(), "] ", "handshake done")

	messageChannel := make(chan message.Message, 10)
	go ReceiverRoutine(address, connection, messageChannel)
	//for {
	//	piece, _ := scheduler.GetPiece(address)
	//	counter := 0
	//	fmt.Println("[worker ", address.String(), "] ", "piece assigned: ", piece)
	//
	//	for {
	//		select {
	//		case msg, ok := <-messageChannel:
	//			if ok {
	//				switch v := msg.(type) {
	//				case *message.Job:
	//					fmt.Println("[worker ", address.String(), "] ", "message received: piece: ", v.Index, " begin: ", v.Begin, " length: ", len(v.Block))
	//				default:
	//					pretty.Println("[worker ", address.String(), "] ", "message received:", v)
	//				}
	//			} else {
	//				fmt.Println("[worker ", address.String(), "] ", "message channel closed")
	//				return
	//			}
	//		case <-time.After(time.Duration(5) * time.Second):
	//			mess := message.Request{
	//				Index:  piece,
	//				Begin:  uint32(counter * 16384),
	//				Length: 16384,
	//			}
	//			err := message.SendMessage(&mess, connection)
	//			pretty.Println("[worker ", address.String(), "] ", "sending message: ", mess)
	//			if err != nil {
	//				fmt.Println("[worker ", address.String(), "] ", "message sending error: ", err)
	//			}
	//			counter++
	//		}
	//	}
	//}
}
