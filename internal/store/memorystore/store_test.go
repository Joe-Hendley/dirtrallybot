package memorystore_test

import (
	"context"
	"testing"
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/memorystore"
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

func TestCancelledContextIsRefused(t *testing.T) {
	store := memorystore.New()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	require.ErrorIs(t, store.PutChallenge(ctx, "123", challenge.Model{}), context.Canceled)

	_, err := store.GetChallenge(ctx, "123")
	require.ErrorIs(t, err, context.Canceled)

	require.ErrorIs(t, store.DeleteChallenge(ctx, "123"), context.Canceled)
	require.ErrorIs(t, store.RegisterCompletion(ctx, "123", challenge.Completion{}), context.Canceled)
}
