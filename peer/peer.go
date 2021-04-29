package peer

import (
	"bittorrent-go/job"
	"bittorrent-go/message"
	"bittorrent-go/util"
	"bytes"
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/kr/pretty"
	"net"
	"sync"
	"time"
)

const protocolIdentifier = "BitTorrent protocol"

type Peer struct {
	Connection net.Conn
	PeerID     util.PeerID
	Address    string
}

func HandshakeRoutine(connection net.Conn, peerID *util.PeerID, infoHash *util.Hash) error {
	handshake := message.Handshake{
		Protocol:  protocolIdentifier,
		InfoHash:  *infoHash,
		PeerID:    *peerID,
		Extension: util.Extension{},
	}
	err := connection.SetReadDeadline(time.Now().Add(5 * time.Second))

	if err != nil {
		return err
	}
	err = message.WriteHandshake(&handshake, connection)
	if err != nil {
		return err
	}
	rec, err := message.ReadHandshake(connection)
	if err != nil {
		return err
	}
	if rec.Protocol != handshake.Protocol || !bytes.Equal(rec.InfoHash.Slice(), handshake.InfoHash.Slice()) {
		return errors.New("bad handshake")
	}
	return nil
}

func BitfieldRoutine(connection net.Conn) (*util.BitField, error) {
	err := connection.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return nil, err
	}

	packet, err := message.ReceiveMessage(connection)
	if err != nil {
		return nil, err
	}

	if packet.GetMessageID() != message.MsgBitfield {
		return nil, errors.New("wrong message type")
	}
	field := packet.(*message.BitField).Field
	return &field, nil
}

const BlockSize uint32 = 16 * 1024

func Disconnect(address *util.Address, disconnect chan *util.Address) {
	disconnect <- address
}

func WorkerRoutine(ctx context.Context, wg *sync.WaitGroup, address *util.Address, peerID *util.PeerID, infoHash *util.Hash, jobs chan *job.Job, results chan *job.Result, disconnect chan *util.Address) {
	wg.Add(1)
	defer wg.Done()

	defer Disconnect(address, disconnect)

	fmt.Println("[worker ", address.String(), "] ", "routine started")
	defer fmt.Println("[worker ", address.String(), "] ", "routine finished")

	connection, err := net.DialTimeout("tcp", address.String(), 10*time.Second)
	if err != nil {
		fmt.Println("[worker ", address.String(), "] ", "connection error: ", err)
		return
	}
	defer connection.Close()

	fmt.Println("[worker ", address.String(), "] ", "connection established")

	err = HandshakeRoutine(connection, peerID, infoHash)
	if err != nil {
		fmt.Println("[worker ", address.String(), "] ", "handshake error: ", err)
		return
	}
	fmt.Println("[worker ", address.String(), "] ", "handshake done")

	bitfield, err := BitfieldRoutine(connection)
	if err != nil {
		fmt.Println("[worker ", address.String(), "] ", "bitfield error: ", err)
		return
	}

	pretty.Println(bitfield)

	messageChannel := make(chan message.Message, 10)
	go ReceiverRoutine(ctx, wg, address, connection, messageChannel)

	unchoke := message.Unchoke{}
	err = message.SendMessage(&unchoke, connection)
	if err != nil {
		fmt.Println("[worker ", address.String(), "] ", "unchoke send failed: ", err)
		return
	}
	fmt.Println("[worker ", address.String(), "] ", "choke done")

	interested := message.Interested{}
	err = message.SendMessage(&interested, connection)
	if err != nil {
		fmt.Println("[worker ", address.String(), "] ", "interested send failed: ", err)
		return
	}
	fmt.Println("[worker ", address.String(), "] ", "interested done")

	choke := true

	for {
		select {
		case j := <-jobs:
			ok, err := bitfield.CheckPiece(int(j.Index))
			if err != nil || !ok {
				jobs <- j
				continue
			}

			fmt.Println("[worker ", address.String(), "] ", "job picked ", j.Index)

			data := make([]byte, j.Length)
			downloaded := uint32(0)
			requested := uint32(0)
			for downloaded < j.Length {
				if !choke {
					for requested < downloaded+BlockSize*3 && requested < j.Length {
						blockSize := j.Length - requested
						if blockSize > BlockSize {
							blockSize = BlockSize
						}
						req := message.Request{
							Index:  j.Index,
							Begin:  requested,
							Length: blockSize,
						}
						err := message.SendMessage(&req, connection)
						if err != nil {
							jobs <- j
							fmt.Println("[worker ", address.String(), "] ", "send message error ", err)
							return
						}
						requested += blockSize
					}
				}

				select {
				case msg, ok := <-messageChannel:
					if !ok || msg == nil {
						jobs <- j
						fmt.Println("[worker ", address.String(), "] ", "message channel closed")
						return
					}
					switch msg.GetMessageID() {
					case message.MsgUnchoke:
						choke = false
						fmt.Println("[worker ", address.String(), "] ", "unchoked")
					case message.MsgChoke:
						choke = true
						fmt.Println("[worker ", address.String(), "] ", "choked")
					case message.MsgPiece:
						piece := msg.(*message.Piece)
						end := piece.Begin + uint32(len(piece.Block))
						if downloaded < end {
							downloaded = end
							fmt.Println("[worker ", address.String(), "] ", "downloaded ", downloaded)
						}
						copy(data[piece.Begin:], piece.Block)
					}
				case <-ctx.Done():
					jobs <- j
					fmt.Println("[worker ", address.String(), "] ", "context closed")
					return
				default:
					continue
				}
			}

			hash := util.Hash{
				Value: sha1.Sum(data),
			}

			if !hash.Match(&j.Hash) {
				jobs <- j
				continue
			}

			fmt.Println("[worker ", address.String(), "] ", "job done: ", j.Index, "hash: ", hash, "match: ", hash.Match(&j.Hash))

			//go filesystem.WriteRoutine(wg, j.Index, data)
			results <- &job.Result{
				Index: j.Index,
				Data:  data,
			}
			continue
		case <-ctx.Done():
			return
		}
	}
}

func CreateDisconnectChannel() chan *util.Address {
	return make(chan *util.Address)
}
