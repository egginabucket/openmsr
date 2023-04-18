package cards

const future = " and other future industry assigments"

func MII(d int) string {
	switch d {
	case 0:
		return "ISO/TC 68 and other industry assignments"
	case 1:
		return "Airlines"
	case 2:
		return "Airlines, financial" + future
	case 3:
		return "Travel and entertainment"
	case 4, 5:
		return "Banking and financial"
	case 6:
		return "Merchandising and banking/financial"
	case 7:
		return "Petroleum" + future
	case 8:
		return "Healthcare, telecommunications" + future
	case 9:
		return "For assignment by national standards bodies"
	}
	return ""
}
