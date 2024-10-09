package game

type Model int

const (
	NotSet Model = iota
	DR2
	WRC
)

func (m Model) String() string {
	switch m {
	case DR2:
		return "Dirt Rally 2"
	case WRC:
		return "WRC"
	}
	return "invalid game"
}
