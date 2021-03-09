package tracker

import (
	"bittorrent-go/torrent"
	"bittorrent-go/util"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// BuildTrackerURL builds the url
func buildTrackerURL(tor *torrent.Torrent, peerID *util.PeerID, port uint16) (string, error) {
	baseUrl, err := url.Parse(tor.Announce)
	if err != nil {
		return "", errors.New("announce url broken")
	}
	params := url.Values{
		"info_hash":  []string{string(tor.InfoHash.Slice()[:])}, // info-hash of the given torrent file
		"peer_id":    []string{string(peerID.Value[:])},         // peer id for this client
		"port":       []string{strconv.Itoa(int(port))},         // port number the client is listening on
		"uploaded":   []string{"0"},                             // amount uploaded so far
		"downloaded": []string{"0"},                             // amount downloaded so far
		"compact":    []string{"0"},                             // bep_0023: change compact mode
		"left":       []string{strconv.Itoa(int(tor.Length()))}, // amount of data left to be downloaded
	}
	baseUrl.RawQuery = params.Encode()
	return baseUrl.String(), nil
}

// requestTrackerHttp makes the http request to the tracker
func requestTrackerHttp(tor *torrent.Torrent, peerID *util.PeerID, port uint16) ([]byte, error) {
	trackerUrl, err := buildTrackerURL(tor, peerID, port)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(trackerUrl)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := io.ReadAll(resp.Body)
	return body, nil
}

// RequestTracker maker the request to tracker
func RequestTracker(tor *torrent.Torrent, peerID *util.PeerID, port uint16) (*Response, error) {
	responseData, err := requestTrackerHttp(tor, peerID, port)
	if err != nil {
		return nil, err
	}
	response, err := parse(responseData)
	if err != nil {
		return nil, err
	}
	return response, nil
}
