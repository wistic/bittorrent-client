package filesystem

import (
	"bittorrent-go/torrent"
	"bittorrent-go/util"
	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
)

const BlockSize = int64(16384)

type Block struct {
	PieceIndex int
	BlockIndex int
	Data       []byte
}

type Piece struct {
	Block      []bool
	Data       []byte
	BlockCount int
}

func newPiece(size int64) *Piece {
	extra := size % BlockSize
	count := size / BlockSize
	if extra != 0 {
		count++
	}
	return &Piece{
		Block:      make([]bool, count),
		Data:       make([]byte, size),
		BlockCount: 0,
	}
}

func (piece *Piece) write(index int, data []byte) error {
	if index < 0 || index >= len(piece.Block) {
		return errors.New("block index out of range")
	}
	if piece.Block[index] {
		return nil
	}
	start := int64(index) * BlockSize
	copy(piece.Data[start:], data[:])
	piece.Block[index] = true
	piece.BlockCount++
	return nil
}

func (piece *Piece) finished() bool {
	if piece.BlockCount < len(piece.Block) {
		return false
	}
	// TODO: Verify and skip this part
	for i := 0; i < len(piece.Block); i++ {
		if !piece.Block[i] {
			return false
		}
	}
	return true
}

func (piece *Piece) hash() *util.Hash {
	return &util.Hash{Value: sha1.Sum(piece.Data)}
}

func getPiece(torrent *torrent.Torrent, pieceMap map[int]*Piece, index int) (*Piece, error) {
	p, ok := pieceMap[index]
	if ok {
		return p, nil
	}
	if index < 0 || index >= len(torrent.PieceHashes) {
		return nil, errors.New("piece index out of range")
	}
	var piece *Piece = nil
	if index < len(torrent.PieceHashes)-1 {
		piece = newPiece(torrent.PieceLength)
	} else {
		piece = newPiece(torrent.Length() % torrent.PieceLength)
	}
	pieceMap[index] = piece
	return piece, nil
}

func StartWriter(torrent *torrent.Torrent, output string) (chan *Block, chan int, chan error, chan struct{}) {
	block := make(chan *Block)
	finish := make(chan int)
	err := make(chan error)
	done := make(chan struct{})
	go writerRoutine(torrent, output, block, finish, err, done)
	return block, finish, err, done
}

func StopWriter(done chan struct{}) {
	close(done)
}

func writerRoutine(torrent *torrent.Torrent, output string, blockChannel chan *Block, finishChannel chan int, errorChannel chan error, done chan struct{}) {
	fmt.Println("opening writer")
	defer fmt.Println("closing writer")
	complete := make([]bool, len(torrent.PieceHashes))
	pieceMap := make(map[int]*Piece)

	for {
		select {
		case block := <-blockChannel:
			if complete[block.PieceIndex] {
				fmt.Println("piece already completed")
				break
			}
			index := block.PieceIndex
			piece, err := getPiece(torrent, pieceMap, index)
			if err != nil {
				errorChannel <- err
				break
			}
			err = piece.write(block.BlockIndex, block.Data)
			if err != nil {
				errorChannel <- err
				break
			}
			if piece.finished() {
				hash := piece.hash()
				if !hash.Match(&torrent.PieceHashes[index]) {
					delete(pieceMap, index)
					errorChannel <- errors.New("piece hash mismatch")
					break
				}
				go dumpRoutine(piece, index, path.Join(output, fmt.Sprint(index, ".piece")), finishChannel, errorChannel)
				complete[index] = true
			}

		case <-done:
			return
		}
	}
}

func dumpRoutine(piece *Piece, index int, path string, finishChannel chan int, errorChannel chan error) {
	fmt.Println("writing piece: ", index)
	err := ioutil.WriteFile(path, piece.Data, 0)
	if err != nil {
		errorChannel <- err
		return
	}
	finishChannel <- index
}
