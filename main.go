package main

import (
	"bittorrent-go/cli"
	"bittorrent-go/filesystem"
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

	err = filesystem.AssembleRoutine(tor)
	if err != nil {
		fmt.Println(err)
	}
	//
	//ctx, cancel := context.WithCancel(context.Background())
	//wg := sync.WaitGroup{}
	//
	//jobs := job.CreateJobQueue(tor)
	//
	//peerID := util.GeneratePeerID()
	//
	//responses := tracker.StartTrackerRoutine(ctx, &wg, tor, peerID, 9969)
	//response := <-responses
	//
	//// Connect to first 40 peers
	//for i := 0; i < len(response.Peers); i += 1 {
	//	if i == 40 {
	//		break
	//	}
	//	go peer.WorkerRoutine(ctx, &wg, &response.Peers[i], peerID, &tor.InfoHash, jobs)
	//}
	//
	//wg.Wait()
	//cancel()
}
