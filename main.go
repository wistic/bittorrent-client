package main

import (
	"bittorrent-go/cli"
	"bittorrent-go/filesystem"
	"bittorrent-go/job"
	"bittorrent-go/peer"
	"bittorrent-go/torrent"
	"bittorrent-go/tracker"
	"bittorrent-go/util"
	"context"
	"sync"

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

	jobs, jobCount := job.CreateJobChannel(tor)
	results := job.CreateResultChannel()
	disconnect := peer.CreateDisconnectChannel()

	peerID := util.GeneratePeerID()

	responses := tracker.StartTrackerRoutine(ctx, &wg, tor, peerID, 9969)

	peerScheduler := peer.NewPeerScheduler()

	// Connect to first 40 peers
	//for i := 0; i < len(response.Peers); i += 1 {
	//	if i == 40 {
	//		break
	//	}
	//	go peer.WorkerRoutine(ctx, &wg, &response.Peers[i], peerID, &tor.InfoHash, jobs, results)
	//}

	for {
		connectNewPeers := false
		select {
		case response := <-responses:
			peerScheduler.UpdateList(response.Peers)
			connectNewPeers = true

		case result := <-results:
			go filesystem.WriteRoutine(&wg, result.Index, result.Data)
			jobCount -= 1
			fmt.Println(jobCount)

		case address := <-disconnect:
			peerScheduler.RemoveAddress(*address)
			connectNewPeers = true

		default:
		}

		if jobCount == 0 {
			break
		}

		if connectNewPeers {
			for peerScheduler.Total() < 40 {
				newAddress := peerScheduler.Next()
				if newAddress == nil {
					break
				}
				peerScheduler.AddAddress(*newAddress)
				go peer.WorkerRoutine(ctx, &wg, newAddress, peerID, &tor.InfoHash, jobs, results, disconnect)
			}
		}
	}

	cancel()
	wg.Wait()

	err = filesystem.AssembleRoutine(tor)
	if err != nil {
		fmt.Println(err)
	}
}
