package torrent

import (
	"fmt"
	"github.com/IncSW/go-bencode"
)

// Torrent holds .torrent file data
type Torrent struct {
}

// Parse parses the torrent code
func Parse(code string) (Torrent, error) {
	data, err := bencode.Unmarshal([]byte(code))
	if err != nil {
		return Torrent{}, err
	}
	fmt.Println(data)
	return Torrent{}, nil
}
