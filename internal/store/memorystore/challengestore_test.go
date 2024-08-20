package memorystore

import (
	"sync"
	"testing"
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/model"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/event"
)

func TestApplyEvent(t *testing.T) {
	challengeID := "123"
	newStoreWithEvent := func(t *testing.T, e model.Event) Store {
		t.Helper()
		originalChallenge := &challenge.Model{}
		store := Store{
			lock:         &sync.Mutex{},
			challengeMap: map[string]*challenge.Model{},
		}

		store.Put(challengeID, originalChallenge)
		err := store.ApplyEvent(challengeID, e)
		if err != nil {
			t.Errorf("couldn't apply event: %v", err)
		}

		return store
	}

	t.Run("event is added to slice", func(t *testing.T) {
		emptyEvent := event.Completion{}
		store := newStoreWithEvent(t, emptyEvent)

		gotChallenge, ok := store.Get(challengeID)

		if !ok {
			t.Fatalf("couldn't find challenge")
		}

		if len(gotChallenge.Events()) != 1 {
			t.Errorf("event not applied: %v", gotChallenge)
		}
	})

	t.Run("completion is applied to challenge", func(t *testing.T) {
		challengeEvent := event.New("", 0).AsCompletion("userID", time.Minute)

		store := newStoreWithEvent(t, challengeEvent)

		gotChallenge, ok := store.Get(challengeID)

		if !ok {
			t.Fatalf("couldn't find challenge")
		}

		if len(gotChallenge.Completions()) != 1 {
			t.Errorf("completion not applied")
		}
	})
}
