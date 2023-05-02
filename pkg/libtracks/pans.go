package libtracks

import (
	"errors"
	"unicode"
)

// Primary account number (eg. credit card number), 1-19 digits
type PAN struct {
	digits []int
}

// MII returns the major industry identifier, identified from the first digit.
func (n *PAN) MII() string {
	const future = " and other future industry assigments"
	switch n.digits[0] {
	case 0:
		return "ISO/TC 68 and other industry assignments"
	case 1:
		return "Airlines"
	case 2:
		return "Airlines, financial" + future
	case 3:
		return "Travel and entertainment"
	case 4, 5:
		return "Banking and financial"
	case 6:
		return "Merchandising and banking/financial"
	case 7:
		return "Petroleum" + future
	case 8:
		return "Healthcare, telecommunications" + future
	case 9:
		return "For assignment by national standards bodies"
	}
	return ""
}

func (n *PAN) checksum(shift bool) int {
	sum := 0
	parity := n.Len() % 2
	for i, d := range n.digits {
		if d < 0 {
			return -1
		} else if shift != (i%2 == parity) {
			sum += d
		} else if d > 4 {
			sum += 2*d - 9
		} else {
			sum += 2 * d
		}
	}
	return sum % 10
}

// IsLuhnValid checks against accidental errors with the Luhn algorithm.
func (n *PAN) IsLuhnValid() bool {
	return n.checksum(true) == 0
}

// WriteLRC writes a longitudinal redundancy check digit calculated with the Luhn alg.
func (n *PAN) WriteLRC() {
	n.digits = append(n.digits, (10-n.checksum(false))%10)
}

func (n *PAN) Len() int {
	return len(n.digits)
}

func (n *PAN) String() string {
	b := make([]byte, n.Len())
	for i, d := range n.digits {
		if d < 0 {
			b[i] = 'x'
		} else {
			b[i] = byte(d) + '0'
		}
	}
	return string(b)
}

func (n *PAN) Info() *Info {
	return NewInfo("PAN", n.String(),
		NewInfo("MII", n.MII()),
		NewInfo("Luhn valid", n.IsLuhnValid()),
	)
}

func NewPAN(s string) (PAN, error) {
	var n PAN
	n.digits = make([]int, 0, 19)
	for _, r := range s {
		if '0' <= r && r <= '9' {
			n.digits = append(n.digits, int(r-'0'))
		} else if r != '-' && r != '_' && !unicode.IsSpace(r) {
			n.digits = append(n.digits, -1)
		}
	}
	if n.Len() < 1 {
		return n, errors.New("libtracks.NewPAN: no data")
	}
	return n, nil
}
