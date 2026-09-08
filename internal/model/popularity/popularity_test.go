package popularity_test

import (
	"testing"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/drivetrain"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/popularity"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
	"github.com/stretchr/testify/assert"
)

func TestWeight(t *testing.T) {
	key := popularity.LocationKey(location.WAL)

	t.Run("no feedback sits at the midpoint", func(t *testing.T) {
		assert.InDelta(t, 0.5, popularity.Snapshot{}.Weight(key), 1e-9)
	})

	t.Run("net score moves the weight symmetrically", func(t *testing.T) {
		liked := popularity.Snapshot{key: {Up: 5}}.Weight(key)
		disliked := popularity.Snapshot{key: {Down: 5}}.Weight(key)

		assert.InDelta(t, 0.68, liked, 0.02)
		assert.InDelta(t, 0.32, disliked, 0.02)
		assert.InDelta(t, 1.0, liked+disliked, 1e-9)
	})

	t.Run("stays within the bounds at the extremes", func(t *testing.T) {
		loved := popularity.Snapshot{key: {Up: 1000}}.Weight(key)
		hated := popularity.Snapshot{key: {Down: 1000}}.Weight(key)

		assert.InDelta(t, 0.9, loved, 1e-9)
		assert.InDelta(t, 0.1, hated, 1e-9)
		assert.LessOrEqual(t, loved, 0.9)
		assert.GreaterOrEqual(t, hated, 0.1)
	})
}

func TestKeysFor(t *testing.T) {
	c := challenge.NewChallenge(
		stage.New("Sweet Lamb", location.WAL, stage.Long),
		weather.WET,
		car.New("Lancia Delta S4", class.GroupB4WD),
		nil,
		nil,
	)

	assert.ElementsMatch(t, []popularity.Key{
		popularity.StageKey(stage.New("Sweet Lamb", location.WAL, stage.Long)),
		popularity.LocationKey(location.WAL),
		popularity.DistanceKey(stage.Long),
		popularity.WeatherKey(weather.WET),
		popularity.CarKey(car.New("Lancia Delta S4", class.GroupB4WD)),
		popularity.ClassKey(class.GroupB4WD),
		popularity.DrivetrainKey(drivetrain.AWD),
	}, popularity.KeysFor(c))
}
