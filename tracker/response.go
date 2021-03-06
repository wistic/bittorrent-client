package tracker

import (
	"fmt"
	"github.com/IncSW/go-bencode"
)

func Parse(resp []byte) {
	data, err := bencode.Unmarshal(resp)
	if err != nil {

	}
	fmt.Println(data)
}
