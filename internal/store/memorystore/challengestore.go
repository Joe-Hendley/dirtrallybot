package memorystore

import (
	"fmt"
	"sync"

	"github.com/Joe-Hendley/dirtrallybot/internal/model"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
)

var _ model.Store = &Store{}

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

func (s *Store) PutChallenge(id string, challenge challenge.Model) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.challengeMap[id] = challenge

	return nil
}

func (s *Store) GetChallenge(challengeID string) (c challenge.Model, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	got, ok := s.challengeMap[challengeID]
	if !ok {
		return challenge.Model{}, fmt.Errorf("challenge %s not found", challengeID)
	}
	return got, nil
}

func (s *Store) DeleteChallenge(id string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.challengeMap, id)

	return nil
}

func (s *Store) RegisterCompletion(challengeID string, completion challenge.Completion) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	challenge, ok := s.challengeMap[challengeID]
	if !ok {
		return fmt.Errorf("challenge %s not found", challengeID) //fmt.Errorf("challenge id %s not found", challengeID)
	}

	challenge.RegisterCompletion(completion)

	s.challengeMap[challengeID] = challenge

	return nil
}
