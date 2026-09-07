package memorystore

import (
	"context"
	"fmt"
	"sync"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/port"
)

var _ port.Store = &Store{}

type Store struct {
	lock         *sync.Mutex
	challengeMap map[string]challenge.Model
}

func New() *Store {
	return &Store{
		lock:         &sync.Mutex{},
		challengeMap: map[string]challenge.Model{},
	}
}

func (s *Store) PutChallenge(ctx context.Context, id string, challenge challenge.Model) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.challengeMap[id] = challenge

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

	return got, nil
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

	stored.RegisterCompletion(completion)
	s.challengeMap[challengeID] = stored

	return nil
}
