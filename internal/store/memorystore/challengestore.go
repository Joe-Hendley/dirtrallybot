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
	challengeMap map[string]*challenge.Model
}

func New() *Store {
	return &Store{
		lock:         &sync.Mutex{},
		challengeMap: map[string]*challenge.Model{},
	}
}

func (s *Store) Get(id string) (c challenge.Model, ok bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	got, ok := s.challengeMap[id]
	if !ok {
		return challenge.Model{}, false
	}
	return *got, ok
}

func (s *Store) Put(id string, challenge *challenge.Model) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.challengeMap[id] = challenge
}

func (s *Store) ApplyEvent(id string, event model.Event) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, ok := s.challengeMap[id]
	if !ok {
		return fmt.Errorf("challenge id %s not found", id)
	}

	return s.challengeMap[id].ApplyEvent(event)
}

func (s *Store) Delete(id string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.challengeMap, id)
}
