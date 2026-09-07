package challenge

import (
	"reflect"
	"slices"
	"testing"
	"time"

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
		assert.Empty(t, challenge.TopThree())
	})

	t.Run("fewer than two completions -> returned as-is", func(t *testing.T) {
		only := NewCompletion("alice", time.Minute)
		challenge := newChallengeWithCompletions([]Completion{only})
		assert.Equal(t, []Completion{only}, challenge.TopThree())
	})

	t.Run("fastest per user, ascending, input slice not mutated", func(t *testing.T) {
		aliceSlow := NewCompletion("alice", 5*time.Minute)
		bob := NewCompletion("bob", 4*time.Minute)
		carol := NewCompletion("carol", 6*time.Minute)
		aliceFast := NewCompletion("alice", 3*time.Minute)
		dave := NewCompletion("dave", 7*time.Minute)

		completions := []Completion{aliceSlow, bob, carol, aliceFast, dave}
		original := slices.Clone(completions)

		challenge := newChallengeWithCompletions(completions)
		got := challenge.TopThree()

		// Three fastest distinct users, ascending by duration.
		assert.Equal(t, []Completion{aliceFast, bob, carol}, got)

		// Unexported fields survive the internal copy.
		assert.Equal(t, "alice", got[0].UserID())
		assert.Equal(t, 3*time.Minute, got[0].Duration())

		// TopThree copies before sorting, so the caller's slice is untouched.
		assert.Equal(t, original, completions)
	})
}
