package libtracks

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	aamvad20 "github.com/egginabucket/openmsr/pkg/libtracks/aamva_d20"
)

type (
	AAMVATrack1 struct {
		StateOrProv string //[2]byte
		City        string
		Name        []string
		Address     []string
	}
	AAMVATrack2 struct {
		IIN           byte
		ID            string // includes overflow
		ExpDate       time.Time
		NonExp        bool
		ExpBirthMonth bool
		ExpBirthDay   bool
		BirthDate     time.Time
	}
	AAMVATrack3 struct {
		CDSVersion    int
		JdxVersion    byte
		PostalCode    string //[11]byte
		Class         string //[2]byte
		Restrictions  string //[10]byte
		Endorsements  string //[4]byte
		Sex           aamvad20.Sex
		Height        *aamvad20.Height
		Weight        int
		HairColor     aamvad20.HairColor //[3]byte
		EyeColor      aamvad20.EyeColor  //[3]byte
		Discretionary string
	}
	AAMVATracks struct {
		Track1 *AAMVATrack1
		Track2 *AAMVATrack2
		Track3 *AAMVATrack3
	}
)

const (
	expMonthNonExp   = "77"
	expMonthBirth    = "88"
	expMonthBirthDay = "99"
	birthDateLayout  = "20060102"
)

// ^%([A-Z]{2})([^\^]{1,12}\^|[^\^]{13})([^\^]{1,34}\^|[^\^]{35})([^\^]{1,28}\^|[^\^]{29})\?$
var AAMVATrack1Re = regexp.MustCompile(`^%([A-Z]{2})([^\^]{2,12}\^|[^\^]{13})([^\^]{2,34}\^|[^\^]{35})([^\^]{2,28}\^|[^\^]{29,})\?$`)
var AAMVATrack2Re = regexp.MustCompile(`^;([0-9])([0-9]{1,13})=([0-9]{4})([0-9]{8})([0-9]{1,5}|=)\?$`)
var AAMVATrack3Re = regexp.MustCompile(`^%([0-9])(.)([A-Z0-9 ]{11})([A-Z0-9 ]{2})([A-Z0-9 ]{10})([A-Z0-9 ]{4})([A-Z0-9])([0-9]{3}|   )([0-9 ]{3}|   )([A-Z]{3}|   )([A-Z]{3}|   )([^\?]+)\?$`)

func (t *AAMVATrack1) String() string {
	var b strings.Builder
	b.WriteByte('%')
	b.WriteString(t.StateOrProv)
	b.WriteString(t.City)
	if len(t.City) < 13 {
		b.WriteByte('^')
	}
	name := strings.Join(t.Name, "$")
	b.WriteString(name)
	if len(name) < 35 {
		b.WriteByte('^')
	}
	address := strings.Join(t.Address, "$")
	b.WriteString(address)
	if len(address) < 29 {
		b.WriteByte('^')
	}
	b.WriteByte('?')
	return b.String()
}

func (t *AAMVATrack1) Info() *Info {
	if t == nil {
		return NewInfo("Track 1", "nil")
	}
	return NewInfo("Track 1", t.String(),
		NewInfo("State or province", t.StateOrProv),
		NewInfo("City", t.City),
		NewInfo("Name", strings.Join(t.Name, " / ")),
		NewInfo("Address", strings.Join(t.Address, " / ")),
	)
}

func (t *AAMVATrack2) String() string {
	var b strings.Builder
	b.WriteByte(';')
	b.WriteByte(t.IIN)
	if len(t.ID) > 13 {
		b.WriteString(t.ID[:13])
	} else {
		b.WriteString(t.ID)
	}
	b.WriteByte('=')
	switch {
	case t.NonExp, t.ExpBirthMonth, t.ExpBirthDay:
		b.WriteString(t.ExpDate.Format("06"))
		switch {
		case t.NonExp:
			b.WriteString(expMonthNonExp)
		case t.ExpBirthMonth:
			b.WriteString(expMonthBirth)
		case t.ExpBirthDay:
			b.WriteString(expMonthBirthDay)
		}
	default:
		b.WriteString(t.ExpDate.Format(expDateLayout))
	}
	b.WriteString(t.BirthDate.Format(birthDateLayout))
	if len(t.ID) > 13 { // overflow field
		b.WriteString(t.ID[13:])
	} else {
		b.WriteByte('=')
	}
	b.WriteByte('?')
	return b.String()
}

func (t *AAMVATrack2) Info() *Info {
	if t == nil {
		return NewInfo("Track 2", "nil")
	}
	return NewInfo("Track 2", t.String(),
		NewInfo("IIN", string(t.IIN)),
		NewInfo("DL/ID", t.ID),
		NewInfo("Expiration date", t.ExpDate.Format("2006-01")),
		NewInfo("Non-Expiring", t.NonExp),
		NewInfo("Expires on birth month", t.ExpBirthMonth),
		NewInfo("Expires one month after birth month", t.ExpBirthDay),
		NewInfo("Birth date", t.BirthDate.Format("2006-01-02")),
	)
}

