package model

import (
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
)

type Event interface {
	ID() string
	Timestamp() time.Time
	ChallengeID() string
}

type Store interface {
	Get(id string) (c challenge.Model, ok bool)
	Put(id string, challenge *challenge.Model)
	ApplyEvent(id string, event Event) error
	Delete(id string)
}
