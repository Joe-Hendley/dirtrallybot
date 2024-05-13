package store

import (
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/internal/memorystore"
)

var store = memorystore.DefaultStore

func Get(id string) (c challenge.Model, ok bool) {
	return store.Get(id)
}

func Put(id string, challenge *challenge.Model) {
	store.Put(id, challenge)
}

func ApplyEvent(id string, event any) error {
	return store.ApplyEvent(id, event)
}

func Delete(id string) {
	store.Delete(id)
}
