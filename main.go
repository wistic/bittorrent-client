package main

import (
	"bittorrent-go/cli"
	"bittorrent-go/peer/attribute"
	"bittorrent-go/torrent"
	"bittorrent-go/tracker"
	"fmt"
	"io/ioutil"

	"github.com/kr/pretty"
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

	peerID := attribute.GeneratePeerID()

	resp, err := tracker.RequestTracker(tor, peerID, 9969)
	if err != nil {
		fmt.Println("Tracker request error:", err)
		return
	}
	response, err := tracker.Parse(resp)
	if err != nil {
		fmt.Println("Tracker response error:", err)
		return
	}
	_, _ = pretty.Println(response)
}
