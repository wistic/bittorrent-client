package main

import (
	"bittorrent-go/cli"
	"bittorrent-go/torrent"
	"fmt"
	"io/ioutil"
)

func main() {
	args, err := cli.Parse()
	if err != nil {
		fmt.Println("Argument parsing error: ", err)
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

	fmt.Println(tor)
}
