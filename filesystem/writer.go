package filesystem

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

func fileName(index uint32) string {
	return strconv.Itoa(int(index)) + ".piece"
}

func WriteRoutine(wg *sync.WaitGroup, index uint32, data []byte) {
	wg.Add(1)
	defer wg.Done()

	file, err := os.Create(fileName(index))
	if err != nil {
		fmt.Println("[ writer ] ", "file create error: ", err)
		return
	}
	_, err = file.Write(data)
	if err != nil {
		fmt.Println("[ writer ] ", "file write error: ", err)
		return
	}

	err = file.Sync()
	if err != nil {
		fmt.Println("[ writer ] ", "file sync error: ", err)
		return
	}

	err = file.Close()
	if err != nil {
		fmt.Println("[ writer ] ", "file close error: ", err)
		return
	}

}
