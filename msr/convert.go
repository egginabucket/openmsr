package msr

func hexToBin(hex []byte) []bool {
	bin := make([]bool, len(hex)*4)
	for i, h := range hex {
		if h >= 'A' {
			if h >= 'a' {
				h -= 'a' - 'A'
			}
			h -= 'A' - '9' - 1
		}
		h -= '0'
		for b := 0; b < 4; b++ {
			bin[i*4+b] = h&(1<<3) == 1<<3
			h <<= 1
		}
	}
	return bin
}

func binToChars(bin []bool, bpc int, parityEven bool) (chars []byte, ok []bool) {
	chars = make([]byte, len(bin)/bpc)
	ok = make([]bool, len(chars))
	for i := range chars {
		binVal := bin[i*bpc : (i+1)*bpc-1]
		parity := parityEven != bin[(i+1)*bpc-1]
		var c byte
		for _, bit := range binVal {
			c <<= 1
			if bit {
				c |= 1
				parity = !parity
			}
		}
		switch bpc { // TODO
		case 7:
			c += ' '
		case 4, 5, 6:
			c += '0'
		}
		chars[i], ok[i] = c, parity
	}
	return chars, ok
}

func HexToChar(hex []byte, bpc int, parityEven bool) (chars []byte, ok []bool) {
	return binToChars(hexToBin(hex), bpc, parityEven)
}
