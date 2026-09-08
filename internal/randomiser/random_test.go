package randomiser_test

import (
	"testing"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/randomiser"
	"github.com/stretchr/testify/assert"
)

// Unlike Deterministic, Random should not settle on a single choice; over many
// draws it must produce more than one location.
func TestRandomVariesAcrossDraws(t *testing.T) {
	r := randomiser.NewRandom(game.DR2)

	seen := map[location.Model]struct{}{}
	for range 100 {
		seen[r.Loc()] = struct{}{}
	}

	assert.Greater(t, len(seen), 1)
}

func TestRandomStageOfDistanceHonoursDistance(t *testing.T) {
	r := randomiser.NewRandom(game.WRC)
	loc := location.List(game.WRC)[0]

	for _, distance := range []stage.Distance{stage.Short, stage.Long, stage.ReallyLong} {
		for range 50 {
			assert.Equal(t, distance, r.StageOfDistance(loc, distance).Distance())
		}
	}
}

func TestRandomStageOfDistanceFallsBackWhenNoneMatch(t *testing.T) {
	r := randomiser.NewRandom(game.DR2)
	loc := location.List(game.DR2)[0]

	// DR2 has no ReallyLong stages, so this must fall back to any stage at the
	// location rather than panic on an empty slice.
	assert.Equal(t, loc, r.StageOfDistance(loc, stage.ReallyLong).Location())
}
