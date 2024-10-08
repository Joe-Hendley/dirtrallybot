package drivetrain

import "github.com/Joe-Hendley/dirtrallybot/internal/model/game"

type Model int

const (
	FWD Model = iota
	AWD       // interchangable with 4WD - can't begin consts with numbers
	RWD
	AWDHYBRID
)

func List(g game.Model) []Model {
	switch g {
	case game.DR2:
		return listDR2()
	case game.WRC:
		return listWRC()
	}
	return []Model{}
}

func listDR2() []Model {
	return []Model{
		FWD,
		AWD,
		RWD,
	}
}

func listWRC() []Model {
	return []Model{
		FWD,
		AWD,
		AWDHYBRID,
		RWD,
	}

}

func (m Model) String() string {
	switch m {
	case FWD:
		return "Front Wheel Drive"
	case AWD:
		return "Four Wheel Drive"
	case AWDHYBRID:
		return "Four Wheel Drive (Hybrid)"
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
	case AWDHYBRID:
		return "⚡"
	case RWD:
		return "🏎️"
	}

	return "invalid drivetrain"
}

func (m Model) FancyString() string {
	return m.Emoji() + " " + m.String()
}
