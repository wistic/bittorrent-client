package filesystem

import (
	"bittorrent-go/torrent"
	"bittorrent-go/util"
	"crypto/sha1"
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
)

func AssembleRoutine(tor *torrent.Torrent, outputPath string) error {
	reader := SerialReader{
		torrent:    tor,
		outputPath: outputPath,
		current:    nil,
		offset:     0,
		index:      0,
	}

	for _, file := range tor.Files {
		logrus.Infoln("creating:", file.Path)

		f, err := os.Create(path.Join(outputPath, file.Path))
		if err != nil {
			return err
		}

		_, err = io.CopyN(f, &reader, file.Length)
		if err != nil {
			err2 := f.Close()
			if err2 != nil {
				return err2
			}
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}

		logrus.Infoln("finished:", file.Path)
	}

	return nil
}

type SerialReader struct {
	torrent    *torrent.Torrent
	outputPath string
	current    []byte
	offset     int
	index      uint32
}

func (sr *SerialReader) Read(p []byte) (n int, err error) {
	if sr.current == nil || sr.offset >= len(sr.current) {
		logrus.Debugln("reading:", sr.index)

		sr.current, err = ioutil.ReadFile(fileName(sr.index, sr.outputPath))
		if err != nil {
			logrus.Errorln("file read error:", err)
			return 0, err
		}
		hash := util.Hash{Value: sha1.Sum(sr.current)}

		if !hash.Match(&sr.torrent.PieceHashes[sr.index]) {
			logrus.Errorln("hash mismatch")
			return 0, errors.New("hash mismatch")
		}

		sr.index += 1
		sr.offset = 0
	}

	cp := copy(p[:], sr.current[sr.offset:])
	sr.offset += cp
	return cp, nil
}
