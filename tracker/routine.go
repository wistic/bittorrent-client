package tracker

import (
	"bittorrent-go/torrent"
	"bittorrent-go/util"
	"fmt"
	"strings"
	"time"
)

type RoutineChannel struct {
	Response chan *Response
	Error    chan error
	Done     chan struct{}
}

func RoutineHTTP(torrent *torrent.Torrent, peerID *util.PeerID, port uint16, channel RoutineChannel) {
	fmt.Println("[tracker http] routine started")
	defer fmt.Println("[tracker http] routine finished")
	for {
		fmt.Println("[tracker http] requesting")
		response, err := RequestTrackerHTTP(torrent, peerID, port)
		if err != nil {
			fmt.Println("[tracker http] request error: ", err)
			channel.Error <- err
		} else {
			fmt.Println("[tracker http] response received: ", response)
			channel.Response <- response
		}

		select {
		case <-channel.Done:
			close(channel.Response)
			close(channel.Error)
			return
		case <-time.After(time.Duration(response.Interval) * time.Second):
			continue
		}
	}
}

func RoutineUDP(torrent *torrent.Torrent, peerID *util.PeerID, port uint16, channel RoutineChannel) {
	fmt.Println("[tracker udp] routine started")
	defer fmt.Println("[tracker udp] routine finished")

	url := strings.TrimPrefix(torrent.Announce, "udp://")
	conn, id, err := ConnectUDP(url)
	if err != nil {
		fmt.Println("[tracker udp] connection error: ", err)
		channel.Error <- err
		close(channel.Response)
		close(channel.Error)
		return
	}
	fmt.Println("[tracker udp] id received: ", id)
	for {
		fmt.Println("[tracker udp] requesting")
		response, err := AnnounceUDP(&torrent.InfoHash, peerID, port, conn, id)
		if err != nil {
			fmt.Println("[tracker udp] request error: ", err)
			channel.Error <- err
		} else {
			fmt.Println("[tracker udp] response received: ", response)
			channel.Response <- response
		}

		select {
		case <-channel.Done:
			close(channel.Response)
			close(channel.Error)
			return
		case <-time.After(time.Duration(response.Interval) * time.Second):
			continue
		}
	}
}

func StartTrackerRoutine(torrent *torrent.Torrent, peerID *util.PeerID, port uint16) RoutineChannel {
	channel := RoutineChannel{
		Response: make(chan *Response, 10),
		Error:    make(chan error, 10),
		Done:     make(chan struct{}),
	}
	url := torrent.Announce
	if strings.HasPrefix(url, "http://") {
		go RoutineHTTP(torrent, peerID, port, channel)
	} else {
		fmt.Print("udp")
		go RoutineUDP(torrent, peerID, port, channel)
	}
	return channel
}
