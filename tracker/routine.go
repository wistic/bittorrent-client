package tracker

import (
	"bittorrent-go/torrent"
	"bittorrent-go/util"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

func RoutineHTTP(ctx context.Context, wg *sync.WaitGroup, torrent *torrent.Torrent, peerID *util.PeerID, port uint16, channel chan *Response) {
	wg.Add(1)
	defer wg.Done()

	fmt.Println("[tracker http] routine started")
	defer fmt.Println("[tracker http] routine finished")
	for {
		fmt.Println("[tracker http] requesting")
		response, err := RequestTrackerHTTP(torrent, peerID, port)
		if err != nil {
			fmt.Println("[tracker http] request error: ", err)
		} else {
			fmt.Println("[tracker http] response received: ", response)
			channel <- response
		}

		select {
		case <-ctx.Done():
			close(channel)
			return
		case <-time.After(time.Duration(response.Interval) * time.Second):
			continue
		}
	}
}

func RoutineUDP(ctx context.Context, wg *sync.WaitGroup, torrent *torrent.Torrent, peerID *util.PeerID, port uint16, channel chan *Response) {
	wg.Add(1)
	defer wg.Done()

	fmt.Println("[tracker udp] routine started")
	defer fmt.Println("[tracker udp] routine finished")

	url := strings.TrimPrefix(torrent.Announce, "udp://")
	conn, id, err := ConnectUDP(url)
	if err != nil {
		fmt.Println("[tracker udp] connection error: ", err)
		close(channel)
		return
	}
	fmt.Println("[tracker udp] id received: ", id)
	for {
		fmt.Println("[tracker udp] requesting")
		response, err := AnnounceUDP(&torrent.InfoHash, peerID, port, conn, id)
		if err != nil {
			fmt.Println("[tracker udp] request error: ", err)
		} else {
			fmt.Println("[tracker udp] response received: ", response)
			channel <- response
		}

		select {
		case <-ctx.Done():
			close(channel)
			return
		case <-time.After(time.Duration(response.Interval) * time.Second):
			continue
		}
	}
}

func StartTrackerRoutine(ctx context.Context, wg *sync.WaitGroup, torrent *torrent.Torrent, peerID *util.PeerID, port uint16) chan *Response {
	response := make(chan *Response)
	url := torrent.Announce
	if strings.HasPrefix(url, "http://") {
		go RoutineHTTP(ctx, wg, torrent, peerID, port, response)
	} else {
		fmt.Print("udp")
		go RoutineUDP(ctx, wg, torrent, peerID, port, response)
	}
	return response
}
