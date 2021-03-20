package peer

import (
	"bittorrent-go/message"
	"bittorrent-go/scheduler"
	"bittorrent-go/util"
	"bytes"
	"errors"
	"fmt"
	"github.com/kr/pretty"
	"net"
	"time"
)

const protocolIdentifier = "BitTorrent protocol"

type Peer struct {
	Connection net.Conn
	PeerID     util.PeerID
	Address    string
}

func AttemptConnection(address string, peerID *util.PeerID, infoHash *util.Hash) (*Peer, error) {
	connection, err := net.DialTimeout("tcp", address, 6*time.Second)
	if err != nil {
		return nil, err
	}

	connection.SetDeadline(time.Now().Add(5 * time.Second))
	defer connection.SetDeadline(time.Time{})

	handshake := message.Handshake{Protocol: protocolIdentifier, Extension: util.Extension{}, InfoHash: *infoHash, PeerID: *peerID}
	err = message.WriteHandshake(&handshake, connection)
	if err != nil {
		_ = connection.Close()
		return nil, err
	}
	receivedHandshake, err := message.ReadHandshake(connection)
	if err != nil {
		_ = connection.Close()
		return nil, err
	}
	if receivedHandshake.Protocol != protocolIdentifier || !bytes.Equal(receivedHandshake.InfoHash.Slice(), infoHash.Slice()) {
		connection.Close()
		return nil, errors.New("bad handshake")
	}
	peer := Peer{
		Connection: connection,
		PeerID:     receivedHandshake.PeerID,
		Address:    address,
	}
	return &peer, nil
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

func PeerRoutine(address *util.Address, peerID *util.PeerID, infoHash *util.Hash, scheduler *scheduler.Scheduler) {
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
	errorChannel := make(chan error, 10)
	go ReceiverCoroutine(address, connection, messageChannel, errorChannel)
	for {
		piece, _ := scheduler.GetPiece(address)
		counter := 0
		fmt.Println("[worker ", address.String(), "] ", "piece assigned: ", piece)

		for {
			select {
			case msg, ok := <-messageChannel:
				if ok {
					switch v := msg.(type) {
					case *message.Piece:
						fmt.Println("[worker ", address.String(), "] ", "message received: piece: ", v.Index, " begin: ", v.Begin, " length: ", len(v.Block))
					default:
						pretty.Println("[worker ", address.String(), "] ", "message received:", v)
					}
				} else {
					fmt.Println("[worker ", address.String(), "] ", "message channel closed")
					return
				}
			case err, ok := <-errorChannel:
				if ok {
					fmt.Println("[worker ", address.String(), "] ", "error received: ", err)
				} else {
					fmt.Println("[worker ", address.String(), "] ", "error channel closed")
				}
			case <-time.After(time.Duration(5) * time.Second):
				mess := message.Request{
					Index:  piece,
					Begin:  uint32(counter * 16384),
					Length: 16384,
				}
				err := message.SendMessage(&mess, connection)
				pretty.Println("[worker ", address.String(), "] ", "sending message: ", mess)
				if err != nil {
					fmt.Println("[worker ", address.String(), "] ", "message sending error: ", err)
				}
				counter++
			}
		}
	}
}
