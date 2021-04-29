package filesystem

import (
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strconv"
	"sync"
)

func fileName(index uint32, outputPath string) string {
	return path.Join(outputPath, strconv.Itoa(int(index))+".piece")
}

func WriteRoutine(wg *sync.WaitGroup, index uint32, data []byte, outputPath string) {
	wg.Add(1)
	defer wg.Done()

	file, err := os.Create(fileName(index, outputPath))
	if err != nil {
		logrus.Errorln("file create error:", err)
		return
	}
	_, err = file.Write(data)
	if err != nil {
		logrus.Errorln("file write error:", err)
		return
	}

	err = file.Sync()
	if err != nil {
		logrus.Errorln("file sync error:", err)
		return
	}

	err = file.Close()
	if err != nil {
		logrus.Errorln("file close error:", err)
		return
	}

}
