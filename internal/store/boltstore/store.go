package boltstore

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"strconv"
	"strings"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/popularity"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/dto"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/port"

	bolt "go.etcd.io/bbolt"
)

const (
	ChallengeBucketID  = "challenges"
	PopularityBucketID = "popularity"
)

// popularityKeyDelim separates the Kind from the ID in a popularity bucket key.
// It also appears inside stage/car IDs, so parsing splits on the first one only.
const popularityKeyDelim = "\x1f"

// TODO - implement some sort of backup like this https://github.com/treeder/bolt-backup/blob/master/backup.go

var _ port.Store = &Store{}

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
		if _, err := tx.CreateBucketIfNotExists([]byte(ChallengeBucketID)); err != nil {
			return fmt.Errorf("creating challenge bucket: %w", err)
		}

		if _, err := tx.CreateBucketIfNotExists([]byte(PopularityBucketID)); err != nil {
			return fmt.Errorf("creating popularity bucket: %w", err)
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

func (s *Store) PutChallenge(ctx context.Context, challengeID string, challenge challenge.Model) error {
	if err := ctx.Err(); err != nil {
		return err
	}

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

func (s *Store) GetChallenge(ctx context.Context, challengeID string) (challenge.Model, error) {
	if err := ctx.Err(); err != nil {
		return challenge.Model{}, err
	}

	dto := dto.Challenge{}

	err := s.db.View(func(tx *bolt.Tx) error {
		buf := tx.Bucket([]byte(ChallengeBucketID)).Get([]byte(challengeID))
		if buf == nil {
			return fmt.Errorf("challenge %s not found", challengeID)
		}

		return gob.NewDecoder(bytes.NewBuffer(buf)).Decode(&dto)
	})

	return dto.ToChallenge(), err
}

func (s *Store) DeleteChallenge(ctx context.Context, challengeID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	err := s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(ChallengeBucketID)).Delete([]byte(challengeID))
	})

	return err
}

func (s *Store) RegisterCompletion(ctx context.Context, challengeID string, completion challenge.Completion) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	err := s.db.Update(func(tx *bolt.Tx) error {
		buf := tx.Bucket([]byte(ChallengeBucketID)).Get([]byte(challengeID))
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

func (s *Store) RegisterVote(ctx context.Context, challengeID string, vote challenge.Vote) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		challenges := tx.Bucket([]byte(ChallengeBucketID))

		buf := challenges.Get([]byte(challengeID))
		if buf == nil {
			return fmt.Errorf("challenge %s not found", challengeID)
		}

		storedDTO := dto.Challenge{}
		if err := gob.NewDecoder(bytes.NewBuffer(buf)).Decode(&storedDTO); err != nil {
			return err
		}

		challenge := storedDTO.ToChallenge()
		deltas := popularity.RegisterVote(&challenge, vote)

		encoded := bytes.Buffer{}
		if err := gob.NewEncoder(&encoded).Encode(dto.FromChallenge(challenge)); err != nil {
			return err
		}
		if err := challenges.Put([]byte(challengeID), encoded.Bytes()); err != nil {
			return err
		}

		return applyPopularityDeltas(tx.Bucket([]byte(PopularityBucketID)), deltas)
	})
}

func (s *Store) Popularity(ctx context.Context) (popularity.Snapshot, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	snapshot := popularity.Snapshot{}

	err := s.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(PopularityBucketID)).ForEach(func(k, v []byte) error {
			key, err := parsePopularityKey(k)
			if err != nil {
				return err
			}

			tally := popularity.Tally{}
			if err := gob.NewDecoder(bytes.NewBuffer(v)).Decode(&tally); err != nil {
				return err
			}

			snapshot[key] = tally
			return nil
		})
	})

	return snapshot, err
}

// applyPopularityDeltas folds each delta into its stored Tally, removing the
// entry when it nets back to zero.
func applyPopularityDeltas(bucket *bolt.Bucket, deltas []popularity.Delta) error {
	for _, delta := range deltas {
		key := popularityKey(delta.Key)

		tally := popularity.Tally{}
		if existing := bucket.Get(key); existing != nil {
			if err := gob.NewDecoder(bytes.NewBuffer(existing)).Decode(&tally); err != nil {
				return err
			}
		}

		tally.Up += delta.Up
		tally.Down += delta.Down

		if tally == (popularity.Tally{}) {
			if err := bucket.Delete(key); err != nil {
				return err
			}
			continue
		}

		encoded := bytes.Buffer{}
		if err := gob.NewEncoder(&encoded).Encode(tally); err != nil {
			return err
		}
		if err := bucket.Put(key, encoded.Bytes()); err != nil {
			return err
		}
	}

	return nil
}

func popularityKey(k popularity.Key) []byte {
	return []byte(strconv.Itoa(int(k.Kind)) + popularityKeyDelim + k.ID)
}

func parsePopularityKey(b []byte) (popularity.Key, error) {
	kind, id, found := strings.Cut(string(b), popularityKeyDelim)
	if !found {
		return popularity.Key{}, fmt.Errorf("malformed popularity key %q", b)
	}

	n, err := strconv.Atoi(kind)
	if err != nil {
		return popularity.Key{}, fmt.Errorf("malformed popularity key %q: %w", b, err)
	}

	return popularity.Key{Kind: popularity.Kind(n), ID: id}, nil
}
