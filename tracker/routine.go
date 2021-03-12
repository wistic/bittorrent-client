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
	for {
		fmt.Println("requesting")
		response, err := RequestTrackerHTTP(torrent, peerID, port)
		if err != nil {
			channel.Error <- err
		} else {
			channel.Response <- response
		}

		select {
		case <-channel.Done:
			close(channel.Response)
			close(channel.Error)
			return
		case <-time.After(time.Duration(2) * time.Second): //TODO: Replace duration
			continue
		}
	}
}

func RoutineUDP(torrent *torrent.Torrent, peerID *util.PeerID, port uint16, channel RoutineChannel) {
	url := strings.TrimPrefix(torrent.Announce, "udp://")
	conn, id, err := ConnectUDP(url)
	if err != nil {
		channel.Error <- err
		close(channel.Response)
		close(channel.Error)
		return
	}
	fmt.Println("connection id", id)
	for {
		fmt.Println("requesting")
		response, err := AnnounceUDP(&torrent.InfoHash, peerID, port, conn, id)
		if err != nil {
			channel.Error <- err
			fmt.Println(err)

		} else {
			channel.Response <- response
			fmt.Println("response")
		}

		select {
		case <-channel.Done:
			close(channel.Response)
			close(channel.Error)
			return
		case <-time.After(time.Duration(2) * time.Second): //TODO: Replace duration
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
