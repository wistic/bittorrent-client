package job

import "bittorrent-go/util"

type Job struct {
	Index  uint32
	Length uint32
	Hash   util.Hash
}
