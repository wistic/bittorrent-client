package peer

import (
	"bittorrent-go/job"
	"bittorrent-go/message"
	"bittorrent-go/util"
	"bytes"
	"context"
	"crypto/sha1"
	"errors"
	"github.com/sirupsen/logrus"
	"net"
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

func CloseConnection(connection net.Conn) {
	err := connection.Close()
	if err != nil {
		logrus.Debugln("close error:", err)
	}
}

func WorkerRoutine(ctx context.Context, address *util.Address, peerID *util.PeerID, infoHash *util.Hash, jobs chan *job.Job, results chan *job.Result, disconnect chan *util.Address) {
	defer Disconnect(address, disconnect)

	logrus.Debugln("worker", address, "started")
	defer logrus.Debugln("worker", address, "finished")

	connection, err := net.DialTimeout("tcp", address.String(), 10*time.Second)
	if err != nil {
		logrus.Debugln("worker", address, "dial error:", err)
		return
	}
	defer CloseConnection(connection)

	logrus.Debugln("worker", address, "connection established")

	err = HandshakeRoutine(connection, peerID, infoHash)
	if err != nil {
		logrus.Debugln("worker", address, "handshake error:", err)
		return
	}

	logrus.Debugln("worker", address, "handshake done")

	bitfield, err := BitfieldRoutine(connection)
	if err != nil {
		logrus.Debugln("worker", address, "bitfield error:", err)
		return
	}

	messageChannel := make(chan message.Message, 10)
	go ReceiverRoutine(ctx, address, connection, messageChannel)

	unchoke := message.Unchoke{}
	err = message.SendMessage(&unchoke, connection)
	if err != nil {
		logrus.Debugln("worker", address, "unchoke send error:", err)
		return
	}
	logrus.Debugln("worker", address, "unchoke done")

	interested := message.Interested{}
	err = message.SendMessage(&interested, connection)
	if err != nil {
		logrus.Debugln("worker", address, "interested send error:", err)
		return
	}
	logrus.Debugln("worker", address, "interested done")

	logrus.Infoln("peer connected:", address)
	defer logrus.Infoln("peer disconnected:", address)

	choke := true

	for {
		select {
		case j := <-jobs:
			{
				ok, err := bitfield.CheckPiece(int(j.Index))
				if err != nil || !ok {
					jobs <- j
					continue
				}

				logrus.Debugln("worker", address, "job picked:", j.Index)

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
								logrus.Debugln("worker", address, "request send error:", err)
								return
							}
							requested += blockSize
						}
					}

					select {
					case msg, ok := <-messageChannel:
						if !ok || msg == nil {
							jobs <- j
							logrus.Debugln("worker", address, "message channel closed")
							return
						}

						switch msg.GetMessageID() {
						case message.MsgUnchoke:
							choke = false
							logrus.Debugln("worker", address, "unchoked")

						case message.MsgChoke:
							choke = true
							logrus.Debugln("worker", address, "choked")

						case message.MsgHave:
							have := msg.(*message.Have)
							err := bitfield.SetPiece(int(have.Index))
							if err != nil {
								jobs <- j
								logrus.Debugln("worker", address, "bitfield set error")
								return
							}

						case message.MsgPiece:
							piece := msg.(*message.Piece)
							end := piece.Begin + uint32(len(piece.Block))
							if downloaded < end {
								downloaded = end
								logrus.Debugln("worker", address, "piece message")
							}
							copy(data[piece.Begin:], piece.Block)
						}
					case <-ctx.Done():
						jobs <- j
						logrus.Debugln("worker", address, "context closed")
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
					logrus.Debugln("worker", address, "hash mismatch")
					continue
				}

				logrus.Debugln("worker", address, "job done:", j.Index)

				results <- &job.Result{
					Index: j.Index,
					Data:  data,
				}
			}

		case <-ctx.Done():
			return
		}
	}
}

func CreateDisconnectChannel() chan *util.Address {
	return make(chan *util.Address)
}
