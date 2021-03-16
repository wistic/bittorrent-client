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
	fmt.Println("start worker")
	defer fmt.Println("end worker")
	connection, err := net.DialTimeout("tcp", address.String(), 6*time.Second)
	if err != nil {
		fmt.Println("timeout set", err)
		return
	}

	err = HandshakeRoutine(connection, &message.Handshake{Protocol: protocolIdentifier, InfoHash: *infoHash, PeerID: *peerID, Extension: util.Extension{}})
	if err != nil {
		fmt.Println("handshake ", err)
		return
	}
	fmt.Println("handshake done")

	messageChannel := make(chan message.Message, 10)
	errorChannel := make(chan error, 10)
	go ReceiverCoroutine(connection, messageChannel, errorChannel)
	for {
		piece, _ := scheduler.GetPiece(address)
		counter := 0
		fmt.Println("piece assigned: ", piece)
		for {
			select {
			case msg, ok := <-messageChannel:
				if ok {
					switch v := msg.(type) {
					case *message.Piece:
						fmt.Println("message channel receive: Piece:", v.Index, " Begin:", v.Begin, " Length:", len(v.Block))
					default:
						pretty.Println("message channel receive:", v)
					}
				} else {
					fmt.Println("message channel closed")
					return
				}
			case err, ok := <-errorChannel:
				if ok {
					pretty.Println("error channel received: ", err)
				} else {
					fmt.Println("error channel closed")
				}
			case <-time.After(time.Duration(5) * time.Second):
				mess := message.Request{
					Index:  piece,
					Begin:  uint32(counter * 16384),
					Length: 16384,
				}
				err := message.SendMessage(&mess, connection)
				pretty.Println("message sent:", mess)
				if err != nil {
					fmt.Println("message sending error :", err)
				}
				counter++
			}
		}
	}
}
