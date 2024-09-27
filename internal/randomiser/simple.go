package randomiser

import (
	"math/rand/v2"
	"sync"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/drivetrain"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
)

type Simple struct {
	r *rand.Rand

	mu sync.Mutex
}

func NewSimple() *Simple {
	return &Simple{
		r:  rand.New(rand.NewPCG(0, 0)),
		mu: sync.Mutex{},
	}
}

func (s *Simple) Car() car.Model {
	s.mu.Lock()
	defer s.mu.Unlock()

	classes := class.List()
	class := classes[s.r.IntN(len(classes))]

	cars := car.InClass(class)
	return cars[s.r.IntN(len(cars))]
}

func (s *Simple) CarFromClass(class class.Model) car.Model {
	s.mu.Lock()
	defer s.mu.Unlock()

	cars := car.InClass(class)
	return cars[s.r.IntN(len(cars))]
}

func (s *Simple) CarFromDrivetrain(drivetrain drivetrain.Model) car.Model {
	s.mu.Lock()
	defer s.mu.Unlock()

	class.F2.Drivetrain()

	classes := class.List()
	class := classes[s.r.IntN(len(classes))]

	cars := car.InClass(class)
	return cars[s.r.IntN(len(cars))]
}

func (s *Simple) Loc() location.Model {
	s.mu.Lock()
	defer s.mu.Unlock()

	locs := location.List()
	return locs[s.r.IntN(len(locs))]
}

func (s *Simple) Stage(loc location.Model) stage.Model {
	s.mu.Lock()
	defer s.mu.Unlock()

	stages := stage.AtLocation(loc)
	return stages[s.r.IntN(len(stages))]
}

func (s *Simple) Weather(loc location.Model) weather.Model {
	s.mu.Lock()
	defer s.mu.Unlock()

	weathers := loc.Weather()
	return weathers[s.r.IntN(len(weathers))]
}
