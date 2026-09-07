package location

import (
	"slices"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
)

type Model int

// In theory this is based on ISO 3166
// should change it to actual words
const (
	ARG Model = iota
	AUS
	FIN
	DEU
	GRC
	MCO
	NZL
	POL
	SCO
	ESP
	SWE
	USA
	WAL

	MCO_WRC
	SWE_WRC
	MEX
	HRV
	PRT
	ITA
	KEN
	EST
	FIN_WRC
	GRC_WRC
	CHL
	CER // central european rally
	JPN
	MED // rally mediterraneo
	PAC // rally pacifico
	UUU // rally oceania
	SCA // rally scandia
	IBE // rally iberia
)

var byGame = map[game.Model][]Model{
	game.DR2: {
		ARG,
		AUS,
		FIN,
		DEU,
		GRC,
		MCO,
		NZL,
		POL,
		SCO,
		ESP,
		SWE,
		USA,
		WAL,
	},
	game.WRC: {
		MCO_WRC,
		SWE_WRC,
		MEX,
		HRV,
		PRT,
		ITA,
		KEN,
		EST,
		FIN_WRC,
		GRC_WRC,
		CHL,
		CER,
		JPN,
		MED,
		PAC,
		UUU,
		SCA,
		IBE,
	},
}

func List(g game.Model) []Model {
	return slices.Clone(byGame[g])
}

func (m Model) String() string {
	switch m {
	// DR2
	case ARG:
		return "Argentina"
	case AUS:
		return "Australia"
	case FIN:
		return "Finland"
	case DEU:
		return "Germany"
	case GRC:
		return "Greece"
	case MCO:
		return "Monaco"
	case NZL:
		return "New Zealand"
	case POL:
		return "Poland"
	case SCO:
		return "Scotland"
	case ESP:
		return "Spain"
	case SWE:
		return "Sweden"
	case USA:
		return "USA"
	case WAL:
		return "Wales"

	// WRC
	case MCO_WRC:
		return "Monte Carlo"
	case SWE_WRC:
		return "Sweden"
	case MEX:
		return "Mexico"
	case HRV:
		return "Croatia"
	case PRT:
		return "Portugal"
	case ITA:
		return "Italy"
	case KEN:
		return "Kenya"
	case EST:
		return "Estonia"
	case FIN_WRC:
		return "Finland"
	case GRC_WRC:
		return "Greece"
	case CHL:
		return "Chile"
	case CER:
		return "Central Europe"
	case JPN:
		return "Japan"
	case MED:
		return "Mediterranean"
	case PAC:
		return "Pacific"
	case UUU:
		return "Oceania"
	case SCA:
		return "Scandinavia"
	case IBE:
		return "Iberia"
	}
	return "invalid location"
}

func (m Model) DetailedString() string {
	switch m {
	// DR2
	case ARG:
		return "Catamarca Province, Argentina"
	case AUS:
		return "Monaro, Australia"
	case FIN:
		return "Jämsä, Finland"
	case DEU:
		return "Baumholder, Germany"
	case GRC:
		return "Argolis, Greece"
	case MCO:
		return "Monte Carlo, Monaco"
	case NZL:
		return "Hawkes Bay, New Zealand"
	case POL:
		return "Łęczna County, Poland"
	case SCO:
		return "Perth and Kinross, Scotland"
	case ESP:
		return "Ribadelles, Spain"
	case SWE:
		return "Värmland, Sweden"
	case USA:
		return "New England, USA"
	case WAL:
		return "Powys, Wales"

	// WRC
	case MCO_WRC:
		return "Rallye Monte Carlo"
	case SWE_WRC:
		return "Rally Sweden"
	case MEX:
		return "Guanajuato Rally México"
	case HRV:
		return "Croatia Rally"
	case PRT:
		return "Vodafone Rally de Portugal"
	case ITA:
		return "Rally Italia Sardegna"
	case KEN:
		return "Safari Rally Kenya"
	case EST:
		return "Rally Estonia"
	case FIN_WRC:
		return "Secto Rally Finland"
	case GRC_WRC:
		return "EKO Acropolis Rally Greece"
	case CHL:
		return "Bio Bío Rally Chile"
	case CER:
		return "Central Europe Rally"
	case JPN:
		return "Forum8 Rally Japan"
	case MED:
		return "Rally Mediterraneo (Fictional)"
	case PAC:
		return "Agon by AOC Rally Pacifico (Fictional)"
	case UUU:
		return "Fanatec Rally Oceania (Fictional)"
	case SCA:
		return "Rally Scandia (Fictional)"
	case IBE:
		return "Rally Iberia (Fictional)"
	}
	return "invalid location"
}

func (m Model) Weather() []weather.Model {
	switch m {
	case MCO:
		return []weather.Model{weather.DRY}
	case SWE, SWE_WRC, SCA:
		return []weather.Model{weather.SNOW}
	default:
		return []weather.Model{weather.DRY, weather.WET}
	}
}
