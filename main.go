package main

import (
	"bittorrent-go/cli"
	"bittorrent-go/filesystem"
	"bittorrent-go/peer"
	"bittorrent-go/scheduler"
	"bittorrent-go/torrent"
	"bittorrent-go/tracker"
	"bittorrent-go/util"
	"fmt"
	"io/ioutil"
	"time"
)

func main() {
	args, err := cli.Parse()
	if err != nil {
		fmt.Println("Argument parsing error: ", err)
		fmt.Println("Usage: bittorrent-go -o <path-to-download-directory> <path-to-torrent-file> ")
		return
	}

	content, err := ioutil.ReadFile(args.Torrent)
	if err != nil {
		fmt.Println("File reading error: ", err)
		return
	}
	tor, err := torrent.Parse(content)
	if err != nil {
		fmt.Println("Torrent parsing error: ", err)
	}

	peerID := util.GeneratePeerID()

	trackerChannel := tracker.StartTrackerRoutine(tor, peerID, 9969)

	response := <-trackerChannel.Response
	//errorChannel := make(chan *message.ErrorMessage, 10)
	//receiver := make(chan *message.Directive, 10)
	//sender := make(chan *message.Directive, 10)
	//channel := peer.Channel{
	//	Sender:   peer.SenderChannel{Data: sender, Error: errorChannel},
	//	Receiver: peer.ReceiverChannel{Data: receiver, Error: errorChannel},
	//}
	//go peer.StartSender(response.Peers[0].String(), *peerID, tor.InfoHash, channel)
	//sender <- message.NewDirective(&message.Interested{}, response.Peers[0].String())
	//counter := 0
	//for {
	//	select {
	//	case a := <-receiver:
	//		pretty.Println(a.Message)
	//	case <-time.After(time.Duration(5) * time.Second):
	//		mess := message.Request{
	//			Index:  0,
	//			Begin:  uint32(counter * 16384),
	//			Length: 16384,
	//		}
	//		sender <- message.NewDirective(&mess, response.Peers[0].String())
	//		counter++
	//	case b := <-errorChannel:
	//		fmt.Println("error", b.Value, "from", b.Address)
	//		return
	//	}
	//}
	_, writerFinishChannel, writerErrorChannel, writerDoneChannel := filesystem.StartWriter(tor, args.Output)
	sch := scheduler.Scheduler{}
	go peer.PeerRoutine(&response.Peers[0], peerID, &tor.InfoHash, &sch)
	for {
		select {
		case <-time.After(time.Second * 60):
			filesystem.StopWriter(writerDoneChannel)
			return
		case pieceIndex := <-writerFinishChannel:
			fmt.Println("finish writing piece: ", pieceIndex)
		case err := <-writerErrorChannel:
			fmt.Println("writer error: ", err)

		}
	}
}
