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
		return "Dry"
	case WET:
		return "Wet"
	case SNOW:
		return "Snow"
	}
	return "invalid weather"
}
