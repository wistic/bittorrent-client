package main

import (
	"bittorrent-go/cli"
	"bittorrent-go/peer"
	"bittorrent-go/torrent"
	"bittorrent-go/tracker"
	"bittorrent-go/util"
	"fmt"
	"io/ioutil"
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

	peer.StartSender(response.Peers[0].String(), *peerID, tor.InfoHash)
	close(trackerChannel.Done)
}
