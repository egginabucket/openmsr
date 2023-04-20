package msr

import (
	"errors"
	"fmt"
)

const (
	statusOK            byte = '0'
	statusReadWriteErr  byte = '1'
	statusCommandFmtErr byte = '2'
	statusCommandErr    byte = '4'
	statusWriteSwipeErr byte = '9'
	statusFail          byte = 'A'
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
	case statusOK:
		return nil
	case statusReadWriteErr:
		return ErrReadWrite
	case statusCommandFmtErr:
		return ErrInvalidCommandFmt
	case statusCommandErr:
		return ErrInvalidCommand
	case statusWriteSwipeErr:
		return ErrWriteSwipe
	case statusFail:
		return ErrFail
	}
	return fmt.Errorf("unknown status byte %X", b)
}
