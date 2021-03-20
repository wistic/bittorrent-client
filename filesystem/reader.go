package filesystem

import "bittorrent-go/torrent"

func Read(torrent *torrent.Torrent, path string) []bool {
	res := make([]bool, len(torrent.PieceHashes))
	for i := 0; i < len(res); i++ {
		res[i] = false
	}
	return res
}
