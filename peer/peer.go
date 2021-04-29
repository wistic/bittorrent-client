package peer

import (
	"bittorrent-go/filesystem"
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

const BLOCK_SIZE uint32 = 16 * 1024

func WorkerRoutine(ctx context.Context, wg *sync.WaitGroup, address *util.Address, peerID *util.PeerID, infoHash *util.Hash, jobs chan *job.Job) {
	wg.Add(1)
	defer wg.Done()

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
		case job := <-jobs:
			ok, err := bitfield.CheckPiece(int(job.Index))
			if err != nil || !ok {
				jobs <- job
				continue
			}

			fmt.Println("[worker ", address.String(), "] ", "job picked ", job.Index)

			data := make([]byte, job.Length)
			downloaded := uint32(0)
			requested := uint32(0)
			for downloaded < job.Length {
				if !choke {
					for requested < downloaded+BLOCK_SIZE*3 && requested < job.Length {
						blockSize := job.Length - requested
						if blockSize > BLOCK_SIZE {
							blockSize = BLOCK_SIZE
						}
						req := message.Request{
							Index:  job.Index,
							Begin:  requested,
							Length: blockSize,
						}
						err := message.SendMessage(&req, connection)
						if err != nil {
							jobs <- job
							fmt.Println("[worker ", address.String(), "] ", "send message error ", err)
							return
						}
						requested += blockSize
					}
				}

				select {
				case msg, ok := <-messageChannel:
					if !ok || msg == nil {
						jobs <- job
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
					jobs <- job
					fmt.Println("[worker ", address.String(), "] ", "context closed")
					return
				default:
					continue
				}
			}

			hash := util.Hash{
				Value: sha1.Sum(data),
			}

			fmt.Println("[worker ", address.String(), "] ", "job done: ", job.Index, "hash: ", hash, "match: ", hash.Match(&job.Hash))

			go filesystem.WriteRoutine(wg, job.Index, data)

			continue
		case <-ctx.Done():
			return
		}
	}

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
