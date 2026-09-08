package randomiser

import (
	"math/rand/v2"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/drivetrain"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/popularity"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
)

// Biased picks stages, cars and weather with probability skewed by a popularity
// Snapshot: items with a positive net score come up more often, disliked ones
// less, but nothing is excluded. It is non-deterministic and safe for concurrent
// use, and is the default outside tests. An empty snapshot weights everything
// equally, so a fresh install behaves like Random until feedback arrives.
type Biased struct {
	game     game.Model
	snapshot popularity.Snapshot
}

// NewBiased returns a Biased randomiser for a game, biased by the given
// popularity snapshot.
func NewBiased(game game.Model, snapshot popularity.Snapshot) *Biased {
	return &Biased{game: game, snapshot: snapshot}
}

func (b *Biased) Car() car.Model {
	classes := class.List(b.game)
	chosen := classes[pick(b.snapshot, classes, popularity.ClassKey)]

	cars := car.InClass(chosen, b.game)
	return cars[pick(b.snapshot, cars, popularity.CarKey)]
}

func (b *Biased) CarFromClass(class class.Model) car.Model {
	cars := car.InClass(class, b.game)
	return cars[pick(b.snapshot, cars, popularity.CarKey)]
}

func (b *Biased) CarFromDrivetrain(drivetrain drivetrain.Model) car.Model {
	classes := class.WithDrivetrain(drivetrain, b.game)
	chosen := classes[pick(b.snapshot, classes, popularity.ClassKey)]

	cars := car.InClass(chosen, b.game)
	return cars[pick(b.snapshot, cars, popularity.CarKey)]
}

func (b *Biased) Loc() location.Model {
	locs := location.List(b.game)
	return locs[pick(b.snapshot, locs, popularity.LocationKey)]
}

func (b *Biased) Stage(loc location.Model) stage.Model {
	stages := stage.AtLocation(loc)
	return stages[pick(b.snapshot, stages, popularity.StageKey)]
}

func (b *Biased) StageOfDistance(loc location.Model, distance stage.Distance) stage.Model {
	stages := stage.AtLocationWithDistance(loc, distance)
	if len(stages) == 0 {
		// No stage of that length here; fall back to any stage at the location.
		stages = stage.AtLocation(loc)
	}

	return stages[pick(b.snapshot, stages, popularity.StageKey)]
}

func (b *Biased) Weather(loc location.Model) weather.Model {
	weathers := loc.Weather()
	return weathers[pick(b.snapshot, weathers, popularity.WeatherKey)]
}

// pick chooses an index into items with probability proportional to each item's
// popularity weight.
func pick[T any](snapshot popularity.Snapshot, items []T, key func(T) popularity.Key) int {
	weights := make([]float64, len(items))
	total := 0.0
	for i, item := range items {
		weights[i] = snapshot.Weight(key(item))
		total += weights[i]
	}

	target := rand.Float64() * total
	for i, w := range weights {
		target -= w
		if target < 0 {
			return i
		}
	}

	return len(weights) - 1
}
