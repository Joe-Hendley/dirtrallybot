package challenge

import (
	"reflect"
	"slices"
	"testing"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
	"github.com/stretchr/testify/assert"
)

func TestSortUser(t *testing.T) {
	users := []string{"Bob", "Alice", "Carol"}
	slices.SortFunc(users, func(a, b string) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	})
	want := []string{"Alice", "Bob", "Carol"}

	if !reflect.DeepEqual(users, want) {
		t.Errorf("got %v want %v", users, want)
	}
}

func TestTopThree(t *testing.T) {
	newChallengeWithCompletions := func(completions []Completion) Model {
		return NewChallenge(stage.Model{}, weather.DRY, car.Model{}, completions)
	}
	t.Run("no completions -> empty array", func(t *testing.T) {
		completions := []Completion{}
		challenge := newChallengeWithCompletions(completions)
		assert.Len(t, challenge.TopThree(), 0)
	})
}
