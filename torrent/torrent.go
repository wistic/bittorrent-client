package torrent

import (
	"fmt"
	"github.com/IncSW/go-bencode"
)

// Torrent holds .torrent file data
type Torrent struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int64
	Length      int64
	Name        string
}

// Parse parses the torrent code
func Parse(code []byte) (Torrent, error) {
	data, err := bencode.Unmarshal(code)
	if err != nil {
		return Torrent{}, err
	}
	dataMap := data.(map[string]interface{})
	infoMap := dataMap["info"].(map[string]interface{})
	fmt.Println(infoMap)
	tor := Torrent{
		Announce:    string(dataMap["announce"].([]byte)),
		InfoHash:    [20]byte{},
		PieceHashes: nil,
		PieceLength: infoMap["piece length"].(int64),
		Length:      infoMap["length"].(int64),
		Name:        string(infoMap["name"].([]byte)),
	}
	return tor, nil
}
