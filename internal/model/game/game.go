package game

type Model int

const (
	DR2 Model = iota
	WRC
)

func (m Model) String() string {
	switch m {
	case DR2:
		return "Dirt Rally 2"
	case WRC:
		return "WRC"
	}
}
