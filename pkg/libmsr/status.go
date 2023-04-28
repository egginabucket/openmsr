package libmsr

import (
	"fmt"
)

type Status byte

const (
	StatusOK                Status = '0'
	StatusReadWriteErr      Status = '1'
	StatusInvalidCommandFmt Status = '2'
	StatusInvalidCommand    Status = '4'
	StatusWriteSwipeErr     Status = '9'
	StatusFail              Status = 'A'
)

func (s Status) Error() string {
	const pre = "MSR status: "
	switch s {
	case StatusOK:
		return ""
	case StatusReadWriteErr:
		return pre + "read/write error"
	case StatusInvalidCommandFmt:
		return pre + "invalid command format"
	case StatusInvalidCommand:
		return pre + "invalid command"
	case StatusWriteSwipeErr:
		return pre + "write mode swipe error"
	case StatusFail:
		return pre + "fail"
	}
	return fmt.Sprintf("unknown status byte %c", s)
}
