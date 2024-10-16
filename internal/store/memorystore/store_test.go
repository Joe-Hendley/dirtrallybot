package memorystore_test

import (
	"testing"
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/memorystore"
)

func TestRegisterCompletion(t *testing.T) {
	challengeID := "123"
	myChallenge := challenge.Model{}
	completion := challenge.NewCompletion("user", time.Second)

	store := memorystore.New()
	err := store.PutChallenge(challengeID, myChallenge)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	err = store.RegisterCompletion(challengeID, completion)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	gotChallenge, err := store.GetChallenge(challengeID)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	if len(gotChallenge.Completions()) != 1 {
		t.Errorf("completion not applied")
	}

}
