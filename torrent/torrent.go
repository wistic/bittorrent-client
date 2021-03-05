package torrent

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
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

// String implements Stringer interface to properly print Torrent struct
func (tor Torrent) String() string {
	return fmt.Sprintf("Torrent of file %v of length %v with infoHash %v", tor.Name, tor.Length, hex.EncodeToString(tor.InfoHash[:]))
}

// splitHashes splits the give array of bytes to SHA1 hashes
func splitHashes(pieceArray []byte) ([][20]byte, error) {
	if len(pieceArray)%20 != 0 {
		err := errors.New("piece hash information is corrupt")
		return nil, err
	}
	numHashes := len(pieceArray) / 20
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], pieceArray[i*20:(i+1)*20])
	}
	return hashes, nil
}

// infoHash returns the 20 byte sha1 hash of the bencoded form of the info value
// for more information refer https://www.bittorrent.org/beps/bep_0003.html
func infoHash(info interface{}) ([20]byte, error) {
	data, err := bencode.Marshal(info)
	if err != nil {
		return [20]byte{}, err
	}
	return sha1.Sum(data), nil
}

// Parse parses the torrent code
func Parse(code []byte) (Torrent, error) {
	data, err := bencode.Unmarshal(code)
	if err != nil {
		return Torrent{}, err
	}

	dataMap := data.(map[string]interface{})
	infoMap := dataMap["info"].(map[string]interface{})

	hashes, err := splitHashes(infoMap["pieces"].([]byte))
	if err != nil {
		return Torrent{}, err
	}

	infoHash, err := infoHash(dataMap["info"])
	if err != nil {
		return Torrent{}, err
	}

	tor := Torrent{
		Announce:    string(dataMap["announce"].([]byte)),
		InfoHash:    infoHash,
		PieceHashes: hashes,
		PieceLength: infoMap["piece length"].(int64),
		Length:      infoMap["length"].(int64),
		Name:        string(infoMap["name"].([]byte)),
	}
	return tor, nil
}
