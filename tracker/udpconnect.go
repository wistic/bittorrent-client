package tracker

import (
	"encoding/binary"
	"errors"
	"net"
	"time"
)

// returns Conn interface and connection id to be used for Announce Part and requires url of tracker in string form

func ConnectUDP(url string) (net.Conn, uint64, error) {
	wbuffer := make([]byte, 16)                                    //buffer with connect request packet
	rbuffer := make([]byte, 16)                                    //buffer to accept connect Response packet
	binary.BigEndian.PutUint64(wbuffer[0:], uint64(0x41727101980)) // protocol id magic constant
	binary.BigEndian.PutUint32(wbuffer[8:], uint32(0))             // action(connect) 0
	binary.BigEndian.PutUint32(wbuffer[12:], uint32(0))            //transaction id set as 0

	conn, err := net.DialTimeout("udp", url, 15*time.Second)
	if err != nil {
		return nil, 0, err
	}

	var n, retries int
	for {
		retries++
		// time out for writting packet added
		conn.SetWriteDeadline(time.Now().Add(15 * time.Second))

		n, err = conn.Write(wbuffer)
		if err != nil {
			return nil, 0, err
		}
		if n != len(wbuffer) {
			return nil, 0, errors.New("udp packet was not entirely written")
		}

		conn.SetReadDeadline(time.Now().Add(15 * time.Second))

		n, err = conn.Read(rbuffer)
		if e, ok := err.(net.Error); ok && e.Timeout() {
			if retries > 3 {
				return nil, 0, errors.New("Retries limit reached")
			}
			continue
		} else if err != nil {
			return nil, 0, err
		}
		break
	}

	if n != len(rbuffer) {
		return nil, 0, errors.New("invalid Response received from tracker")
	}
	action := binary.BigEndian.Uint32(rbuffer[0:])
	if action != uint32(0) {
		return nil, 0, errors.New("invalid action")
	}
	tid := binary.BigEndian.Uint32(rbuffer[4:])
	if tid != uint32(0) {
		return nil, 0, errors.New("Transaction Id donot match")
	}
	connectionid := binary.BigEndian.Uint64(rbuffer[8:])

	return conn, connectionid, nil
}
