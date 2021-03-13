package tracker

import (
	"bittorrent-go/util"
	"encoding/binary"
	"errors"
	"net"
	"time"
)

type AnnounceResponse struct {
	action   uint32
	tid      uint32
	interval uint32
	leechers uint32
	seeders  uint32
}

// Return Announce Response and Peer list and requires torrent infohash , Conn Interface from udpconnect.go and connection id

func AnnounceUDP(infohash *util.Hash, peerId *util.PeerID, port uint16, conn net.Conn, cid uint64) (*Response, error) {

	buffer := make([]byte, 98)                         // Announce Request packet
	binary.BigEndian.PutUint64(buffer[0:], cid)        // connection id that we got from connect Response packet
	binary.BigEndian.PutUint32(buffer[8:], uint32(1))  // action (announce) 1
	binary.BigEndian.PutUint32(buffer[12:], uint32(0)) // transaction id 0
	copy(buffer[16:], infohash.Value[:])               // info hash in bytes
	copy(buffer[36:], peerId.Value[:])
	binary.BigEndian.PutUint64(buffer[56:], uint64(0)) //downloaded 0
	binary.BigEndian.PutUint64(buffer[64:], uint64(0)) // left 0
	binary.BigEndian.PutUint64(buffer[72:], uint64(0)) // uploaded 0
	binary.BigEndian.PutUint32(buffer[80:], uint32(0)) // 0: none; 1: completed; 2: started; 3: stopped
	binary.BigEndian.PutUint32(buffer[84:], uint32(0)) //IP address (default) 0
	binary.BigEndian.PutUint32(buffer[88:], uint32(0)) // key 0
	copy(buffer[92:], []byte{0xFF, 0xFF, 0xFF, 0xFF})  // num_want (default) -1
	binary.BigEndian.PutUint16(buffer[96:], port)      // port no.

	//Creating a duplicate packet for future
	dbuffer := make([]byte, 98)
	copy(dbuffer, buffer)

	var n, retries int
	var err error
	for {
		retries++

		conn.SetWriteDeadline(time.Now().Add(15 * time.Second))

		n, err = conn.Write(buffer)
		if err != nil {
			return nil, err
		}
		if n != len(buffer) {
			return nil, errors.New("udp packet was not entirely written")
		}

		conn.SetReadDeadline(time.Now().Add(15 * time.Second))

		n, err = conn.Read(buffer)
		if e, ok := err.(net.Error); ok && e.Timeout() {
			if retries > 3 {
				return nil, errors.New("retries limit reached")
			}
			continue
		} else if err != nil {
			return nil, err
		}
		break
	}
	//Doubt to keep as Response packet>98
	if n != len(buffer) {
		return nil, errors.New("invalid Response received from tracker 1")
	}
	action := binary.BigEndian.Uint32(buffer[0:])
	if action != uint32(1) {
		return nil, errors.New("invalid action")
	}
	tid := binary.BigEndian.Uint32(buffer[4:])
	if tid != uint32(0) {
		return nil, errors.New("transaction Id don't match")
	}

	leechers := binary.BigEndian.Uint32(buffer[12:])
	seeders := binary.BigEndian.Uint32(buffer[16:])
	N := seeders + leechers
	buf := make([]byte, 20+6*N)

	for {
		//Sending Announce Packet to tracker
		conn.SetWriteDeadline(time.Now().Add(15 * time.Second))
		n, err = conn.Write(dbuffer)
		if err != nil {
			return nil, err
		}
		if n != len(dbuffer) {
			return nil, errors.New("udp packet was not entirely written")
		}
		//Reading the Response packet

		conn.SetReadDeadline(time.Now().Add(15 * time.Second))

		n, err = conn.Read(buf)
		if e, ok := err.(net.Error); ok && e.Timeout() {
			if retries > 3 {
				return nil, errors.New("retries limit reached")
			}
			continue
		} else if err != nil {
			return nil, err
		}
		break
	}

	if (n-20)%6 != 0 {
		return nil, errors.New("invalid Response received from tracker 2")
	}

	action = binary.BigEndian.Uint32(buf[0:])
	if action != uint32(1) {
		return nil, errors.New("invalid action")
	}
	tid = binary.BigEndian.Uint32(buf[4:])
	if tid != uint32(0) {
		return nil, errors.New("transaction Id do not match")
	}

	interval := binary.BigEndian.Uint32(buf[8:])
	_ = binary.BigEndian.Uint32(buf[12:]) // leechers
	_ = binary.BigEndian.Uint32(buf[16:]) // seeders
	peer := make([]util.Address, (n-20)/6)

	for i := 0; i < (n-20)/6; i++ {
		offset := 20 + (i * 6)
		peer[i].IP = buf[offset : offset+4]
		peer[i].Port = binary.BigEndian.Uint16(buf[offset+4:])
	}

	return &Response{uint64(interval), peer}, nil // returns AnnounceResponse Struct and Peer info
}
