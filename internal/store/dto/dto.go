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

const (
	V1 int = 1
	// V2 completions also carry the submitter's display name and submission
	// time. Older records decode with both fields zero-valued.
	V2 int = 2
)

type Challenge struct {
	Version     int
	Stage       Stage
	Weather     weather.Model
	Car         Car
	Completions []Completion
	Votes       []Vote
}

func FromChallenge(c challenge.Model) Challenge {
	completions := []Completion{}
	for _, completion := range c.Completions() {
		completions = append(completions, FromCompletion(completion))
	}

	votes := []Vote{}
	for _, vote := range c.Votes() {
		votes = append(votes, FromVote(vote))
	}

	return Challenge{
		Version:     V1,
		Stage:       FromStage(c.Stage()),
		Weather:     c.Weather(),
		Car:         FromCar(c.Car()),
		Completions: completions,
		Votes:       votes,
	}
}

func (dto Challenge) ToChallenge() challenge.Model {
	completions := []challenge.Completion{}
	for _, c := range dto.Completions {
		completions = append(completions, c.ToCompletion())
	}

	votes := []challenge.Vote{}
	for _, v := range dto.Votes {
		votes = append(votes, v.ToVote())
	}

	return challenge.NewChallenge(
		dto.Stage.ToStage(),
		dto.Weather,
		dto.Car.ToCar(),
		completions,
		votes,
	)
}

type Completion struct {
	Version     int
	UserID      string
	DisplayName string
	Duration    time.Duration
	SubmittedAt time.Time
}

func FromCompletion(c challenge.Completion) Completion {
	return Completion{
		Version:     V2,
		UserID:      c.UserID(),
		DisplayName: c.DisplayName(),
		Duration:    c.Duration(),
		SubmittedAt: c.SubmittedAt(),
	}
}

func (dto Completion) ToCompletion() challenge.Completion {
	return challenge.NewCompletionAt(dto.UserID, dto.DisplayName, dto.Duration, dto.SubmittedAt)
}

type Vote struct {
	Version   int
	UserID    string
	Sentiment challenge.Sentiment
}

func FromVote(v challenge.Vote) Vote {
	return Vote{
		Version:   V1,
		UserID:    v.UserID(),
		Sentiment: v.Sentiment(),
	}
}

func (dto Vote) ToVote() challenge.Vote {
	return challenge.NewVote(dto.UserID, dto.Sentiment)
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
