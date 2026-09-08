package memorystore_test

import (
	"context"
	"testing"
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/popularity"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/memorystore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterCompletion(t *testing.T) {
	challengeID := "123"
	myChallenge := challenge.Model{}
	completion := challenge.NewCompletion("user", time.Second)

	store := memorystore.New()
	err := store.PutChallenge(context.Background(), challengeID, myChallenge)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	err = store.RegisterCompletion(context.Background(), challengeID, completion)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	gotChallenge, err := store.GetChallenge(context.Background(), challengeID)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	if len(gotChallenge.Completions()) != 1 {
		t.Errorf("completion not applied")
	}

}

func TestReadsAreIsolatedFromStoredState(t *testing.T) {
	ctx := context.Background()
	store := memorystore.New()

	stored := challenge.NewChallenge(stage.Model{}, weather.DRY, car.Model{}, []challenge.Completion{
		challenge.NewCompletion("alice", time.Minute),
	}, nil)
	require.NoError(t, store.PutChallenge(ctx, "c1", stored))

	// Mutating a value read back must not reach the store.
	got, err := store.GetChallenge(ctx, "c1")
	require.NoError(t, err)
	got.RegisterCompletion(challenge.NewCompletion("mallory", time.Hour))

	again, err := store.GetChallenge(ctx, "c1")
	require.NoError(t, err)
	require.Len(t, again.Completions(), 1)
}

func TestCancelledContextIsRefused(t *testing.T) {
	store := memorystore.New()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	require.ErrorIs(t, store.PutChallenge(ctx, "123", challenge.Model{}), context.Canceled)

	_, err := store.GetChallenge(ctx, "123")
	require.ErrorIs(t, err, context.Canceled)

	require.ErrorIs(t, store.DeleteChallenge(ctx, "123"), context.Canceled)
	require.ErrorIs(t, store.RegisterCompletion(ctx, "123", challenge.Completion{}), context.Canceled)
	require.ErrorIs(t, store.RegisterVote(ctx, "123", challenge.Vote{}), context.Canceled)

	_, err = store.Popularity(ctx)
	require.ErrorIs(t, err, context.Canceled)
}

func TestRegisterVote(t *testing.T) {
	ctx := context.Background()
	store := memorystore.New()

	c := challenge.NewChallenge(
		stage.New("Sweet Lamb", location.WAL, stage.Long),
		weather.WET,
		car.New("Lancia Delta S4", class.GroupB4WD),
		nil, nil,
	)
	require.NoError(t, store.PutChallenge(ctx, "c1", c))

	locationKey := popularity.LocationKey(location.WAL)

	// A first vote counts towards every constituent item of the challenge.
	require.NoError(t, store.RegisterVote(ctx, "c1", challenge.NewVote("alice", challenge.Up)))
	snapshot, err := store.Popularity(ctx)
	require.NoError(t, err)
	assert.Len(t, snapshot, len(popularity.KeysFor(c)))
	assert.Equal(t, popularity.Tally{Up: 1}, snapshot[locationKey])

	// Voting the same way again withdraws it, leaving no trace.
	require.NoError(t, store.RegisterVote(ctx, "c1", challenge.NewVote("alice", challenge.Up)))
	snapshot, err = store.Popularity(ctx)
	require.NoError(t, err)
	assert.Empty(t, snapshot)

	// Voting the other way flips an existing vote.
	require.NoError(t, store.RegisterVote(ctx, "c1", challenge.NewVote("alice", challenge.Up)))
	require.NoError(t, store.RegisterVote(ctx, "c1", challenge.NewVote("alice", challenge.Down)))
	snapshot, err = store.Popularity(ctx)
	require.NoError(t, err)
	assert.Equal(t, popularity.Tally{Down: 1}, snapshot[locationKey])

	// The vote is recorded on the challenge too.
	got, err := store.GetChallenge(ctx, "c1")
	require.NoError(t, err)
	assert.Equal(t, -1, got.Score())
}
