package filesystem

import (
	"bittorrent-go/torrent"
	"bittorrent-go/util"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func AssembleRoutine(tor *torrent.Torrent) error {
	reader := SerialReader{
		torrent: tor,
		current: nil,
		offset:  0,
		index:   0,
	}

	for _, file := range tor.Files {
		fmt.Println("creating", file.Path)

		f, err := os.Create(file.Path)
		if err != nil {
			return err
		}

		_, err = io.CopyN(f, &reader, file.Length)
		if err != nil {
			f.Close()
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}

		fmt.Println("created", file.Path)
	}

	return nil
}

type SerialReader struct {
	torrent *torrent.Torrent
	current []byte
	offset  int
	index   uint32
}

func (sr *SerialReader) Read(p []byte) (n int, err error) {
	if sr.current == nil || sr.offset >= len(sr.current) {
		fmt.Println("reading", sr.index)

		sr.current, err = ioutil.ReadFile(fileName(sr.index))
		if err != nil {
			return 0, err
		}
		hash := util.Hash{Value: sha1.Sum(sr.current)}

		if !hash.Match(&sr.torrent.PieceHashes[sr.index]) {
			return 0, errors.New("hash mismatch")
		}

		sr.index += 1
		sr.offset = 0
	}

	cp := copy(p[:], sr.current[sr.offset:])
	sr.offset += cp
	return cp, nil
}
