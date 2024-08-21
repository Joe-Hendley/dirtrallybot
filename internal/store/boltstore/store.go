package boltstore

import (
	"fmt"

	"github.com/Joe-Hendley/dirtrallybot/internal/model"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"

	bolt "go.etcd.io/bbolt"
)

// TODO - implement some sort of backup like this https://github.com/treeder/bolt-backup/blob/master/backup.go

var _ model.Store = &Store{}

type Store struct {
	db *bolt.DB
}

func New() (*Store, error) {
	db, err := bolt.Open("rallybot.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("challenges"))
		if err != nil {
			return fmt.Errorf("creating challenge bucket: %w", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte("feedback"))
		if err != nil {
			return fmt.Errorf("creating feedback bucket: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) ApplyEvent(id string, event model.Event) error {
	panic("unimplemented")
}

func (s *Store) Delete(id string) {
	panic("unimplemented")
}

func (s *Store) Get(id string) (c challenge.Model, ok bool) {
	panic("unimplemented")
}

func (s *Store) Put(id string, challenge *challenge.Model) {
	panic("unimplemented")
}
