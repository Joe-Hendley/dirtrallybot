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

// ID is the short, URL-safe slug used in Discord component custom IDs.
func (m Model) ID() string {
	switch m {
	case DR2:
		return "dr2"
	case WRC:
		return "wrc"
	}
	return ""
}

// FromID is the inverse of Model.ID. It returns NotSet for an unknown slug.
func FromID(id string) Model {
	switch id {
	case DR2.ID():
		return DR2
	case WRC.ID():
		return WRC
	}
	return NotSet
}
