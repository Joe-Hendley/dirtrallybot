// Package buildersession holds the in-progress challenge Config for an open
// challenge builder. State lives only in memory: a builder is a short-lived UI
// interaction, and Discord discards the interaction token after 15 minutes, so
// there is nothing worth persisting.
package buildersession

import (
	"sync"
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
)

// defaultTTL matches Discord's interaction token lifetime. A builder that has
// not been touched for this long can no longer be edited, so its state is dead.
const defaultTTL = 15 * time.Minute

type entry struct {
	config  challenge.Config
	expires time.Time
}

// Store keeps one Config per open builder, keyed by the builder message ID.
type Store struct {
	mu      sync.Mutex
	ttl     time.Duration
	now     func() time.Time
	entries map[string]entry
}

// New returns an empty Store with the default TTL.
func New() *Store {
	return &Store{
		ttl:     defaultTTL,
		now:     time.Now,
		entries: map[string]entry{},
	}
}

// Get returns the stored Config for a builder, or false if there is none or it
// has expired.
func (s *Store) Get(builderID string) (challenge.Config, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	got, ok := s.entries[builderID]
	if !ok {
		return challenge.Config{}, false
	}

	if s.now().After(got.expires) {
		delete(s.entries, builderID)
		return challenge.Config{}, false
	}

	return got.config, true
}

// Put stores the Config for a builder and refreshes its expiry.
func (s *Store) Put(builderID string, config challenge.Config) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.entries[builderID] = entry{config: config, expires: s.now().Add(s.ttl)}
	s.reapLocked()
}

// Delete drops a builder's state, e.g. once its challenge has been created.
func (s *Store) Delete(builderID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.entries, builderID)
}

// reapLocked removes expired entries so an abandoned builder does not leak. Put
// is frequent enough to keep the map bounded without a background goroutine.
func (s *Store) reapLocked() {
	now := s.now()
	for id, e := range s.entries {
		if now.After(e.expires) {
			delete(s.entries, id)
		}
	}
}
