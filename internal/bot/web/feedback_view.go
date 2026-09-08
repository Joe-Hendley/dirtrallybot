package web

import (
	"sort"
	"strconv"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/popularity"
)

// kindOrder is the order the feedback tables appear in, mirroring how a
// challenge reads: where it is, then what you drive.
var kindOrder = []popularity.Kind{
	popularity.KindLocation,
	popularity.KindStage,
	popularity.KindDistance,
	popularity.KindWeather,
	popularity.KindClass,
	popularity.KindCar,
	popularity.KindDrivetrain,
}

// feedbackTable is one Kind's tally, ready to render.
type feedbackTable struct {
	Kind string
	Rows []feedbackRow
}

type feedbackRow struct {
	Name string
	Up   int
	Down int
	Net  int
}

// feedbackTables groups the snapshot by Kind, ordering each table by net score
// descending then name. A Kind nobody has voted on is left out.
func feedbackTables(snapshot popularity.Snapshot) []feedbackTable {
	rowsByKind := make(map[popularity.Kind][]feedbackRow)
	for key, tally := range snapshot {
		rowsByKind[key.Kind] = append(rowsByKind[key.Kind], feedbackRow{
			Name: key.Label(),
			Up:   tally.Up,
			Down: tally.Down,
			Net:  tally.Net(),
		})
	}

	tables := make([]feedbackTable, 0, len(rowsByKind))
	for _, kind := range kindOrder {
		rows, ok := rowsByKind[kind]
		if !ok {
			continue
		}

		sort.Slice(rows, func(i, j int) bool {
			if rows[i].Net != rows[j].Net {
				return rows[i].Net > rows[j].Net
			}
			return rows[i].Name < rows[j].Name
		})

		tables = append(tables, feedbackTable{Kind: kind.String(), Rows: rows})
	}

	return tables
}

// feedbackSummary is the count line shown under the heading.
func feedbackSummary(tables []feedbackTable) string {
	items := 0
	for _, t := range tables {
		items += len(t.Rows)
	}

	noun := "items"
	if items == 1 {
		noun = "item"
	}
	return strconv.Itoa(items) + " " + noun + " with feedback"
}

// signed formats a net score with an explicit sign, so a positive tally is not
// mistaken for a neutral one.
func signed(n int) string {
	if n > 0 {
		return "+" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}
