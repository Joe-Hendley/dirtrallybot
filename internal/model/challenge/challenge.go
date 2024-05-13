package challenge

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/event"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
	"github.com/bwmarrin/discordgo"
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
	completions []completion // ordered by submission time, but timestamp is on the event so we could add it to the completion
}

type completion struct {
	userID   string
	duration time.Duration
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
		m.completions = append(m.completions, completion{userID: e.(event.Completion).UserID(), duration: e.(event.Completion).Duration()})
		return nil
	default:
		return fmt.Errorf("type %s not a valid event", t)
	}
}

func (m *Model) TopThreeFancyString(s *discordgo.Session, guildID string) string {
	if len(m.completions) < 1 {
		return ""
	}

	if len(m.completions) == 1 {
		return fmt.Sprintf("🥇 **%s**\t%s", formatDuration(m.completions[0].duration), getCurrentDisplayName(s, guildID, m.completions[0].userID))
	}

	sorted := make([]completion, len(m.completions))
	//lint:ignore S1001 copy doesn't work on unexported struct fields
	for i := range m.completions { //nolint:gosimple // copy doesn't work on unexported struct fields
		sorted[i] = m.completions[i]
	}

	// go-staticcheck & golangci-lint fighting it out? sheesh
	// TODO - sort out this linter nonsense, probably refactor the above
	// could sort a list of indices instead?

	slices.SortFunc(sorted, func(a, b completion) int { return int(a.duration - b.duration) })

	if len(sorted) == 2 {
		return strings.Join([]string{
			fmt.Sprintf("🥇 **%s**\t%s", formatDuration(sorted[0].duration), getCurrentDisplayName(s, guildID, m.completions[0].userID)),
			fmt.Sprintf("🥈 **%s**\t%s", formatDuration(sorted[1].duration), getCurrentDisplayName(s, guildID, m.completions[1].userID)),
		},
			"\n")
	}
	return strings.Join([]string{
		fmt.Sprintf("🥇\t**%-s**\t\t%s", formatDuration(sorted[0].duration), getCurrentDisplayName(s, guildID, m.completions[0].userID)),
		fmt.Sprintf("🥈\t**%-s**\t\t%s", formatDuration(sorted[1].duration), getCurrentDisplayName(s, guildID, m.completions[1].userID)),
		fmt.Sprintf("🥉\t**%-s**\t\t%s", formatDuration(sorted[2].duration), getCurrentDisplayName(s, guildID, m.completions[2].userID)),
	},
		"\n")
}

func formatDuration(d time.Duration) string {
	var (
		minutes      = d.Truncate(time.Minute)
		seconds      = (d - minutes).Truncate(time.Second)
		milliseconds = (d - minutes - seconds).Truncate(time.Millisecond)
	)

	minuteComponent := fmt.Sprintf("%2.f", minutes.Minutes())
	secondComponent := fmt.Sprintf("%02.f", seconds.Seconds())
	millisecondComponent := fmt.Sprintf("%d", milliseconds.Milliseconds())

	return minuteComponent + ":" + secondComponent + "." + millisecondComponent
}

func (m *Model) FancyListCompletions(s *discordgo.Session, guildID string) string {
	if len(m.completions) == 0 {
		return "no completions logged"
	}

	userCompletions := make(map[string][]time.Duration)

	for _, completion := range m.completions {
		userCompletions[completion.userID] = append(userCompletions[completion.userID], completion.duration)
	}

	type user struct {
		id          string
		displayName string
	}

	users := []user{}
	for userID := range userCompletions {
		users = append(users, user{id: userID, displayName: getCurrentDisplayName(s, guildID, userID)})
	}

	slices.SortFunc(users, func(a user, b user) int {
		if a.displayName < b.displayName {
			return -1
		}
		if a.displayName > b.displayName {
			return 1
		}
		return 0
	})

	buf := strings.Builder{}

	for _, user := range users {
		buf.Write([]byte("**" + user.displayName + "**\n"))
		for _, completion := range userCompletions[user.id] {
			buf.Write([]byte(formatDuration(completion) + "\n"))
		}
	}

	return buf.String()
}

func getCurrentDisplayName(s *discordgo.Session, guildID, userID string) string {
	u, err := s.GuildMember(guildID, userID)
	if err != nil {
		slog.Error("getting display name", "guildID", guildID, "userID", userID, "err", err)
		return userID
	}

	return u.DisplayName()
}
