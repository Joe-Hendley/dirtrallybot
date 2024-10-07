package boltstore_test

import (
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/randomiser"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/boltstore"
)

func TestBoltStore(t *testing.T) {
	store := MustCreateStore(t)
	myChallenge := challenge.Model{}
	challengeID := "12345"

	err := store.PutChallenge(challengeID, myChallenge)
	if err != nil {
		t.Errorf("error putting challenge: %s", err)
	}

	_, err = store.GetChallenge(challengeID)
	if err != nil {
		t.Errorf("error getting challenge: %s", err)
	}

	err = store.RegisterCompletion(challengeID, challenge.Completion{})
	if err != nil {
		t.Errorf("error adding completion: %s", err)
	}

	gotChallenge, err := store.GetChallenge(challengeID)
	if err != nil {
		t.Errorf("error getting challenge: %s", err)
	}

	completions := gotChallenge.Completions()
	if len(completions) != 1 {
		t.Errorf("unexpected # of completions")
	}
}

func TestPutAndGet(t *testing.T) {
	r := randomiser.NewSimple()
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
		err := store.PutChallenge(challengeID, challenge)
		if err != nil {
			t.Errorf("error putting challenge: %s", err)
		}
	}

	for _, challengeID := range challengeIDs {
		got, err := store.GetChallenge(challengeID)
		if err != nil {
			t.Errorf("error getting challenge: %s", err)
		}

		want := challenges[challengeID]
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	}
}

func TestRegisterCompletion(t *testing.T) {
	store := MustCreateStore(t)
	challengeID := "123"
	myChallenge := challenge.NewRandomChallenge(challenge.Config{}, randomiser.NewSimple())

	err := store.PutChallenge(challengeID, myChallenge)
	if err != nil {
		t.Errorf("error putting challenge: %s", err)
	}

	wantCompletions := []challenge.Completion{
		challenge.NewCompletion("someUser", time.Minute),
		challenge.NewCompletion("someUser", time.Hour),
		challenge.NewCompletion("someOtherUser", time.Hour),
	}

	for _, completion := range wantCompletions {
		err = store.RegisterCompletion(challengeID, completion)
		if err != nil {
			t.Errorf("error registering completion: %s", err)
		}
	}

	gotChallenge, err := store.GetChallenge(challengeID)
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
