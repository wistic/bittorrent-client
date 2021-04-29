package peer

import "bittorrent-go/util"

type PeerScheduler struct {
	Map   map[string]bool
	Index int
	List  []util.Address
}

func NewPeerScheduler() *PeerScheduler {
	return &PeerScheduler{
		Map:   make(map[string]bool),
		Index: 0,
		List:  nil,
	}
}

func (ps *PeerScheduler) UpdateList(list []util.Address) {
	ps.Index = 0
	ps.List = list
}

func (ps *PeerScheduler) AddAddress(address util.Address) {
	ps.Map[address.String()] = true
}

func (ps *PeerScheduler) RemoveAddress(address util.Address) {
	delete(ps.Map, address.String())
}

func (ps *PeerScheduler) Total() int {
	return len(ps.Map)
}

func (ps *PeerScheduler) Next() *util.Address {
	if ps.List == nil || len(ps.List) == 0 {
		return nil
	}
	traverse := 0
	for traverse < len(ps.List) {
		address := ps.List[ps.Index]
		ps.Index = (ps.Index + 1) % len(ps.List)
		if _, ok := ps.Map[address.String()]; !ok {
			return &address
		}
		traverse += 1
	}
	return nil
}
