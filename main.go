package main

import (
	"bittorrent-go/cli"
	"bittorrent-go/tracker"
	"bittorrent-go/util"
	"context"
	"github.com/kr/pretty"
	"sync"

	"bittorrent-go/torrent"
	//"bittorrent-go/tracker"
	//"bittorrent-go/util"
	"fmt"
	//"time"
)

func main() {
	args, err := cli.Parse()
	if err != nil {
		fmt.Println("[cli] argument parsing error: ", err)
		fmt.Println("[cli] usage: bittorrent-go -o <path-to-download-directory> <path-to-torrent-file> ")
		return
	}

	tor, err := torrent.FromFile(args.Torrent)
	if err != nil {
		fmt.Println("[parser] torrent parsing error: ", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	//jobs := job.CreateJobQueue(tor)
	//
	//for j := range jobs {
	//	fmt.Println(j.Index, j.Length)
	//}

	peerID := util.GeneratePeerID()

	responses := tracker.StartTrackerRoutine(ctx, &wg, tor, peerID, 9969)
	response := <-responses

	pretty.Println(response)
	//
	//go peer.WorkerRoutine(&response.Peers[0], peerID, &tor.InfoHash, &sch)
	//for {
	//	select {
	//	case <-time.After(time.Second * 60):
	//		filesystem.StopWriter(writerDoneChannel)
	//		return
	//	}
	//}

	cancel()
	wg.Wait()
}
