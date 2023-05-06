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
		return "front wheel drive"
	case AWD:
		return "four wheel drive"
	case RWD:
		return "rear wheel drive"
	}

	return "invalid drivetrain"
}
