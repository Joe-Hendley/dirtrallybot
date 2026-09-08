package challenge

import (
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/drivetrain"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeRandomiser records the arguments it is called with and returns fixed
// values, so tests can assert which path NewRandomChallenge took.
type fakeRandomiser struct {
	stage            stage.Model
	stageOfDistance  stage.Model
	distanceAskedFor *stage.Distance
	stageCalled      bool
}

func (f *fakeRandomiser) Car() car.Model                               { return car.Model{} }
func (f *fakeRandomiser) CarFromClass(class.Model) car.Model           { return car.Model{} }
func (f *fakeRandomiser) CarFromDrivetrain(drivetrain.Model) car.Model { return car.Model{} }
func (f *fakeRandomiser) Loc() location.Model                          { return location.SCO }
func (f *fakeRandomiser) Weather(location.Model) weather.Model         { return weather.DRY }

func (f *fakeRandomiser) Stage(location.Model) stage.Model {
	f.stageCalled = true
	return f.stage
}

func (f *fakeRandomiser) StageOfDistance(_ location.Model, d stage.Distance) stage.Model {
	f.distanceAskedFor = &d
	return f.stageOfDistance
}

func TestNewRandomChallengeUsesDistance(t *testing.T) {
	longStage := stage.New("Sweet Lamb", location.WAL, stage.Long)
	fake := &fakeRandomiser{
		stage:           stage.New("Fferm Wynt", location.WAL, stage.Short),
		stageOfDistance: longStage,
	}

	loc := location.WAL
	distance := stage.Long
	got := NewRandomChallenge(Config{Location: &loc, Distance: &distance}, fake)

	require.NotNil(t, fake.distanceAskedFor)
	assert.Equal(t, stage.Long, *fake.distanceAskedFor)
	assert.False(t, fake.stageCalled, "should not fall back to Stage when a distance is set")
	assert.Equal(t, longStage, got.Stage())
}

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

func TestVotes(t *testing.T) {
	newChallenge := func() Model {
		return NewChallenge(stage.Model{}, weather.DRY, car.Model{}, nil, nil)
	}

	t.Run("one vote per user, latest wins", func(t *testing.T) {
		challenge := newChallenge()
		challenge.SetVote(NewVote("alice", Up))
		challenge.SetVote(NewVote("bob", Down))
		challenge.SetVote(NewVote("alice", Down))

		assert.Len(t, challenge.Votes(), 2)
		got, ok := challenge.VoteFor("alice")
		require.True(t, ok)
		assert.Equal(t, Down, got.Sentiment())
		assert.Equal(t, -2, challenge.Score())
	})

	t.Run("withdraw removes the user's vote", func(t *testing.T) {
		challenge := newChallenge()
		challenge.SetVote(NewVote("alice", Up))
		challenge.SetVote(NewVote("bob", Up))
		challenge.WithdrawVote("alice")

		_, ok := challenge.VoteFor("alice")
		assert.False(t, ok)
		assert.Equal(t, 1, challenge.Score())
	})
}

func TestTopThree(t *testing.T) {
	newChallengeWithCompletions := func(completions []Completion) Model {
		return NewChallenge(stage.Model{}, weather.DRY, car.Model{}, completions, nil)
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
