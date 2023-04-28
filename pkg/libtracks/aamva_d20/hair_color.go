package aamvad20

type HairColor string

const (
	HairBald    HairColor = "BAL"
	HairBlack   HairColor = "BLK"
	HairBlond   HairColor = "BLN"
	HairBrown   HairColor = "BRO"
	HairGray    HairColor = "GRY"
	HairRed     HairColor = "RED"
	HairSandy   HairColor = "SDY"
	HairWhite   HairColor = "WHI"
	HairUnknown HairColor = "UNK"
)

func (hc HairColor) String() string {
	switch hc {
	case HairBald:
		return "bald"
	case HairBlack:
		return "black"
	case HairBlond:
		return "blond"
	case HairBrown, "BRN":
		return "brown"
	case HairGray:
		return "gray"
	case HairRed:
		return "red/auburn"
	case HairSandy:
		return "sandy"
	case HairWhite:
		return "white"
	case HairUnknown:
		return "unknown"
	}
	return string(hc)
}
