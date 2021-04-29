package job

import (
	"bittorrent-go/torrent"
)

func CreateJobChannel(tor *torrent.Torrent) (chan *Job, int) {
	pieceChan := make(chan *Job, len(tor.PieceHashes))
	total := tor.Length()
	for i := 0; i < len(tor.PieceHashes); i += 1 {
		if total > tor.PieceLength {
			pieceChan <- &Job{
				Index:  uint32(i),
				Length: uint32(tor.PieceLength),
				Hash:   tor.PieceHashes[i],
			}
			total -= tor.PieceLength
		} else {
			pieceChan <- &Job{
				Index:  uint32(i),
				Length: uint32(total),
				Hash:   tor.PieceHashes[i],
			}
		}
	}
	return pieceChan, len(tor.PieceHashes)
}

func CreateResultChannel() chan *Result {
	return make(chan *Result)
}
