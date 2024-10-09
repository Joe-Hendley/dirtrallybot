package class

import (
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
	}
}

func listWRC() []Model {
	return []Model{
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
	}
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

func WithDrivetrain(dt drivetrain.Model, g game.Model) []Model {
	switch g {
	case game.DR2:
		return withDrivetrainDR2(dt)
	case game.WRC:
		return withDrivetrainWRC(dt)
	}
	return []Model{}
}

func withDrivetrainDR2(dt drivetrain.Model) []Model {
	switch dt {
	case drivetrain.FWD:
		return []Model{H1, H2FWD, R2, F2}
	case drivetrain.AWD:
		return []Model{GroupB4WD, GroupA, NR4, WRC, R5}
	case drivetrain.RWD:
		return []Model{H2RWD, H3, GroupBRWD, RGT}
	}
	return []Model{}
}

func withDrivetrainWRC(dt drivetrain.Model) []Model {
	switch dt {
	case drivetrain.FWD:
		return []Model{H1_WRC, H2FWD_WRC, F2_WRC, S1600, Rally4}
	case drivetrain.AWD:
		return []Model{GroupB4WD_WRC, GroupA_WRC, NR4_WRC, S2000, Rally2, WRC1997to2011, WRC2017to2021, WRC2, JuniorWRC}
	case drivetrain.RWD:
		return []Model{H2RWD_WRC, H3RWD_WRC, GroupBRWD_WRC}
	case drivetrain.AWDHYBRID:
		return []Model{WRC_WRC}
	}
	return []Model{}
}
