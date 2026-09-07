package challenge

import (
	"slices"
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/drivetrain"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
)

type Randomiser interface {
	Car() car.Model
	CarFromClass(class class.Model) car.Model
	CarFromDrivetrain(drivetrain drivetrain.Model) car.Model

	Loc() location.Model
	Weather(loc location.Model) weather.Model
	Stage(loc location.Model) stage.Model
	StageOfDistance(loc location.Model, distance stage.Distance) stage.Model
}

type Model struct {
	stage       stage.Model
	weather     weather.Model
	car         car.Model
	completions []Completion
}

func NewChallenge(s stage.Model, w weather.Model, car car.Model, completions []Completion) Model {
	return Model{
		stage:       s,
		weather:     w,
		car:         car,
		completions: completions,
	}
}

func NewRandomChallenge(c Config, r Randomiser) Model {
	var loc location.Model

	challenge := Model{}

	if c.Location != nil {
		loc = *c.Location
	} else {
		loc = r.Loc()
	}

	switch {
	case c.Stage != nil:
		challenge.stage = *c.Stage
	case c.Distance != nil:
		challenge.stage = r.StageOfDistance(loc, *c.Distance)
	default:
		challenge.stage = r.Stage(loc)
	}

	if c.Weather != nil {
		challenge.weather = *c.Weather
	} else {
		challenge.weather = r.Weather(loc)
	}

	switch {
	case c.Car != nil:
		challenge.car = *c.Car
		return challenge

	case c.Class != nil:
		challenge.car = r.CarFromClass(*c.Class)
		return challenge

	case c.Drivetrain != nil:
		challenge.car = r.CarFromDrivetrain(*c.Drivetrain)
		return challenge
	}

	challenge.car = r.Car()

	challenge.completions = []Completion{}

	return challenge
}

func (m *Model) Stage() stage.Model {
	return m.stage
}

func (m *Model) Weather() weather.Model {
	return m.weather
}

func (m *Model) Car() car.Model {
	return m.car
}

func (m *Model) Completions() []Completion {
	return m.completions
}

func (m *Model) RegisterCompletion(c Completion) {
	m.completions = append(m.completions, c)
}

func (m *Model) TopThree() []Completion {
	if len(m.completions) < 2 {
		return m.completions
	}

	sorted := make([]Completion, len(m.completions))
	copy(sorted, m.completions)

	slices.SortFunc(sorted, func(a, b Completion) int { return int(a.duration - b.duration) })

	topThree := make([]Completion, 0, 3)
	listedUsers := map[string]struct{}{}
	for _, completion := range sorted {
		_, ok := listedUsers[completion.userID]
		if !ok {
			topThree = append(topThree, completion)
			listedUsers[completion.userID] = struct{}{}
		}

		if len(topThree) == 3 {
			break
		}
	}

	return topThree
}

func (m *Model) UserCompletions() map[string][]time.Duration {
	if len(m.completions) == 0 {
		return map[string][]time.Duration{}
	}

	userCompletions := make(map[string][]time.Duration)

	for _, completion := range m.completions {
		userCompletions[completion.userID] = append(userCompletions[completion.userID], completion.duration)
	}

	return userCompletions
}

type Config struct {
	Game game.Model

	Location *location.Model
	Distance *stage.Distance
	Stage    *stage.Model
	Weather  *weather.Model

	Car        *car.Model
	Class      *class.Model
	Drivetrain *drivetrain.Model
}

type Completion struct {
	userID   string
	duration time.Duration
}

func NewCompletion(userID string, duration time.Duration) Completion {
	return Completion{
		userID:   userID,
		duration: duration,
	}
}

func (c Completion) UserID() string {
	return c.userID
}

func (c Completion) Duration() time.Duration {
	return c.duration
}
