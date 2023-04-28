package aamvad20

type EyeColor string

const (
	EyeBlack       EyeColor = "BLK"
	EyeBlue        EyeColor = "BLU"
	EyeBrown       EyeColor = "BRO"
	EyeDichromatic EyeColor = "DIC"
	EyeGray        EyeColor = "GRY"
	EyeGreen       EyeColor = "GRN"
	EyeHazel       EyeColor = "HAZ"
	EyePink        EyeColor = "PNK"
	EyeUnknown     EyeColor = "UNK"
)

func (ec EyeColor) String() string {
	switch ec {
	case EyeBlack:
		return "black"
	case EyeBlue:
		return "blue"
	case EyeBrown, "BRN":
		return "brown"
	case EyeDichromatic:
		return "dichromatic"
	case EyeGray:
		return "gray"
	case EyeGreen:
		return "green"
	case EyeHazel:
		return "hazel"
	case EyePink:
		return "pink"
	case EyeUnknown:
		return "unknown"
	}
	return string(ec)
}
