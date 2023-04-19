package msr

import (
	"errors"
	"fmt"
)

const (
	StatusOK            byte = '0'
	StatusReadWriteErr  byte = '1'
	StatusCommandFmtErr byte = '2'
	StatusCommandErr    byte = '4'
	StatusWriteSwipeErr byte = '9'
	StatusFail          byte = 'A'
)

var (
	ErrReadWrite         = errors.New("read/write error")
	ErrInvalidCommandFmt = errors.New("invalid command format")
	ErrInvalidCommand    = errors.New("invalid command")
	ErrWriteSwipe        = errors.New("write mode swipe error")
	ErrFail              = errors.New("MSR fail")
)

func statusErr(b byte) error {
	switch b {
	case StatusOK:
		return nil
	case StatusReadWriteErr:
		return ErrReadWrite
	case StatusCommandFmtErr:
		return ErrInvalidCommandFmt
	case StatusCommandErr:
		return ErrInvalidCommand
	case StatusWriteSwipeErr:
		return ErrWriteSwipe
	case StatusFail:
		return ErrFail
	}
	return fmt.Errorf("unknown status byte %X", b)
}
