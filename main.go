package main

import (
	"bittorrent-go/cli"
	"bittorrent-go/torrent"
	"bittorrent-go/tracker"
	"fmt"
	"io/ioutil"
)

func main() {
	args, err := cli.Parse()
	if err != nil {
		fmt.Println("Argument parsing error: ", err)
		fmt.Println("Usage: bittorrent-go <path-to-torrent-file> <path-to-download-directory>")
		return
	}
	fmt.Println(args)

	content, err := ioutil.ReadFile(args.Torrent)
	if err != nil {
		fmt.Println("File reading error: ", err)
		return
	}

	tor, err := torrent.Parse(content)
	if err != nil {
		fmt.Println("Torrent parsing error: ", err)
	}

	resp, err := tracker.RequestTracker(tor, [20]byte{}, 9969)
	if err != nil {
		fmt.Println("Tracker request error:", err)
		return
	}
	response, err := tracker.Parse(resp)
	if err != nil {
		fmt.Println("Tracker response error:", err)
		return
	}
	fmt.Println(response)
}
