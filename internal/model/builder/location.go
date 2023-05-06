package builder

import (
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
)

func Locations() map[location.Model]int {
	locations := make(map[location.Model]int)
	for _, l := range location.List() {
		locations[l] = 0
	}

	return locations
}
