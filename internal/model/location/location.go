package location

import (
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
)

type Model int

const (
	ARG Model = iota
	AUS
	FIN
	DEU
	GRC
	MCO
	NZL
	POL
	SCO // not valid ISO3166, but not taken
	ESP
	SWE
	USA
	WAL // not valid ISO3166, but not taken
)

func List() []Model {
	return []Model{
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
	}
}

func (m Model) LongString() string {
	switch m {
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
	}
	return "invalid location"
}

func (m Model) String() string {
	switch m {
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
	}
	return "invalid location"
}

func (m Model) Flag() string {
	switch m {
	case ARG:
		return "🇦🇷"
	case AUS:
		return "🇦🇺"
	case FIN:
		return "🇫🇮"
	case DEU:
		return "🇩🇪"
	case GRC:
		return "🇬🇷"
	case MCO:
		return "🇲🇨"
	case NZL:
		return "🇳🇿"
	case POL:
		return "🇵🇱"
	case SCO:
		return "🏴󠁧󠁢󠁳󠁣󠁴󠁿"
	case ESP:
		return "🇪🇸"
	case SWE:
		return "🇸🇪"
	case USA:
		return "🇺🇸"
	case WAL:
		return "🏴󠁧󠁢󠁷󠁬󠁳󠁿"
	}
	return "invalid location"
}

func (m Model) Weather() []weather.Model {
	switch m {
	case ARG:
		return []weather.Model{weather.DRY, weather.WET}
	case AUS:
		return []weather.Model{weather.DRY, weather.WET}
	case FIN:
		return []weather.Model{weather.DRY, weather.WET}
	case DEU:
		return []weather.Model{weather.DRY, weather.WET}
	case GRC:
		return []weather.Model{weather.DRY, weather.WET}
	case MCO:
		return []weather.Model{weather.DRY}
	case NZL:
		return []weather.Model{weather.DRY, weather.WET}
	case POL:
		return []weather.Model{weather.DRY, weather.WET}
	case SCO:
		return []weather.Model{weather.DRY, weather.WET}
	case ESP:
		return []weather.Model{weather.DRY, weather.WET}
	case SWE:
		return []weather.Model{weather.SNOW}
	case USA:
		return []weather.Model{weather.DRY, weather.WET}
	case WAL:
		return []weather.Model{weather.DRY, weather.WET}
	}
	return []weather.Model{}
}
