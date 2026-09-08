package web

import (
	"sort"
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/timestamp"
	"github.com/bwmarrin/discordgo"
)

const (
	dateFormat = "2006-01-02 15:04"
	noValue    = "—"

	sortByTime      = "time"
	sortBySubmitted = "submitted"
	dirAscending    = "asc"
	dirDescending   = "desc"
)

// challengeView is one stored challenge in the listing, together with its
// completions panel.
type challengeView struct {
	ID          string
	created     time.Time
	Date        string
	Track       string
	Distance    string
	Weather     string
	Car         string
	Class       string
	Completions completionsView
}

// completionsView is the sortable completions panel for one challenge. The same
// value backs both the initial page render and the htmx re-sort fragment.
type completionsView struct {
	ChallengeID string
	Sort        string
	Dir         string
	Rows        []completionRow
}

type completionRow struct {
	Name      string
	Time      string
	Submitted string
}

// challengeViews builds the listing, newest challenge first. A challenge ID is a
// Discord snowflake, so its date comes from the ID itself; an unparseable ID
// sorts last with no date.
func challengeViews(challenges map[string]challenge.Model) []challengeView {
	views := make([]challengeView, 0, len(challenges))

	for id, c := range challenges {
		s := c.Stage()

		view := challengeView{
			ID:          id,
			Date:        noValue,
			Track:       s.Location().String() + " » " + s.Name(),
			Distance:    s.Distance().String(),
			Weather:     c.Weather().String(),
			Car:         c.Car().Name(),
			Class:       c.Car().Class().String(),
			Completions: newCompletionsView(id, c.Completions(), sortByTime, dirAscending),
		}

		if created, err := discordgo.SnowflakeTimestamp(id); err == nil {
			view.created = created
			view.Date = created.Format(dateFormat)
		}

		views = append(views, view)
	}

	sort.Slice(views, func(i, j int) bool {
		return views[i].created.After(views[j].created)
	})

	return views
}

// newCompletionsView orders the completions by the requested key and direction,
// falling back to time ascending for anything unrecognised.
func newCompletionsView(challengeID string, completions []challenge.Completion, sortKey, dir string) completionsView {
	if sortKey != sortBySubmitted {
		sortKey = sortByTime
	}
	if dir != dirDescending {
		dir = dirAscending
	}

	ordered := make([]challenge.Completion, len(completions))
	copy(ordered, completions)

	less := completionLess(sortKey)
	sort.SliceStable(ordered, func(i, j int) bool {
		if dir == dirDescending {
			return less(ordered[j], ordered[i])
		}
		return less(ordered[i], ordered[j])
	})

	rows := make([]completionRow, len(ordered))
	for i, completion := range ordered {
		rows[i] = completionRow{
			Name:      completionName(completion),
			Time:      timestamp.Format(completion.Duration()),
			Submitted: formatSubmitted(completion.SubmittedAt()),
		}
	}

	return completionsView{ChallengeID: challengeID, Sort: sortKey, Dir: dir, Rows: rows}
}

func completionLess(sortKey string) func(a, b challenge.Completion) bool {
	if sortKey == sortBySubmitted {
		return func(a, b challenge.Completion) bool {
			return a.SubmittedAt().Before(b.SubmittedAt())
		}
	}
	return func(a, b challenge.Completion) bool {
		return a.Duration() < b.Duration()
	}
}

func completionsLabel(count int) string {
	if count == 1 {
		return "completion"
	}
	return "completions"
}

func completionName(c challenge.Completion) string {
	if c.DisplayName() != "" {
		return c.DisplayName()
	}
	return c.UserID()
}

func formatSubmitted(t time.Time) string {
	if t.IsZero() {
		return noValue
	}
	return t.Format(dateFormat)
}

// Count is the number of completions in the panel.
func (v completionsView) Count() int { return len(v.Rows) }

// active reports whether the panel is currently ordered by column.
func (v completionsView) active(column string) bool { return v.Sort == column }

// nextDir is the direction a header link for column should request: the opposite
// of the current one when the column is already active, ascending otherwise.
func (v completionsView) nextDir(column string) string {
	if v.Sort == column && v.Dir == dirAscending {
		return dirDescending
	}
	return dirAscending
}

// arrow is the sort indicator shown next to an active column's label.
func (v completionsView) arrow(column string) string {
	switch {
	case v.Sort != column:
		return ""
	case v.Dir == dirDescending:
		return " ▾"
	default:
		return " ▴"
	}
}
