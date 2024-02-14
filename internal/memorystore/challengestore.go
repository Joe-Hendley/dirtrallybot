package memorystore

import (
	"fmt"
	"sync"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
)

// TODO: remove the exported var, replace it with NewStore()
var DefaultStore = Store{
	lock:         &sync.Mutex{},
	challengeMap: map[string]*challenge.Model{},
}

type Store struct {
	lock         *sync.Mutex
	challengeMap map[string]*challenge.Model
}

// Get returns a shallow copy of the stored challenge if present.
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

func (s *Store) ApplyEvent(id string, event any) error { // TODO EVENTTYPE
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
