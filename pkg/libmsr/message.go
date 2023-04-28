package libmsr

import (
	"bytes"
	"errors"
)

const (
	seqStartBit byte = 1 << 7
	seqEndBit   byte = 1 << 6
	escByte     byte = 0x1B
	fsByte      byte = 0x1C
)

func esc(data ...byte) []byte {
	b := make([]byte, len(data)+1)
	b[0] = escByte
	copy(b[1:], data)
	return b
}

func makePackets(msg []byte) [][]byte {
	n := len(msg) / 63
	if len(msg)%63 != 0 {
		n++
	}
	pkts := make([][]byte, n)
	for i := range pkts {
		pkt := make([]byte, 64)
		if i == 0 {
			pkt[0] |= seqStartBit
		}
		if i == n-1 {
			pkt[0] |= seqEndBit | byte(len(msg[i:]))
		} else {
			pkt[0] |= 63
		}
		copy(pkt[1:], msg[i*63:])
		pkts[i] = pkt
	}
	return pkts
}

func parsePackets(pkts [][]byte) ([]byte, error) {
	n := 0
	for _, pkt := range pkts {
		n += int(pkt[0] & 63)
	}
	msg := make([]byte, n)
	i := 0
	for pktI, pkt := range pkts {
		if pktI == 0 != (pkt[0]&seqStartBit == seqStartBit) {
			return nil, errors.New("invalid start bit")
		}
		pktLen := int(pkt[0] & 63)
		copy(msg[i:i+pktLen], pkt[1:1+pktLen])
		i += pktLen
	}
	return msg, nil
}

func decode(msg []byte) (data, result []byte, err error) {
	escI := bytes.LastIndexByte(msg, escByte)
	err = Status(msg[escI+1])
	if err == StatusOK {
		err = nil
	} else {
		return
	}
	data, result = msg[:escI], msg[escI+2:]
	return
}

func encodeTrack(data []byte, num int) []byte {
	if num < 1 || num > 3 {
		panic("invalid track")
	}
	return append(esc(byte(num), byte(len(data))), data...)
}

func encodeRawTracks(tracks ...[]byte) []byte {
	data := esc('s')
	for i, t := range tracks {
		if true || len(t) > 0 { // TODO
			data = append(data, encodeTrack(t, i+1)...)
		}
	}
	return append(data, '?', fsByte)
}

func encodeISOTracks(t1, t2, t3 []byte) []byte {
	data := esc('s')
	data = append(data, encodeTrack(t1, 1)...)
	data = append(data, encodeTrack(t2, 2)...)
	data = append(data, encodeTrack(t3, 3)...)
	return append(data, '?', fsByte)
}

func decodeTracks(data []byte) ([3][]byte, error) {
	var tracks [3][]byte
	if data[0] != escByte || data[1] != 's' {
		return tracks, errors.New("invalid track block start")
	}
	data = data[2:]
	for i := 0; i < 3; i++ {
		if data[0] != escByte || data[1] != byte(i+1) {
			return tracks, errors.New("invalid track data")
		}
		trackLen := int(data[2])
		data = data[3:]
		tracks[i] = data[:trackLen]
		data = data[trackLen:]
	}
	if data[0] != '?' || data[1] != fsByte {
		return tracks, errors.New("invalid track block end")
	}
	return tracks, nil
}
