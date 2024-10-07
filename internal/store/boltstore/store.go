package boltstore

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/Joe-Hendley/dirtrallybot/internal/model"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/boltstore/internal/dto"

	bolt "go.etcd.io/bbolt"
)

const ChallengeBucketID = "challenges"

// TODO - implement some sort of backup like this https://github.com/treeder/bolt-backup/blob/master/backup.go

var _ model.Store = &Store{}

type Store struct {
	db *bolt.DB
}

// "rallybot.db"
func New(filename string) (*Store, error) {
	db, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(ChallengeBucketID))
		if err != nil {
			return fmt.Errorf("creating challenge bucket: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) PutChallenge(challengeID string, challenge challenge.Model) error {
	dto := dto.FromChallenge(challenge)

	err := s.db.Update(func(tx *bolt.Tx) error {
		encoded := bytes.Buffer{}
		err := gob.NewEncoder(&encoded).Encode(dto)
		if err != nil {
			return err
		}

		err = tx.Bucket([]byte(ChallengeBucketID)).Put([]byte(challengeID), encoded.Bytes())
		return err
	})

	return err
}

func (s *Store) GetChallenge(challengeID string) (challenge.Model, error) {
	dto := dto.Challenge{}

	err := s.db.View(func(tx *bolt.Tx) error {
		buf := tx.Bucket([]byte("challenges")).Get([]byte(challengeID))
		if buf == nil {
			return fmt.Errorf("challenge %s not found", challengeID)
		}

		return gob.NewDecoder(bytes.NewBuffer(buf)).Decode(&dto)
	})

	return dto.ToChallenge(), err
}

func (s *Store) DeleteChallenge(challengeID string) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("challenges")).Delete([]byte(challengeID))
	})

	return err
}

func (s *Store) RegisterCompletion(challengeID string, completion challenge.Completion) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		buf := tx.Bucket([]byte("challenges")).Get([]byte(challengeID))
		if buf == nil {
			return fmt.Errorf("challenge %s not found", challengeID)
		}

		originalDTO := dto.Challenge{}
		err := gob.NewDecoder(bytes.NewBuffer(buf)).Decode(&originalDTO)
		if err != nil {
			return err
		}

		challenge := originalDTO.ToChallenge()
		challenge.RegisterCompletion(completion)
		newDTO := dto.FromChallenge(challenge)

		encoded := bytes.Buffer{}
		err = gob.NewEncoder(&encoded).Encode(newDTO)
		if err != nil {
			return err
		}

		err = tx.Bucket([]byte(ChallengeBucketID)).Put([]byte(challengeID), encoded.Bytes())
		return err
	})

	return err
}
