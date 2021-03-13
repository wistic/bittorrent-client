package database

import (
	"bittorrent-go/torrent"
	"bittorrent-go/tracker"
	"bittorrent-go/util"
)

type PeerData struct {
	AmChoking      bool
	AmInterested   bool
	PeerChoking    bool
	PeerInterested bool
	BitField       util.BitField
}

// TODO: Add piece progress measure

type Database struct {
	Interval   uint64
	Torrent    *torrent.Torrent
	Peers      map[string]PeerData
	MyPeerID   util.PeerID
	MyBitField util.BitField
	Port       uint16
	OutputPath string
}

func (database *Database) fill(response *tracker.Response, torrent *torrent.Torrent, peerID util.PeerID, path string) {
	database.Interval = response.Interval
	database.Torrent = torrent
	database.MyPeerID = peerID
	database.MyBitField = util.BitField{}
	database.OutputPath = path
	peerArray := make(map[string]PeerData)
	for _, v := range response.Peers {
		peerArray[v.String()] = PeerData{
			AmChoking:      true,
			AmInterested:   false,
			PeerChoking:    true,
			PeerInterested: false,
			BitField:       util.BitField{},
		}
	}
}
