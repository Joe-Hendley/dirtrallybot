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
	votes       []Vote
}

func NewChallenge(s stage.Model, w weather.Model, car car.Model, completions []Completion, votes []Vote) Model {
	return Model{
		stage:       s,
		weather:     w,
		car:         car,
		completions: completions,
		votes:       votes,
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
	challenge.votes = []Vote{}

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
	userID      string
	displayName string
	duration    time.Duration
	submittedAt time.Time
}

// NewCompletion records a completion known only by user ID and time. It is used
// for legacy data and where the submitter's name and submission time do not
// matter.
func NewCompletion(userID string, duration time.Duration) Completion {
	return Completion{
		userID:   userID,
		duration: duration,
	}
}

// NewCompletionAt records a completion together with the submitter's display
// name as it was at submission and the moment it was submitted.
func NewCompletionAt(userID, displayName string, duration time.Duration, submittedAt time.Time) Completion {
	return Completion{
		userID:      userID,
		displayName: displayName,
		duration:    duration,
		submittedAt: submittedAt,
	}
}

func (c Completion) UserID() string {
	return c.userID
}

// DisplayName is the submitter's display name captured at submission, or "" for
// a completion recorded without one.
func (c Completion) DisplayName() string {
	return c.displayName
}

func (c Completion) Duration() time.Duration {
	return c.duration
}

// SubmittedAt is when the completion was submitted, or the zero time for a
// completion recorded without one.
func (c Completion) SubmittedAt() time.Time {
	return c.submittedAt
}

// Sentiment is a thumbs-up or thumbs-down on a challenge. The zero value is
// invalid so an unset Vote is never mistaken for feedback.
type Sentiment int

const (
	Up Sentiment = iota + 1
	Down
)

// Vote is one user's feedback on a challenge. A user has at most one vote per
// challenge.
type Vote struct {
	userID    string
	sentiment Sentiment
}

func NewVote(userID string, sentiment Sentiment) Vote {
	return Vote{
		userID:    userID,
		sentiment: sentiment,
	}
}

func (v Vote) UserID() string {
	return v.userID
}

func (v Vote) Sentiment() Sentiment {
	return v.sentiment
}

// Votes returns the current votes on the challenge, one per user.
func (m *Model) Votes() []Vote {
	return m.votes
}

// VoteFor returns the user's current vote, or false if they have not voted.
func (m *Model) VoteFor(userID string) (Vote, bool) {
	for _, v := range m.votes {
		if v.userID == userID {
			return v, true
		}
	}
	return Vote{}, false
}

// SetVote records a user's vote, replacing any existing one from that user.
func (m *Model) SetVote(vote Vote) {
	for i, v := range m.votes {
		if v.userID == vote.userID {
			m.votes[i] = vote
			return
		}
	}
	m.votes = append(m.votes, vote)
}

// WithdrawVote removes a user's vote, if any.
func (m *Model) WithdrawVote(userID string) {
	m.votes = slices.DeleteFunc(m.votes, func(v Vote) bool {
		return v.userID == userID
	})
}

// Score is the net feedback on the challenge: thumbs up minus thumbs down.
func (m *Model) Score() int {
	score := 0
	for _, v := range m.votes {
		switch v.sentiment {
		case Up:
			score++
		case Down:
			score--
		}
	}
	return score
}
