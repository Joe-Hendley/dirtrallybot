package memorystore

import (
	"sync"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
)

var DefaultStore = Store{
	lock:         &sync.Mutex{},
	challengeMap: map[string]challenge.Model{},
}

type Store struct {
	lock         *sync.Mutex
	challengeMap map[string]challenge.Model
}

func (s *Store) Get(id string) (challenge challenge.Model, ok bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	challenge, ok = s.challengeMap[id]
	return challenge, ok
}

func (s *Store) Put(id string, challenge challenge.Model) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.challengeMap[id] = challenge
}

func (s *Store) Delete(id string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.challengeMap, id)
}
