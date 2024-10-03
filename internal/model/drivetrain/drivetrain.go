package drivetrain

type Model int

const (
	FWD Model = iota
	AWD       // interchangable with 4WD - can't begin consts with numbers
	RWD
)

func List() []Model {
	return []Model{
		FWD,
		AWD,
		RWD,
	}
}

func (m Model) String() string {
	switch m {
	case FWD:
		return "Front Wheel Drive"
	case AWD:
		return "Four Wheel Drive"
	case RWD:
		return "Rear Wheel Drive"
	}

	return "invalid drivetrain"
}

func (m Model) Emoji() string {
	switch m {
	case FWD:
		return "🚗"
	case AWD:
		return "🚙"
	case RWD:
		return "🏎️"
	}

	return "invalid drivetrain"
}

func (m Model) FancyString() string {
	return m.Emoji() + " " + m.String()
}
