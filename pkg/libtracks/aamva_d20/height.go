package aamvad20

import (
	"fmt"
	"io"
)

type Height struct {
	Feet, Inches int
}

func (h *Height) String() string {
	if h == nil {
		return ""
	}
	return fmt.Sprintf("%d' %d\" or %dcm", h.Feet, h.Inches, h.CM())
}

// Not converted, just the value used in CA
func (h *Height) CM() int {
	return h.Feet*100 + h.Inches
}

func (h *Height) Write3Digits(w io.Writer) {
	fmt.Fprintf(w, "%d%02d", h.Feet, h.Inches)
}

func ParseHeight(s string) *Height {
	var h Height
	if s[0] != ' ' {
		h.Feet = int(s[0] - '0')
	}
	h.Inches = (int(s[1]-'0') * 10) + int(s[2]-'0')
	return &h
}
