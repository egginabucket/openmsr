package libtracks

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

const expDateLayout = "0601" // YYMM

type (
	IEC7813Track1 struct {
		PAN           PAN
		Name          string
		ExpDate       *time.Time
		ServiceCode   string
		Discretionary string
	}
	IEC7813Track2 struct {
		PAN           PAN
		ExpDate       *time.Time
		ServiceCode   string
		Discretionary string
	}
	IEC7813Tracks struct {
		Track1 *IEC7813Track1
		Track2 *IEC7813Track2
	}
)

var iec7813Track1Re = regexp.MustCompile(`^%B([0-9]{1,19})\^([^\^]{2,26})\^([0-9]{4}|\^)([0-9]{3}|\^)([^\?]+)\?$`)
var iec7813Track2Re = regexp.MustCompile(`^;([0-9]{1,19})\=([0-9]{4}|\=)([0-9]{3}|\=)([^\?]+)\?$`)

func (t *IEC7813Track1) String() string {
	var b strings.Builder
	b.WriteString("%B")
	b.WriteString(t.PAN.String())
	b.WriteByte('^')
	b.WriteString(t.Name)
	b.WriteByte('^')
	if t.ExpDate == nil {
		b.WriteByte('^')
	} else {
		b.WriteString(t.ExpDate.Format(expDateLayout))
	}
	if t.ServiceCode == "" {
		b.WriteByte('^')
	} else {
		b.WriteString(t.ServiceCode)
	}
	b.WriteString(t.Discretionary)
	b.WriteByte('?')
	return b.String()
}

func (t *IEC7813Track1) Info() *Info {
	if t == nil {
		return NewInfo("Track 1", "nil")
	}
	return NewInfo("Track 1", t.String(),
		&t.PAN,
		NewInfo("Name", t.Name),
		NewInfo("Expiration date", t.ExpDate.Format("2006-01")),
		NewInfo("Service code", t.ServiceCode),
		NewInfo("Discretionary data", t.Discretionary),
	)
}

func (t *IEC7813Track2) String() string {
	var b strings.Builder
	b.WriteByte(';')
	b.WriteString(t.PAN.String())
	b.WriteByte('=')
	if t.ExpDate == nil {
		b.WriteByte('=')
	} else {
		b.WriteString(t.ExpDate.Format(expDateLayout))
	}
	if t.ServiceCode == "" {
		b.WriteByte('=')
	} else {
		b.WriteString(t.ServiceCode)
	}
	b.WriteString(t.Discretionary)
	b.WriteByte('?')
	return b.String()
}

func (t *IEC7813Track2) Info() *Info {
	if t == nil {
		return NewInfo("Track 2", "nil")
	}
	return NewInfo("Track 2", t.String(),
		&t.PAN,
		NewInfo("Expiration date", t.ExpDate.Format("2006-01")),
		NewInfo("Service code", t.ServiceCode),
		NewInfo("Discretionary data", t.Discretionary),
	)
}

func NewIEC7813Track1(s string) (*IEC7813Track1, error) {
	var err error
	groups := iec7813Track1Re.FindStringSubmatch(s)
	if groups == nil {
		return nil, errors.New("invalid IEC 7813 track 1")
	}
	var t IEC7813Track1
	t.PAN, err = NewPAN(groups[1])
	if err != nil {
		return nil, err
	}
	t.Name = groups[2]
	if groups[3] != "^" {
		time, err := time.Parse(expDateLayout, groups[3])
		if err != nil {
			return nil, err
		}
		t.ExpDate = &time
	}
	if groups[4] != "^" {
		t.ServiceCode = groups[4]
	}
	t.Discretionary = groups[5]
	return &t, nil
}

func NewIEC7813Track2(s string) (*IEC7813Track2, error) {
	var err error
	groups := iec7813Track2Re.FindStringSubmatch(s)
	if groups == nil {
		return nil, errors.New("invalid IEC 7813 track 2")
	}
	var t IEC7813Track2
	t.PAN, err = NewPAN(groups[1])
	if err != nil {
		return nil, err
	}
	if groups[1] != "=" {
		time, err := time.Parse(expDateLayout, groups[2])
		if err != nil {
			return nil, err
		}
		t.ExpDate = &time
	}
	if groups[3] != "=" {
		t.ServiceCode = groups[3]
	}
	t.Discretionary = groups[4]
	return &t, nil
}

func (t *IEC7813Tracks) Info() *Info {
	return NewInfo("IEC 7813 Tracks", "", t.Track1, t.Track2)
}

func NewIEC7813Tracks(t1, t2 string) (*IEC7813Tracks, error) {
	var t IEC7813Tracks
	var err error
	t.Track1, err = NewIEC7813Track1(t1)
	if err != nil {
		return nil, err
	}
	t.Track2, err = NewIEC7813Track2(t2)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
