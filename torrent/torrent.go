package torrent

import (
	"bittorrent-go/util"
	"crypto/sha1"
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/IncSW/go-bencode"
)

// File store file data
type File struct {
	Path   string
	Length int64
}

// Torrent stores .torrent file data
type Torrent struct {
	Announce    string
	InfoHash    util.Hash
	PieceHashes []util.Hash
	PieceLength int64
	Files       []File
	Name        string
}

// Length calculates the total length of all files in .torrent
func (tor *Torrent) Length() int64 {
	var length int64
	for _, v := range tor.Files {
		length += v.Length
	}
	return length
}

// splitHashes splits the give array of bytes to SHA1 hashes
func splitHashes(pieceArray []byte) ([]util.Hash, error) {
	if len(pieceArray)%20 != 0 {
		err := errors.New("piece hash information is corrupt")
		return nil, err
	}
	numHashes := len(pieceArray) / 20
	hashes := make([]util.Hash, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i].Value[:], pieceArray[i*20:(i+1)*20])
	}
	return hashes, nil
}

// infoHash returns the 20 byte sha1 hash of the bencoded form of the info value
// for more information refer https://www.bittorrent.org/beps/bep_0003.html
func infoHash(info interface{}) (*util.Hash, error) {
	data, err := bencode.Marshal(info)
	if err != nil {
		return nil, err
	}
	return &util.Hash{Value: sha1.Sum(data)}, nil
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
func Parse(code []byte) (*Torrent, error) {
	data, err := bencode.Unmarshal(code)
	if err != nil {
		return nil, err
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		err := errors.New("torrent file corrupt: dataMap is not a dictionary")
		return nil, err
	}
	infoMap, ok := dataMap["info"].(map[string]interface{})
	if !ok {
		err := errors.New("torrent file corrupt: infoMap is not a dictionary")
		return nil, err
	}

	pieceArray, ok := infoMap["pieces"].([]byte)
	if !ok {
		err := errors.New("torrent file corrupt: pieces is not an array")
		return nil, err
	}
	hashes, err := splitHashes(pieceArray)
	if err != nil {
		return nil, err
	}

	infoHash, err := infoHash(dataMap["info"])
	if err != nil {
		return nil, err
	}

	files, err := parseFilePaths(infoMap)
	if err != nil {
		return nil, err
	}

	announce, ok := dataMap["announce"].([]byte)
	if !ok {
		err := errors.New("torrent file corrupt: tracker url missing")
		return nil, err
	}
	pieceLength, ok := infoMap["piece length"].(int64)
	if !ok {
		err := errors.New("torrent file corrupt: piece length missing")
		return nil, err
	}
	name, ok := infoMap["name"].([]byte)
	if !ok {
		return nil, errors.New("torrent name is missing")
	}

	tor := Torrent{
		Announce:    string(announce),
		InfoHash:    *infoHash,
		PieceHashes: hashes,
		PieceLength: pieceLength,
		Files:       files,
		Name:        string(name),
	}
	return &tor, nil
}

func FromFile(path string) (*Torrent, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(content)
}
