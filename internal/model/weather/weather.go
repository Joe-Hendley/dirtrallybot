package weather

type Model int

const (
	DRY Model = iota
	WET
	SNOW
)

func List() []Model {
	return []Model{
		DRY,
		WET,
		SNOW,
	}
}

func (m Model) String() string {
	switch m {
	case DRY:
		return "dry"
	case WET:
		return "wet"
	case SNOW:
		return "snow"
	}
	return "invalid weather"
}

func (m Model) Emoji() string {
	switch m {
	case DRY:
		return "☀️"
	case WET:
		return "💧"
	case SNOW:
		return "❄️"
	}
	return "invalid weather"
}
