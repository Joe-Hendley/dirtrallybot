package memorystore

import (
	"context"
	"fmt"
	"sync"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/popularity"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/dto"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/port"
)

var _ port.Store = &Store{}

// Store keeps challenges in memory. It holds DTOs rather than domain values so
// that reads and writes are copies, matching the persistence guarantees of the
// bolt store.
type Store struct {
	lock         sync.Mutex
	challengeMap map[string]dto.Challenge
	popularity   map[popularity.Key]popularity.Tally
}

func New() *Store {
	return &Store{
		challengeMap: map[string]dto.Challenge{},
		popularity:   map[popularity.Key]popularity.Tally{},
	}
}

func (s *Store) PutChallenge(ctx context.Context, id string, c challenge.Model) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.challengeMap[id] = dto.FromChallenge(c)

	return nil
}

func (s *Store) GetChallenge(ctx context.Context, challengeID string) (challenge.Model, error) {
	if err := ctx.Err(); err != nil {
		return challenge.Model{}, err
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	got, ok := s.challengeMap[challengeID]
	if !ok {
		return challenge.Model{}, fmt.Errorf("challenge %s not found", challengeID)
	}

	return got.ToChallenge(), nil
}

func (s *Store) DeleteChallenge(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.challengeMap, id)

	return nil
}

func (s *Store) RegisterCompletion(ctx context.Context, challengeID string, completion challenge.Completion) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	stored, ok := s.challengeMap[challengeID]
	if !ok {
		return fmt.Errorf("challenge %s not found", challengeID)
	}

	updated := stored.ToChallenge()
	updated.RegisterCompletion(completion)
	s.challengeMap[challengeID] = dto.FromChallenge(updated)

	return nil
}

func (s *Store) RegisterVote(ctx context.Context, challengeID string, vote challenge.Vote) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	stored, ok := s.challengeMap[challengeID]
	if !ok {
		return fmt.Errorf("challenge %s not found", challengeID)
	}

	updated := stored.ToChallenge()
	deltas := popularity.RegisterVote(&updated, vote)
	s.challengeMap[challengeID] = dto.FromChallenge(updated)

	for _, delta := range deltas {
		tally := s.popularity[delta.Key]
		tally.Up += delta.Up
		tally.Down += delta.Down
		if tally == (popularity.Tally{}) {
			delete(s.popularity, delta.Key)
			continue
		}
		s.popularity[delta.Key] = tally
	}

	return nil
}

func (s *Store) Popularity(ctx context.Context) (popularity.Snapshot, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	snapshot := make(popularity.Snapshot, len(s.popularity))
	for key, tally := range s.popularity {
		snapshot[key] = tally
	}

	return snapshot, nil
}
