package libmsr

func DecodeRaw(raw []byte, offset byte, bpcRaw, bpcChars int, parityEven bool) (chars []byte, parityOK []bool, lrcOK bool) {
	chars, parityOK = make([]byte, 0), make([]bool, 0)
	if len(raw) < 1 {
		lrcOK = true
		return
	}
	remCount := 0
	var remBits uint16
	var lrc byte
	lastNonNull := -1
	for _, b := range raw {
		remCount += bpcChars
		remBits <<= bpcChars
		remBits |= uint16(b) & ((0x01 << bpcChars) - 1)
		for remCount >= bpcRaw {
			remCount -= bpcRaw
			i := byte(remBits >> remCount)
			remBits &= (0x01 << remCount) - 1

			p := i & 0x01
			var c byte
			for j := bpcRaw - 1; j > 0; j-- {
				i >>= 1
				p ^= i & 0x01
				c |= (i & 0x01) << (j - 1)
			}
			if c != 0x00 {
				lastNonNull = len(chars)
			}
			lrc ^= c
			c += offset
			/*
				switch c {
				case ',':
					c = '`'
				case '-':
					c = ','
				}
			*/
			chars = append(chars, c)
			parityOK = append(parityOK, parityEven != (p&0x01 == 0x01))

		}
	}
	chars, parityOK = chars[:lastNonNull], parityOK[:lastNonNull]
	lrcOK = lrc == 0
	return
}

func EncodeRaw(chars []byte, offset byte, bpcRaw, bpcChars int, parityEven bool) []byte {
	raw := make([]byte, 0)
	remCount := 0
	var remBits uint16
	var lrc byte
	for _, c := range chars {
		/*
			switch c {
			case ',':
				c = '-'
			case '`':
				c = ','
			}
		*/
		c -= offset
		lrc ^= c
		for i := 0; i < bpcRaw-1; i++ {
			c ^= ((c >> i) & 0x01) << (bpcRaw - 1)
		}
		remBits |= uint16(c) << remCount
		remCount += bpcRaw
		if remCount >= bpcChars {
			raw = append(raw, byte(remBits)&((0x01<<bpcChars)-1))
			remBits >>= bpcChars
			remCount -= bpcChars
		}
	}
	for i := 0; i < bpcRaw-1; i++ {
		lrc ^= ((lrc >> i) & 0x01) << (bpcRaw - 1)
	}
	remBits |= uint16(lrc) << remCount
	remCount += bpcRaw
	if remCount >= bpcChars {
		raw = append(raw, byte(remBits)&((0x01<<bpcChars)-1))
		remBits >>= bpcChars
		remCount -= bpcChars
	}
	if remCount > 0 {
		raw = append(raw, byte(remBits))
	}
	return raw
}
