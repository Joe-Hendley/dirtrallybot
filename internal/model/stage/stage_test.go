package stage_test

import (
	"testing"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/stretchr/testify/assert"
)

// Every location offered for a game must have at least one stage.
func TestEveryLocationHasStages(t *testing.T) {
	for _, g := range []game.Model{game.DR2, game.WRC} {
		for _, loc := range location.List(g) {
			stages := stage.AtLocation(loc)
			assert.NotEmptyf(t, stages, "%s has no stages", loc)

			for _, s := range stages {
				assert.Equalf(t, loc, s.Location(), "stage %q filed under the wrong location", s.Name())
			}
		}
	}
}

func TestAtLocationWithDistanceFilters(t *testing.T) {
	for _, s := range stage.AtLocationWithDistance(location.ARG, stage.Short) {
		assert.Equal(t, stage.Short, s.Distance())
	}
}
