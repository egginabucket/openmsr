package cards

import (
	"errors"
	"unicode"
)

type Card struct {
	digits []int
}

func (c *Card) MII() string {
	return MII(c.digits[0])
}

func (c *Card) checksum(shift bool) int {
	sum := 0
	parity := c.Len() % 2
	for i, d := range c.digits {
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

func (c *Card) IsLuhnValid() bool {
	return c.checksum(true) == 0
}

func (c *Card) WriteLRC() {
	c.digits = append(c.digits, (10-c.checksum(false))%10)
}

func (c *Card) Len() int {
	return len(c.digits)
}

func (c *Card) String() string {
	b := make([]byte, c.Len())
	for i, d := range c.digits {
		if d < 0 {
			b[i] = 'x'
		} else {
			b[i] = byte(d) + '0'
		}
	}
	return string(b)
}

func NewCard(s string) (*Card, error) {
	var c Card
	c.digits = make([]int, 0, 19)
	for _, r := range s {
		if '0' <= r && r <= '9' {
			c.digits = append(c.digits, int(r-'0'))
		} else if r != '-' && r != '_' && !unicode.IsSpace(r) {
			c.digits = append(c.digits, -1)
		}
	}
	if c.Len() < 1 {
		return nil, errors.New("empty card")
	}
	return &c, nil
}
