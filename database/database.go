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
	ConnectionInfo *util.Address
}

type Database struct {
	Interval   uint64
	Torrent    *torrent.Torrent
	Peers      []PeerData
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
	peerArray := make([]PeerData, len(response.Peers))
	for i, v := range response.Peers {
		peerArray[i] = PeerData{
			AmChoking:      true,
			AmInterested:   false,
			PeerChoking:    true,
			PeerInterested: false,
			BitField:       util.BitField{},
			ConnectionInfo: &v,
		}
	}
}
