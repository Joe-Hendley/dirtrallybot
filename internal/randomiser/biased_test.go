package randomiser_test

import (
	"testing"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/popularity"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/randomiser"
	"github.com/stretchr/testify/assert"
)

// A challenge for one game must never be built from another game's content.
func TestBiasedStaysWithinItsGame(t *testing.T) {
	for _, g := range []game.Model{game.DR2, game.WRC} {
		r := randomiser.NewBiased(g, popularity.Snapshot{})
		assert.Containsf(t, location.List(g), r.Loc(), "%s randomiser produced a location from another game", g)
	}
}

func TestBiasedWithoutFeedbackSpreadsOut(t *testing.T) {
	r := randomiser.NewBiased(game.DR2, popularity.Snapshot{})

	seen := map[location.Model]struct{}{}
	for range 200 {
		seen[r.Loc()] = struct{}{}
	}

	assert.Greater(t, len(seen), 1)
}

func TestBiasedLeansTowardsPopularWithoutExcludingOthers(t *testing.T) {
	favourite := location.WAL
	snapshot := popularity.Snapshot{
		popularity.LocationKey(favourite): {Up: 40},
	}
	r := randomiser.NewBiased(game.DR2, snapshot)

	const draws = 3000
	otherLocations := len(location.List(game.DR2)) - 1

	favouriteHits, otherHits := 0, 0
	for range draws {
		if r.Loc() == favourite {
			favouriteHits++
			continue
		}
		otherHits++
	}

	// Favourite weight ~0.9 vs ~0.5 for every other location: it should come up
	// well above its uniform share and beat any single other location...
	assert.Greater(t, favouriteHits, draws/(otherLocations+1))
	assert.Greater(t, favouriteHits, otherHits/otherLocations)
	// ...but the field as a whole still dominates.
	assert.Greater(t, otherHits, favouriteHits*2)
}

func TestBiasedStageOfDistanceHonoursDistance(t *testing.T) {
	r := randomiser.NewBiased(game.WRC, popularity.Snapshot{})
	loc := location.List(game.WRC)[0]

	for _, distance := range []stage.Distance{stage.Short, stage.Long, stage.ReallyLong} {
		for range 50 {
			assert.Equal(t, distance, r.StageOfDistance(loc, distance).Distance())
		}
	}
}

func TestBiasedStageOfDistanceFallsBackWhenNoneMatch(t *testing.T) {
	r := randomiser.NewBiased(game.DR2, popularity.Snapshot{})
	loc := location.List(game.DR2)[0]

	assert.Equal(t, loc, r.StageOfDistance(loc, stage.ReallyLong).Location())
}
