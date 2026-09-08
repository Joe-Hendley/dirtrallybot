package randomiser

import (
	"math/rand/v2"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/drivetrain"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
)

// Random picks stages, cars and weather from the package-global math/rand/v2
// source, so challenges differ from run to run. It is the default outside tests;
// use Deterministic where reproducible output matters. The global source is safe
// for concurrent use, so Random needs no locking of its own.
type Random struct {
	game game.Model
}

// NewRandom returns a Random randomiser for a game.
func NewRandom(game game.Model) *Random {
	return &Random{game: game}
}

func (r *Random) Car() car.Model {
	classes := class.List(r.game)
	chosen := classes[rand.IntN(len(classes))]

	cars := car.InClass(chosen, r.game)
	return cars[rand.IntN(len(cars))]
}

func (r *Random) CarFromClass(class class.Model) car.Model {
	cars := car.InClass(class, r.game)
	return cars[rand.IntN(len(cars))]
}

func (r *Random) CarFromDrivetrain(drivetrain drivetrain.Model) car.Model {
	classes := class.WithDrivetrain(drivetrain, r.game)
	chosen := classes[rand.IntN(len(classes))]

	cars := car.InClass(chosen, r.game)
	return cars[rand.IntN(len(cars))]
}

func (r *Random) Loc() location.Model {
	locs := location.List(r.game)
	return locs[rand.IntN(len(locs))]
}

func (r *Random) Stage(loc location.Model) stage.Model {
	stages := stage.AtLocation(loc)
	return stages[rand.IntN(len(stages))]
}

func (r *Random) StageOfDistance(loc location.Model, distance stage.Distance) stage.Model {
	stages := stage.AtLocationWithDistance(loc, distance)
	if len(stages) == 0 {
		// No stage of that length here; fall back to any stage at the location.
		stages = stage.AtLocation(loc)
	}

	return stages[rand.IntN(len(stages))]
}

func (r *Random) Weather(loc location.Model) weather.Model {
	weathers := loc.Weather()
	return weathers[rand.IntN(len(weathers))]
}
