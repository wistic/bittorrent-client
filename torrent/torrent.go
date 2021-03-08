package torrent

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/IncSW/go-bencode"
)

// File hold information of files present in .torrent
type File struct {
	Path   string
	Length int64
}

// Torrent holds .torrent file data
type Torrent struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int64
	Files       []File
	Name        string
}

// String implements Stringer interface to properly print Torrent struct
func (tor Torrent) String() string {
	res := "Torrent Info\n"
	res += fmt.Sprintf("Name: %v\n", tor.Name)
	res += fmt.Sprintf("Announce: %v\n", tor.Announce)
	res += fmt.Sprintf("InfoHash: %v\n", hex.EncodeToString(tor.InfoHash[:]))
	res += fmt.Sprintf("PieceLength: %v\n", tor.PieceLength)
	res += "Files:"
	for i, file := range tor.Files {
		res += fmt.Sprintf("\tIndex:%v FilePath:%v Length:%v\n", i, file.Path, file.Length)
	}
	return res
	//return fmt.Sprintf("Torrent Info\n\tName: %s\n\tNo. of files: %v\n\tTracker URL: %v\n\tinfo-hash: %v", tor.Name, len(tor.Files), tor.Announce, hex.EncodeToString(tor.InfoHash[:]))
}

// Length calculates the total length of all files in .torrent
func (tor Torrent) Length() int64 {
	var length int64
	for _, v := range tor.Files {
		length += v.Length
	}
	return length
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

// parseFilePaths parses the file paths from files field in the infoMap
func parseFilePaths(infoMap map[string]interface{}) ([]File, error) {
	if length, ok := infoMap["length"].(int64); ok {

		return []File{{"", length}}, nil
	}

	if fileArray, ok := infoMap["files"].([]interface{}); ok {
		fileCount := len(fileArray)
		files := make([]File, fileCount)
		for i := 0; i < fileCount; i++ {
			file, ok := fileArray[i].(map[string]interface{})
			if !ok {
				return nil, errors.New("file is not a map")
			}
			files[i].Length, ok = file["length"].(int64)
			if !ok {
				return nil, errors.New("file length is corrupt")
			}
			pathArray, ok := file["path"].([]interface{})
			if !ok {
				return nil, errors.New("file path is corrupt")
			}
			path := ""
			for _, value := range pathArray {
				pathPart, ok := value.([]byte)
				if !ok {
					return nil, errors.New("some part of the file path is corrupt")
				}
				path = filepath.Join(path, string(pathPart))
			}
			files[i].Path = path
		}
		return files, nil
	}

	return nil, errors.New("file information missing")
}

// Parse parses the torrent code
func Parse(code []byte) (Torrent, error) {
	data, err := bencode.Unmarshal(code)
	if err != nil {
		return Torrent{}, err
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		err := errors.New("torrent file corrupt: dataMap is not a dictionary")
		return Torrent{}, err
	}
	infoMap, ok := dataMap["info"].(map[string]interface{})
	if !ok {
		err := errors.New("torrent file corrupt: infoMap is not a dictionary")
		return Torrent{}, err
	}

	pieceArray, ok := infoMap["pieces"].([]byte)
	if !ok {
		err := errors.New("torrent file corrupt: pieces is not an array")
		return Torrent{}, err
	}
	hashes, err := splitHashes(pieceArray)
	if err != nil {
		return Torrent{}, err
	}

	infoHash, err := infoHash(dataMap["info"])
	if err != nil {
		return Torrent{}, err
	}

	files, err := parseFilePaths(infoMap)
	if err != nil {
		return Torrent{}, err
	}

	announce, ok := dataMap["announce"].([]byte)
	if !ok {
		err := errors.New("torrent file corrupt: tracker url missing")
		return Torrent{}, err
	}
	pieceLength, ok := infoMap["piece length"].(int64)
	if !ok {
		err := errors.New("torrent file corrupt: piece length missing")
		return Torrent{}, err
	}
	name, ok := infoMap["name"].([]byte)
	if !ok {
		return Torrent{}, errors.New("torrent name is missing")
	}

	tor := Torrent{
		Announce:    string(announce),
		InfoHash:    infoHash,
		PieceHashes: hashes,
		PieceLength: pieceLength,
		Files:       files,
		Name:        string(name),
	}
	return tor, nil
}
