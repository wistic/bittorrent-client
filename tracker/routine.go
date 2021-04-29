package tracker

import (
	"bittorrent-go/torrent"
	"bittorrent-go/util"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

func RoutineHTTP(ctx context.Context, wg *sync.WaitGroup, torrent *torrent.Torrent, peerID *util.PeerID, port uint16, channel chan *Response) {
	wg.Add(1)
	defer wg.Done()

	logrus.Debugln("tracker started")
	defer logrus.Debugln("tracker finished")

	for {
		logrus.Debugln("tracker requesting")
		response, err := RequestTrackerHTTP(torrent, peerID, port)
		if err != nil {
			logrus.Debugln("tracker request error")

		} else {
			logrus.Debugln("tracker response received")
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

	logrus.Debugln("tracker started")
	defer logrus.Debugln("tracker finished")

	url := strings.TrimPrefix(torrent.Announce, "udp://")
	conn, id, err := ConnectUDP(url)
	if err != nil {
		logrus.Debugln("tracker udp connection error")
		close(channel)
		return
	}
	logrus.Debugln("tracker connection id:", id)

	for {
		logrus.Debugln("tracker requesting")

		response, err := AnnounceUDP(&torrent.InfoHash, peerID, port, conn, id)
		if err != nil {
			logrus.Debugln("tracker request error:", err)

		} else {
			logrus.Debugln("tracker response received")
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
