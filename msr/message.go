package msr

import (
	"bytes"
	"errors"
)

//const esc byte = 0x1B

const (
	seqStartBit byte = 1 << 7
	seqEndBit   byte = 1 << 6
	escByte     byte = 0x1B
)

func esc(data ...byte) []byte {
	b := make([]byte, len(data)+1)
	b[0] = escByte
	copy(b[1:], data)
	return b
}

func makePackets(msg []byte) [][]byte {
	pkts := make([][]byte, 0, len(msg)/63+1)
	for i := 0; i < len(msg); i += 63 {
		pkt := make([]byte, 64)
		if i == 0 {
			pkt[0] |= seqStartBit
		}
		if i+63 >= len(msg) {
			pkt[0] |= seqEndBit | byte(len(msg[i:]))
			copy(pkt[1:], msg[i:])
		} else {
			pkt[0] |= 63
			copy(pkt[1:], msg[i:i+63])
		}
		pkts = append(pkts, pkt)
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
	err = statusErr(msg[escI+1])
	if err != nil {
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

func encodeRaw(tracks ...[]byte) []byte {
	data := esc('s')
	for i, t := range tracks {
		if len(t) > 0 {
			data = append(data, encodeTrack(t, i+1)...)
		}
	}
	return append(data, '?', 0x1C)
}

func encodeISO(t1, t2, t3 []byte) []byte {
	data := esc('s')
	data = append(data, encodeTrack(t1, 1)...)
	data = append(data, encodeTrack(t2, 2)...)
	data = append(data, encodeTrack(t3, 3)...)
	return append(data, '?', 0x1C)
}
