package scheduler

import (
	"bittorrent-go/util"
	"sync"
)

type Scheduler struct {
	Mutex        sync.Mutex
	PieceCounter uint32 // TODO: Replace with better scheduling
}

func (scheduler *Scheduler) GetPiece(address *util.Address) (uint32, bool) {
	scheduler.Mutex.Lock()
	defer scheduler.Mutex.Unlock()
	counter := scheduler.PieceCounter
	scheduler.PieceCounter++
	return counter, true
}
