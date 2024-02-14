package challenge

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/feedback/event"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
)

type Randomiser interface {
	Car() car.Model
	Loc() location.Model
	Weather(loc location.Model) weather.Model
	Stage(loc location.Model) stage.Model
}

type Model struct {
	Stage       stage.Model
	Weather     weather.Model
	Car         car.Model
	events      []any        // TODO EVENTTYPE
	completions []completion // implicitly ordered, maybe change to explicit?
}

type completion struct {
	userID          string
	userDisplayName string
	duration        time.Duration
}

func New(r Randomiser) *Model {
	loc := r.Loc()

	return &Model{
		Stage:   r.Stage(loc),
		Weather: r.Weather(loc),
		Car:     r.Car(),
	}
}

func (m *Model) Events() []any {
	return m.events
}

func (m *Model) Completions() []completion {
	return m.completions
}

func (m *Model) FancyString() string {
	return strings.Join([]string{
		m.Stage.FancyString(),
		location.WeatherStrings()[m.Stage.Location()][m.Weather],
		m.Car.FancyString(),
		""},
		"\n")
}

func (m *Model) ApplyEvent(e any) error { // TODO EVENTTYPE
	switch t := e.(type) {
	case event.Completion:
		m.events = append(m.events, e)
		m.completions = append(m.completions, completion{userID: e.(event.Completion).UserID(), userDisplayName: e.(event.Completion).UserDisplayName(), duration: e.(event.Completion).Duration()})
		return nil
	default:
		return fmt.Errorf("type %s not a valid event", t)
	}
}

func (m *Model) TopThreeFancyString() string {
	if len(m.completions) < 1 {
		return ""
	}

	if len(m.completions) == 1 {
		return fmt.Sprintf("🥇 **%s** %s", formatDuration(m.completions[0].duration), m.completions[0].userDisplayName)
	}

	sorted := make([]completion, len(m.completions))
	//lint:ignore S1001 copy doesn't work on unexported struct fields
	for i := range m.completions { //nolint:gosimple // copy doesn't work on unexported struct fields
		sorted[i] = m.completions[i]
	}

	// VSCode running go-staticcheck & golangci-lint fighting it out? sheesh
	// TODO - sort out this linter nonsense, probably refactor the above

	slices.SortFunc(sorted, func(a, b completion) int { return int(a.duration - b.duration) })

	fmt.Println(sorted)

	if len(sorted) == 2 {
		return strings.Join([]string{
			fmt.Sprintf("🥇 *%s* %s", formatDuration(sorted[0].duration), sorted[0].userDisplayName),
			fmt.Sprintf("🥈 *%s* %s", formatDuration(sorted[1].duration), sorted[1].userDisplayName),
		},
			"\n")
	}
	return strings.Join([]string{
		fmt.Sprintf("🥇 *%s* %s", formatDuration(sorted[0].duration), sorted[0].userDisplayName),
		fmt.Sprintf("🥈 *%s* %s", formatDuration(sorted[1].duration), sorted[1].userDisplayName),
		fmt.Sprintf("🥉 *%s* %s", formatDuration(sorted[2].duration), sorted[2].userDisplayName),
	},
		"\n")
}

func formatDuration(d time.Duration) string {
	// TODO is there a better way?
	minutes := d.Truncate(time.Minute)
	seconds := (d - minutes).Truncate(time.Second)
	milliseconds := (d - d.Truncate(seconds)).Truncate(time.Millisecond)
	return fmt.Sprintf("%2.f:%2.f.%d", minutes.Minutes(), seconds.Seconds(), milliseconds.Milliseconds())
}
