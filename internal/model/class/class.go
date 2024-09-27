package class

import (
	"github.com/Joe-Hendley/dirtrallybot/internal/model/drivetrain"
)

type Model int

const (
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
)

func List() []Model {
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

func (m Model) String() string {
	switch m {
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
	}

	return "invalid Model"
}

func (m Model) Drivetrain() drivetrain.Model {
	switch m {
	case H1, H2FWD, R2, F2:
		return drivetrain.FWD
	case GroupB4WD, GroupA, NR4, WRC, R5:
		return drivetrain.AWD
	case H2RWD, H3, GroupBRWD, RGT:
		return drivetrain.RWD
	}
	return 0 // equal to FWD, but it shouldn't matter
}

func WithDrivetrain(dt drivetrain.Model) []Model {
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
