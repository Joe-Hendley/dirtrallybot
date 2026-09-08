package boltstore_test

import (
	"context"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/drivetrain"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/popularity"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
	"github.com/Joe-Hendley/dirtrallybot/internal/randomiser"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/boltstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoltStore(t *testing.T) {
	store := MustCreateStore(t)
	myChallenge := challenge.Model{}
	challengeID := "12345"

	err := store.PutChallenge(context.Background(), challengeID, myChallenge)
	if err != nil {
		t.Errorf("error putting challenge: %s", err)
	}

	_, err = store.GetChallenge(context.Background(), challengeID)
	if err != nil {
		t.Errorf("error getting challenge: %s", err)
	}

	err = store.RegisterCompletion(context.Background(), challengeID, challenge.Completion{})
	if err != nil {
		t.Errorf("error adding completion: %s", err)
	}

	gotChallenge, err := store.GetChallenge(context.Background(), challengeID)
	if err != nil {
		t.Errorf("error getting challenge: %s", err)
	}

	completions := gotChallenge.Completions()
	if len(completions) != 1 {
		t.Errorf("unexpected # of completions")
	}
}

func TestPutAndGet(t *testing.T) {
	r := randomiser.NewDeterministic(game.DR2)
	challengeIDs := []string{
		"challenge1",
		"challenge2",
		"challenge3",
	}

	challenges := map[string]challenge.Model{
		challengeIDs[0]: challenge.NewRandomChallenge(challenge.Config{}, r),
		challengeIDs[1]: challenge.NewRandomChallenge(challenge.Config{}, r),
		challengeIDs[2]: challenge.NewRandomChallenge(challenge.Config{}, r),
	}

	store := MustCreateStore(t)

	for challengeID, challenge := range challenges {
		err := store.PutChallenge(context.Background(), challengeID, challenge)
		if err != nil {
			t.Errorf("error putting challenge: %s", err)
		}
	}

	for _, challengeID := range challengeIDs {
		got, err := store.GetChallenge(context.Background(), challengeID)
		if err != nil {
			t.Errorf("error getting challenge: %s", err)
		}

		want := challenges[challengeID]
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	}
}

func TestListChallenges(t *testing.T) {
	store := MustCreateStore(t)
	ctx := context.Background()

	empty, err := store.ListChallenges(ctx)
	require.NoError(t, err)
	assert.Empty(t, empty)

	r := randomiser.NewDeterministic(game.DR2)
	want := map[string]challenge.Model{
		"c1": challenge.NewRandomChallenge(challenge.Config{}, r),
		"c2": challenge.NewRandomChallenge(challenge.Config{}, r),
		"c3": challenge.NewRandomChallenge(challenge.Config{}, r),
	}

	for id, c := range want {
		require.NoError(t, store.PutChallenge(ctx, id, c))
	}

	got, err := store.ListChallenges(ctx)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestRegisterCompletionPersistsNameAndSubmissionTime(t *testing.T) {
	store := MustCreateStore(t)
	ctx := context.Background()

	require.NoError(t, store.PutChallenge(ctx, "c1", challenge.Model{}))

	want := challenge.NewCompletionAt("alice-id", "Alice", 90*time.Second,
		time.Date(2024, 3, 4, 5, 6, 7, 0, time.UTC))
	require.NoError(t, store.RegisterCompletion(ctx, "c1", want))

	got, err := store.GetChallenge(ctx, "c1")
	require.NoError(t, err)
	require.Len(t, got.Completions(), 1)
	assert.Equal(t, want, got.Completions()[0])
}

func TestRegisterCompletion(t *testing.T) {
	store := MustCreateStore(t)
	challengeID := "123"
	myChallenge := challenge.NewRandomChallenge(challenge.Config{}, randomiser.NewDeterministic(game.DR2))

	err := store.PutChallenge(context.Background(), challengeID, myChallenge)
	if err != nil {
		t.Errorf("error putting challenge: %s", err)
	}

	wantCompletions := []challenge.Completion{
		challenge.NewCompletion("someUser", time.Minute),
		challenge.NewCompletion("someUser", time.Hour),
		challenge.NewCompletion("someOtherUser", time.Hour),
	}

	for _, completion := range wantCompletions {
		err = store.RegisterCompletion(context.Background(), challengeID, completion)
		if err != nil {
			t.Errorf("error registering completion: %s", err)
		}
	}

	gotChallenge, err := store.GetChallenge(context.Background(), challengeID)
	if err != nil {
		t.Errorf("error getting challenge: %s", err)
	}

	if len(gotChallenge.Completions()) != len(wantCompletions) {
		t.Errorf("got %d completions, want %d completions", len(gotChallenge.Completions()), len(wantCompletions))
	}

	for idx, completion := range gotChallenge.Completions() {
		if !reflect.DeepEqual(completion, wantCompletions[idx]) {
			t.Errorf("got %v expect %v", completion, wantCompletions[idx])
		}
	}
}

func TestRegisterVotePersistsPopularity(t *testing.T) {
	store := MustCreateStore(t)
	ctx := context.Background()

	c := challenge.NewChallenge(
		stage.New("Sweet Lamb", location.WAL, stage.Long),
		weather.WET,
		car.New("Lancia Delta S4", class.GroupB4WD),
		nil, nil,
	)
	require.NoError(t, store.PutChallenge(ctx, "c1", c))
	require.NoError(t, store.RegisterVote(ctx, "c1", challenge.NewVote("alice", challenge.Up)))
	require.NoError(t, store.RegisterVote(ctx, "c1", challenge.NewVote("bob", challenge.Down)))

	snapshot, err := store.Popularity(ctx)
	require.NoError(t, err)

	// Keys survive the encode/decode round-trip, including the name-qualified ones.
	assert.Equal(t, popularity.Tally{Up: 1, Down: 1},
		snapshot[popularity.StageKey(stage.New("Sweet Lamb", location.WAL, stage.Long))])
	assert.Equal(t, popularity.Tally{Up: 1, Down: 1},
		snapshot[popularity.CarKey(car.New("Lancia Delta S4", class.GroupB4WD))])
	assert.Equal(t, popularity.Tally{Up: 1, Down: 1},
		snapshot[popularity.DrivetrainKey(drivetrain.AWD)])

	// bob withdraws; his down is reversed everywhere.
	require.NoError(t, store.RegisterVote(ctx, "c1", challenge.NewVote("bob", challenge.Down)))
	snapshot, err = store.Popularity(ctx)
	require.NoError(t, err)
	assert.Equal(t, popularity.Tally{Up: 1},
		snapshot[popularity.LocationKey(location.WAL)])
}

func MustCreateStore(t *testing.T) *boltstore.Store {
	path := filepath.Join(t.TempDir(), "db")
	t.Logf("opening store at %s", path)

	store, err := boltstore.New(path)
	if err != nil {
		t.Errorf("error starting bolt store at %s: %s", path, err)
	}

	t.Cleanup(func() {
		err := store.Close()
		if err != nil {
			t.Errorf("error closing bolt store at %s: %s", path, err)
		}
	})

	return store
}
