package aamvad20

type Sex byte

const (
	SexUnknown     Sex = '0' // unused in magstripe
	SexMale        Sex = '1'
	SexFemale      Sex = '2'
	SexUnspecified Sex = '9'
)

func (s Sex) String() string {
	switch s {
	case SexUnknown:
		return "unknown"
	case SexMale, 'M': // alpha for CA
		return "male"
	case SexFemale, 'F':
		return "female"
	case SexUnspecified, 'X':
		return "unspecified"
	}
	return string(s)
}
