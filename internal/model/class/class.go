package class

import (
	"slices"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/drivetrain"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
)

type Model int

const (
	// DR2
	H1 Model = iota
	H2FWD
	H2RWD
	H3
	GroupBRWD
	GroupB4WD
	R2
	F2
	GroupA
	NR4
	WRC
	R5
	RGT

	// WRC
	WRC_WRC
	WRC2
	JuniorWRC
	WRC2017to2021
	WRC1997to2011
	Rally2
	Rally4
	NR4_WRC
	S2000
	S1600
	F2_WRC
	GroupA_WRC
	GroupB4WD_WRC
	GroupBRWD_WRC
	H3RWD_WRC
	H2RWD_WRC
	H2FWD_WRC
	H1_WRC
)

var byGame = map[game.Model][]Model{
	game.DR2: {
		H1,
		H2FWD,
		H2RWD,
		H3,
		GroupBRWD,
		GroupB4WD,
		R2,
		F2,
		GroupA,
		NR4,
		WRC,
		R5,
		RGT,
	},
	game.WRC: {
		WRC_WRC,
		WRC2,
		JuniorWRC,
		WRC2017to2021,
		WRC1997to2011,
		Rally2,
		Rally4,
		NR4_WRC,
		S2000,
		S1600,
		F2_WRC,
		GroupA_WRC,
		GroupB4WD_WRC,
		GroupBRWD_WRC,
		H3RWD_WRC,
		H2RWD_WRC,
		H2FWD_WRC,
		H1_WRC,
	},
}

func List(g game.Model) []Model {
	return slices.Clone(byGame[g])
}

func (m Model) String() string {
	switch m {
	// DR2
	case H1:
		return "H1 (FWD)"
	case H2FWD:
		return "H2 (FWD)"
	case H2RWD:
		return "H2 (RWD)"
	case H3:
		return "H3 (RWD)"
	case GroupBRWD:
		return "Group B (RWD)"
	case GroupB4WD:
		return "Group B (4WD)"
	case R2:
		return "R2"
	case F2:
		return "F2 Kit Car"
	case GroupA:
		return "Group A"
	case NR4:
		return "NR4/R4"
	case WRC:
		return "Up to 2000cc"
	case R5:
		return "R5"
	case RGT:
		return "Rally GT"

	// WRC
	case WRC_WRC:
		return "WRC"
	case WRC2:
		return "WRC2 * Cars"
	case JuniorWRC:
		return "Junior WRC"
	case WRC2017to2021:
		return "World Rally Cars 2017-2021"
	case WRC1997to2011:
		return "World Rally Cars 1997-2011"
	case Rally2:
		return "Rally2 *"
	case Rally4:
		return "Rally4 Cars"
	case NR4_WRC:
		return "NR4/R4"
	case S2000:
		return "S2000"
	case S1600:
		return "S1600"
	case F2_WRC:
		return "F2 Kit Cars"
	case GroupA_WRC:
		return "Group A"
	case GroupB4WD_WRC:
		return "Group B(4WD)"
	case GroupBRWD_WRC:
		return "Group B (RWD)"
	case H3RWD_WRC:
		return "H3 (RWD)"
	case H2RWD_WRC:
		return "H2 (RWD)"
	case H2FWD_WRC:
		return "H2 (FWD)"
	case H1_WRC:
		return "H1 (FWD)"
	}

	return "invalid Model"
}

func (m Model) Drivetrain() drivetrain.Model {
	switch m {
	// DR2
	case H1, H2FWD, R2, F2:
		return drivetrain.FWD
	case GroupB4WD, GroupA, NR4, WRC, R5:
		return drivetrain.AWD
	case H2RWD, H3, GroupBRWD, RGT:
		return drivetrain.RWD

	// WRC
	case WRC_WRC:
		return drivetrain.AWDHYBRID
	case H1_WRC, H2FWD_WRC, F2_WRC, S1600, Rally4:
		return drivetrain.FWD
	case GroupB4WD_WRC, GroupA_WRC, S2000, NR4_WRC, Rally2, WRC1997to2011, WRC2017to2021, WRC2, JuniorWRC:
		return drivetrain.AWD
	case H2RWD_WRC, H3RWD_WRC, GroupBRWD_WRC:
		return drivetrain.RWD
	}
	return 0 // equal to FWD, but it shouldn't matter
}

var byGameDrivetrain = map[game.Model]map[drivetrain.Model][]Model{
	game.DR2: {
		drivetrain.FWD: {H1, H2FWD, R2, F2},
		drivetrain.AWD: {GroupB4WD, GroupA, NR4, WRC, R5},
		drivetrain.RWD: {H2RWD, H3, GroupBRWD, RGT},
	},
	game.WRC: {
		drivetrain.FWD:       {H1_WRC, H2FWD_WRC, F2_WRC, S1600, Rally4},
		drivetrain.AWD:       {GroupB4WD_WRC, GroupA_WRC, NR4_WRC, S2000, Rally2, WRC1997to2011, WRC2017to2021, WRC2, JuniorWRC},
		drivetrain.RWD:       {H2RWD_WRC, H3RWD_WRC, GroupBRWD_WRC},
		drivetrain.AWDHYBRID: {WRC_WRC},
	},
}

func WithDrivetrain(dt drivetrain.Model, g game.Model) []Model {
	return slices.Clone(byGameDrivetrain[g][dt])
}
