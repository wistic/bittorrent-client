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
	"github.com/sirupsen/logrus"
	"sync"
)

func main() {
	args, err := cli.Parse()
	if err != nil {
		logrus.Errorln("argument parsing error:", err)
		logrus.Infoln("usage: bittorrent-go -o <path-to-download-directory> <path-to-torrent-file>")
		return
	}

	logrus.Infoln("reading torrent file")
	tor, err := torrent.FromFile(args.Torrent)
	if err != nil {
		logrus.Errorln("torrent parsing error:", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	jobs, jobCount := job.CreateJobChannel(tor)
	results := job.CreateResultChannel()
	disconnect := peer.CreateDisconnectChannel()
	peerScheduler := peer.NewPeerScheduler()
	peerID := util.GeneratePeerID()

	logrus.Infoln("starting tracker")
	responses := tracker.StartTrackerRoutine(ctx, &wg, tor, peerID, 9969)

	for {
		connectNewPeers := false
		select {
		case response := <-responses:
			peerScheduler.UpdateList(response.Peers)
			connectNewPeers = true
			logrus.Infoln("tracker response received")

		case result := <-results:
			go filesystem.WriteRoutine(&wg, result.Index, result.Data, args.Output)
			jobCount -= 1
			logrus.Infoln("writing piece:", result.Index, "remaining pieces:", jobCount)

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
				go peer.WorkerRoutine(ctx, newAddress, peerID, &tor.InfoHash, jobs, results, disconnect)
			}
		}
	}

	logrus.Infoln("all pieces downloaded")

	cancel()
	wg.Wait()

	logrus.Infoln("assembling files")
	err = filesystem.AssembleRoutine(tor, args.Output)
	if err != nil {
		logrus.Errorln("assembling error:", err)
	}

	logrus.Infoln("deleting downloaded pieces")
	filesystem.DeleteRoutine(tor, args.Output)

	logrus.Infoln("finished")
}
