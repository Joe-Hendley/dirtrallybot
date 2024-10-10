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

type Simple struct {
	randomSource *rand.Rand
	game         game.Model

	mu sync.Mutex
}

func NewSimple(game game.Model) *Simple {
	return &Simple{
		randomSource: rand.New(rand.NewPCG(0, 0)),
		game:         game,
		mu:           sync.Mutex{},
	}
}

func (s *Simple) Car() car.Model {
	s.mu.Lock()
	defer s.mu.Unlock()

	classes := class.List(s.game)
	class := classes[s.randomSource.IntN(len(classes))]

	cars := car.InClass(class, s.game)
	return cars[s.randomSource.IntN(len(cars))]
}

func (s *Simple) CarFromClass(class class.Model) car.Model {
	s.mu.Lock()
	defer s.mu.Unlock()

	cars := car.InClass(class, s.game)
	return cars[s.randomSource.IntN(len(cars))]
}

func (s *Simple) CarFromDrivetrain(drivetrain drivetrain.Model) car.Model {
	s.mu.Lock()
	defer s.mu.Unlock()

	classes := class.WithDrivetrain(drivetrain, s.game)
	class := classes[s.randomSource.IntN(len(classes))]

	cars := car.InClass(class, s.game)
	return cars[s.randomSource.IntN(len(cars))]
}

func (s *Simple) Loc() location.Model {
	s.mu.Lock()
	defer s.mu.Unlock()

	locs := location.List(s.game)
	return locs[s.randomSource.IntN(len(locs))]
}

func (s *Simple) Stage(location location.Model) stage.Model {
	s.mu.Lock()
	defer s.mu.Unlock()

	stages := stage.AtLocation(location)
	return stages[s.randomSource.IntN(len(stages))]
}

func (s *Simple) StageOfDistance(location location.Model, distance stage.Distance) stage.Model {
	s.mu.Lock()
	defer s.mu.Unlock()
	stages := stage.AtLocationWithDistance(location, distance)
	return stages[s.randomSource.IntN(len(stages))]
}

func (s *Simple) Weather(location location.Model) weather.Model {
	s.mu.Lock()
	defer s.mu.Unlock()

	weathers := location.Weather()
	return weathers[s.randomSource.IntN(len(weathers))]
}
