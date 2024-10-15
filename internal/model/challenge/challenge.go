package challenge

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/drivetrain"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/timestamp"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
)

var EmojiDelimiter = string([]byte{0x1f})

var (
	RandomEmoji       = "🎲"
	RandomString      = "Random"
	RandomFancyString = RandomEmoji + " " + EmojiDelimiter + RandomString
)

type Randomiser interface {
	Car() car.Model
	CarFromClass(class class.Model) car.Model
	CarFromDrivetrain(drivetrain drivetrain.Model) car.Model

	Loc() location.Model
	Weather(loc location.Model) weather.Model
	Stage(loc location.Model) stage.Model
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

	if c.Stage != nil {
		challenge.stage = *c.Stage
	} else {
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

func (m *Model) FancyString() string {
	return strings.Join([]string{
		m.stage.FancyString(),
		location.WeatherStrings()[m.Stage().Location()][m.Weather()],
		m.car.FancyString(),
		""},
		"\n")
}

func (m *Model) RegisterCompletion(c Completion) {
	m.completions = append(m.completions, c)
}

func (m *Model) TopThree() []Completion {
	if len(m.completions) < 2 {
		return m.completions
	}

	sorted := make([]Completion, len(m.completions))
	//lint:ignore S1001 copy doesn't work on unexported struct fields
	for i := range m.completions { //nolint:gosimple // copy doesn't work on unexported struct fields
		sorted[i] = m.completions[i]
	}

	// go-staticcheck & golangci-lint fighting it out? sheesh
	// TODO - sort out this linter nonsense, probably refactor the above
	// could sort a list of indices instead?

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

func (m *Model) FancyListCompletions() map[string]string {
	if len(m.completions) == 0 {
		return map[string]string{}
	}

	userCompletions := make(map[string][]time.Duration)

	for _, completion := range m.completions {
		userCompletions[completion.userID] = append(userCompletions[completion.userID], completion.duration)
	}

	userCompletionStrings := map[string]string{}
	for userID, completions := range userCompletions {
		buf := strings.Builder{}
		for _, completion := range completions {
			buf.Write([]byte(timestamp.Format(completion) + "\n"))
		}
		userCompletionStrings[userID] = buf.String()
	}

	return userCompletionStrings
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

func (c Config) String() string {
	stringParts := []string{}
	if c.Stage != nil {
		stringParts = append(stringParts, "stage: "+c.Stage.String())
	} else if c.Location != nil {
		stringParts = append(stringParts, "loc: "+c.Location.String())
	}

	if c.Weather != nil {
		stringParts = append(stringParts, "weather: "+c.Weather.String())
	}

	if c.Car != nil {
		stringParts = append(stringParts, "car: "+c.Car.String())
	} else if c.Class != nil {
		stringParts = append(stringParts, "class: "+c.Class.String())
	} else if c.Drivetrain != nil {
		stringParts = append(stringParts, "drivetrain: "+c.Drivetrain.FancyString())
	}

	return strings.Join(stringParts, ", ")
}

func (c Config) FancyStageString() string {

	var (
		locationString = RandomFancyString
		stageString    = RandomFancyString
		weatherString  = RandomFancyString
	)

	if c.Stage != nil {
		stageString = c.Stage.Distance().Emoji() + " " + EmojiDelimiter + c.Stage.Name()
		locationString = c.Location.Flag() + " " + EmojiDelimiter + c.Location.String()
	} else if c.Location != nil {
		locationString = c.Location.Flag() + " " + EmojiDelimiter + c.Location.String()
	}

	if c.Stage == nil && c.Distance != nil {
		stageString = c.Distance.Emoji() + " " + EmojiDelimiter + RandomString + " " + c.Distance.String()
	}

	locationHasOneWeatherType := c.Location != nil && len(c.Location.Weather()) == 1

	switch {
	case c.Weather != nil:
		weatherString = fmt.Sprintf("%s %s%s", c.Weather.Emoji(), EmojiDelimiter, c.Weather.String())
	case c.Weather == nil && !locationHasOneWeatherType:
		weatherString = RandomFancyString

	case c.Weather == nil && locationHasOneWeatherType:
		weatherString = fmt.Sprintf("%s *(probably %s though)*", RandomFancyString, c.Location.Weather()[0].String())
	}
	return fmt.Sprintf("Location: %s\nStage: %s\nWeather: %s", locationString, stageString, weatherString)
}

func (c Config) FancyCarString() string {
	var (
		drivetrainString = RandomFancyString
		classString      = RandomFancyString
		carString        = RandomFancyString
	)

	switch {
	case c.Car != nil:
		drivetrainString = c.Car.Class().Drivetrain().Emoji() + " " + EmojiDelimiter + c.Car.Class().Drivetrain().String()
		classString = EmojiDelimiter + c.Car.Class().String()
		carString = EmojiDelimiter + c.Car.Name()
	case c.Class != nil:
		drivetrainString = c.Class.Drivetrain().Emoji() + " " + EmojiDelimiter + c.Class.Drivetrain().String()
		classString = c.Class.String()
	case c.Drivetrain != nil:
		drivetrainString = c.Drivetrain.Emoji() + " " + EmojiDelimiter + c.Drivetrain.String()
	}

	return fmt.Sprintf("Drivetrain: %s\nClass: %s\nCar: %s", drivetrainString, classString, carString)
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
