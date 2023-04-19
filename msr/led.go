package msr

import "errors"

type LEDMode byte

var ErrInvalidLEDMode = errors.New("invalid LED mode")

const (
	LEDAllOff LEDMode = iota + 0x81
	LEDAllOn
	LEDGreenOn
	LEDYellowOn
	LEDRedOn
)

func (m LEDMode) ok() bool {
	return m >= LEDAllOff && m <= LEDRedOn
}
