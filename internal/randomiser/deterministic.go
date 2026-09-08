package randomiser

import (
	"math/rand/v2"
	"sync"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/drivetrain"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
)

// Deterministic picks stages, cars and weather with a fixed PCG seed, so a given
// sequence of calls always produces the same challenges. This is intentional:
// challenges are meant to be reproducible.
type Deterministic struct {
	randomSource *rand.Rand
	game         game.Model

	mu sync.Mutex
}

// NewDeterministic returns a Deterministic randomiser for a game, seeded with a
// fixed value so its output is reproducible across restarts.
func NewDeterministic(game game.Model) *Deterministic {
	return &Deterministic{
		randomSource: rand.New(rand.NewPCG(0, 0)),
		game:         game,
	}
}

func (s *Deterministic) Car() car.Model {
	s.mu.Lock()
	defer s.mu.Unlock()

	classes := class.List(s.game)
	class := classes[s.randomSource.IntN(len(classes))]

	cars := car.InClass(class, s.game)
	return cars[s.randomSource.IntN(len(cars))]
}

func (s *Deterministic) CarFromClass(class class.Model) car.Model {
	s.mu.Lock()
	defer s.mu.Unlock()

	cars := car.InClass(class, s.game)
	return cars[s.randomSource.IntN(len(cars))]
}

func (s *Deterministic) CarFromDrivetrain(drivetrain drivetrain.Model) car.Model {
	s.mu.Lock()
	defer s.mu.Unlock()

	classes := class.WithDrivetrain(drivetrain, s.game)
	class := classes[s.randomSource.IntN(len(classes))]

	cars := car.InClass(class, s.game)
	return cars[s.randomSource.IntN(len(cars))]
}

func (s *Deterministic) Loc() location.Model {
	s.mu.Lock()
	defer s.mu.Unlock()

	locs := location.List(s.game)
	return locs[s.randomSource.IntN(len(locs))]
}

func (s *Deterministic) Stage(location location.Model) stage.Model {
	s.mu.Lock()
	defer s.mu.Unlock()

	stages := stage.AtLocation(location)
	return stages[s.randomSource.IntN(len(stages))]
}

func (s *Deterministic) StageOfDistance(location location.Model, distance stage.Distance) stage.Model {
	s.mu.Lock()
	defer s.mu.Unlock()

	stages := stage.AtLocationWithDistance(location, distance)
	if len(stages) == 0 {
		// No stage of that length here; fall back to any stage at the location.
		stages = stage.AtLocation(location)
	}

	return stages[s.randomSource.IntN(len(stages))]
}

func (s *Deterministic) Weather(location location.Model) weather.Model {
	s.mu.Lock()
	defer s.mu.Unlock()

	weathers := location.Weather()
	return weathers[s.randomSource.IntN(len(weathers))]
}
