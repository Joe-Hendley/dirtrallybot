package buildersession

import (
	"testing"
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	t.Run("get on an unknown builder returns false", func(t *testing.T) {
		store := New()

		_, ok := store.Get("nope")

		assert.False(t, ok)
	})

	t.Run("put then get round-trips the config", func(t *testing.T) {
		store := New()
		loc := location.SCO
		want := challenge.Config{Game: game.DR2, Location: &loc}

		store.Put("builder-1", want)
		got, ok := store.Get("builder-1")

		require.True(t, ok)
		assert.Equal(t, want, got)
	})

	t.Run("delete removes the config", func(t *testing.T) {
		store := New()
		store.Put("builder-1", challenge.Config{Game: game.DR2})

		store.Delete("builder-1")
		_, ok := store.Get("builder-1")

		assert.False(t, ok)
	})

	t.Run("expired entries are not returned", func(t *testing.T) {
		now := time.Now()
		store := New()
		store.now = func() time.Time { return now }

		store.Put("builder-1", challenge.Config{Game: game.DR2})

		now = now.Add(defaultTTL + time.Second)
		_, ok := store.Get("builder-1")

		assert.False(t, ok)
	})

	t.Run("put refreshes the expiry", func(t *testing.T) {
		now := time.Now()
		store := New()
		store.now = func() time.Time { return now }

		store.Put("builder-1", challenge.Config{Game: game.DR2})
		now = now.Add(defaultTTL - time.Minute)
		store.Put("builder-1", challenge.Config{Game: game.WRC})

		now = now.Add(2 * time.Minute) // past the first expiry, before the second
		got, ok := store.Get("builder-1")

		require.True(t, ok)
		assert.Equal(t, game.WRC, got.Game)
	})

	t.Run("put reaps other expired entries", func(t *testing.T) {
		now := time.Now()
		store := New()
		store.now = func() time.Time { return now }

		store.Put("stale", challenge.Config{Game: game.DR2})
		now = now.Add(defaultTTL + time.Second)
		store.Put("fresh", challenge.Config{Game: game.WRC})

		store.mu.Lock()
		_, staleKept := store.entries["stale"]
		store.mu.Unlock()

		assert.False(t, staleKept)
	})
}
