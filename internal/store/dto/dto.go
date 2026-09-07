package dto

import (
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
)

const V1 int = 1

type Challenge struct {
	Version     int
	Stage       Stage
	Weather     weather.Model
	Car         Car
	Completions []Completion
}

func FromChallenge(c challenge.Model) Challenge {
	completions := []Completion{}
	for _, completion := range c.Completions() {
		completions = append(completions, FromCompletion(completion))
	}

	return Challenge{
		Version:     V1,
		Stage:       FromStage(c.Stage()),
		Weather:     c.Weather(),
		Car:         FromCar(c.Car()),
		Completions: completions,
	}
}

func (dto Challenge) ToChallenge() challenge.Model {
	completions := []challenge.Completion{}
	for _, c := range dto.Completions {
		completions = append(completions, c.ToCompletion())
	}

	return challenge.NewChallenge(
		dto.Stage.ToStage(),
		dto.Weather,
		dto.Car.ToCar(),
		completions,
	)
}

type Completion struct {
	Version  int
	UserID   string
	Duration time.Duration
}

func FromCompletion(c challenge.Completion) Completion {
	return Completion{
		Version:  V1,
		UserID:   c.UserID(),
		Duration: c.Duration(),
	}
}

func (dto Completion) ToCompletion() challenge.Completion {
	return challenge.NewCompletion(dto.UserID, dto.Duration)
}

type Stage struct {
	Version  int
	Name     string
	Location location.Model
	Distance stage.Distance
}

func FromStage(s stage.Model) Stage {
	return Stage{
		Version:  V1,
		Name:     s.Name(),
		Location: s.Location(),
		Distance: s.Distance(),
	}
}

func (dto Stage) ToStage() stage.Model {
	return stage.New(
		dto.Name,
		dto.Location,
		dto.Distance,
	)
}

type Car struct {
	Version int
	Name    string
	Class   class.Model
}

func FromCar(c car.Model) Car {
	return Car{
		Version: V1,
		Name:    c.Name(),
		Class:   c.Class(),
	}
}

func (dto Car) ToCar() car.Model {
	return car.New(dto.Name, dto.Class)
}
