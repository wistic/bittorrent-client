package filesystem

import (
	"bittorrent-go/torrent"
	"github.com/sirupsen/logrus"
	"os"
)

func checkPath(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func DeleteRoutine(tor *torrent.Torrent, outputPath string) {
	for i := 0; i < len(tor.PieceHashes); i++ {
		if checkPath(fileName(uint32(i), outputPath)) {
			err := os.Remove(fileName(uint32(i), outputPath))
			if err != nil {
				logrus.Errorln("delete error:", err)
			}
		}
	}
}