func (t *AAMVATrack3) String() string {
	var b strings.Builder
	b.WriteByte('%')
	fmt.Fprintf(&b, "%d%c%-11s", t.CDSVersion, t.JdxVersion, t.PostalCode)
	b.WriteString(t.Class)
	b.WriteString(t.Restrictions)
	b.WriteString(t.Endorsements)
	b.WriteByte(byte(t.Sex))
	if t.Height != nil {
		t.Height.Write3Digits(&b)
	} else {
		b.WriteString("   ")
	}
	if t.Weight != 0 {
		fmt.Fprintf(&b, "% 3d", t.Weight)
	} else {
		b.WriteString("   ")
	}
	if t.HairColor != "" {
		b.WriteString(string(t.HairColor))
	} else {
		b.WriteString("   ")
	}
	if t.EyeColor != "" {
		b.WriteString(string(t.EyeColor))
	} else {
		b.WriteString("   ")
	}
	b.WriteString(t.Discretionary)
	b.WriteByte('?')
	return b.String()
}

func (t *AAMVATrack3) Info() *Info {
	if t == nil {
		return NewInfo("Track 3", "nil")
	}
	return NewInfo("Track 3", t.String(),
		NewInfo("CDS version", t.CDSVersion),
		NewInfo("Jurisdiction version", string(t.JdxVersion)),
		NewInfo("Postal code", t.PostalCode),
		NewInfo("Class", t.Class),
		NewInfo("Restrictions", t.Restrictions),
		NewInfo("Sex", t.Sex.String()),
		NewInfo("Height", t.Height),
		NewInfo("Weight", t.Weight),
		NewInfo("Hair color", t.HairColor.String()),
		NewInfo("Eye color", t.EyeColor.String()),
		NewInfo("Discretionary data", t.Discretionary),
	)
}

func NewAAMVATrack1(s string) (*AAMVATrack1, error) {
	groups := AAMVATrack1Re.FindStringSubmatch(s)
	if groups == nil {
		return nil, errors.New("invalid AAMVA track 1")
	}
	var t AAMVATrack1
	t.StateOrProv = groups[1]
	t.City = strings.TrimSuffix(groups[2], "^")
	t.Name = strings.Split(strings.TrimSuffix(groups[3], "^"), "$")
	t.Address = strings.Split(strings.TrimSuffix(groups[4], "^"), "$")
	return &t, nil
}

func NewAAMVATrack2(s string) (*AAMVATrack2, error) {
	var err error
	groups := AAMVATrack2Re.FindStringSubmatch(s)
	if groups == nil {
		return nil, errors.New("invalid AAMVA track 2")
	}
	var t AAMVATrack2
	t.IIN = groups[1][0]
	t.ID = groups[2]
	if groups[5] != "=" {
		t.ID += groups[4]
	}
	switch groups[3][2:] {
	case "77", "88", "99":
		t.ExpDate, err = time.Parse("06", groups[2][:2])
		if err != nil {
			return nil, err
		}
		switch groups[3][2:] {
		case expMonthNonExp:
			t.NonExp = true
		case expMonthBirth:
			t.ExpBirthMonth = true
		case expMonthBirthDay:
			t.ExpBirthDay = true
		}
	default:
		t.ExpDate, err = time.Parse(expDateLayout, groups[3])
		if err != nil {
			return nil, err
		}
	}
	t.BirthDate, err = time.Parse(birthDateLayout, groups[4])
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func NewAAMVATrack3(s string) (*AAMVATrack3, error) {
	var err error
	groups := AAMVATrack3Re.FindStringSubmatch(s)
	if groups == nil {
		return nil, errors.New("invalid AAMVA track 3")
	}
	var t AAMVATrack3
	t.CDSVersion = int(groups[1][0] - '0')
	t.JdxVersion = groups[2][0]
	t.PostalCode = strings.TrimRight(groups[3], " ")
	t.Class = groups[4]
	t.Restrictions = groups[5]
	t.Endorsements = groups[6]
	t.Sex = aamvad20.Sex(groups[7][0])
	if groups[8] != "   " {
		t.Height = aamvad20.ParseHeight(groups[8])
	}
	if groups[9] != "   " {
		t.Weight, err = strconv.Atoi(strings.TrimLeft(groups[9], " "))
		if err != nil {
			return nil, err
		}
	}
	if groups[10] != "   " {
		t.HairColor = aamvad20.HairColor(groups[10])
	}
	if groups[11] != "   " {
		t.EyeColor = aamvad20.EyeColor(groups[11])
	}
	t.Discretionary = groups[12]
	return &t, nil
}

func (t *AAMVATracks) Info() *Info {
	return NewInfo("AAMVA Tracks", "", t.Track1, t.Track2, t.Track3)
}

func NewAAMVATracks(t1, t2, t3 string) (*AAMVATracks, error) {
	var t AAMVATracks
	var err error
	t.Track1, err = NewAAMVATrack1(t1)
	if err != nil {
		return nil, err
	}
	t.Track2, err = NewAAMVATrack2(t2)
	if err != nil {
		return nil, err
	}
	t.Track3, err = NewAAMVATrack3(t3)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
