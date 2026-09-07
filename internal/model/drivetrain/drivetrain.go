package drivetrain

import (
	"slices"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
)

type Model int

const (
	FWD Model = iota
	AWD       // interchangable with 4WD - can't begin consts with numbers
	RWD
	AWDHYBRID
)

var byGame = map[game.Model][]Model{
	game.DR2: {FWD, AWD, RWD},
	game.WRC: {FWD, AWD, AWDHYBRID, RWD},
}

func List(g game.Model) []Model {
	return slices.Clone(byGame[g])
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
