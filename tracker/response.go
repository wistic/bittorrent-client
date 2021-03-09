package tracker

import (
	"bittorrent-go/util"
	"encoding/binary"
	"errors"

	"github.com/IncSW/go-bencode"
)

// Response represents the response sent by the tracker
type Response struct {
	Interval int64
	Peers    []util.Address
}

// parseCompactPeerArray parses the compact peerArray
func parseCompactPeerArray(peerArray []byte) ([]util.Address, error) {
	const peerSize = 6 // bep_0023: 4 for ip, 2 for port
	peerCount := len(peerArray) / peerSize
	if len(peerArray)%peerSize != 0 {
		return nil, errors.New("peers string is corrupt")
	}
	peers := make([]util.Address, peerCount)
	for i := 0; i < peerCount; i++ {
		offset := i * peerSize
		peers[i] = util.Address{IP: peerArray[offset : offset+4], Port: binary.BigEndian.Uint16(peerArray[offset+4 : offset+6])}
	}
	return peers, nil
}

// parse parses the response received from the tracker
func parse(resp []byte) (*Response, error) {
	data, err := bencode.Unmarshal(resp)
	if err != nil {
		return nil, err
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("bad response from the tracker")
	}
	failure, ok := dataMap["failure reason"].([]byte)
	if ok {
		return nil, errors.New("tracker rejected request with reason '" + string(failure) + "'")
	}
	interval, ok := dataMap["interval"].(int64)
	if !ok {
		return nil, errors.New("invalid tracker interval")
	}
	peerArray, ok := dataMap["peers"].([]byte)
	if !ok {
		return nil, errors.New("list of peers is corrupt")
	}
	peers, err := parseCompactPeerArray(peerArray)
	if err != nil {
		return nil, err
	}
	response := Response{
		Interval: interval,
		Peers:    peers,
	}
	return &response, nil
}
