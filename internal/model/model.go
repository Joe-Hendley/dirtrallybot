package model

import (
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
)

type Event interface {
	ID() string
	Timestamp() time.Time
}

type Store interface {
	Get(id string) (c challenge.Model, ok bool)
	Put(id string, challenge *challenge.Model)
	ApplyEvent(id string, event Event) error
	Delete(id string)
}

// TODO - split out store into Challenge & completion / Feedback interfaces

type ChallengeStore interface {
	GetChallenge(id string) (c challenge.Model, ok bool)
	AddChallenge(id string, challenge *challenge.Model)
	AddCompletion(id string, event Event) error
	DeleteChallenge(id string) bool
}

// I'm thinking this should add feedback, get the net feedback, get all feedback, delete individual feedback
// might be more suited to a relational db than bolt but that's a problem for later
type FeedbackStore interface {
	AddFeedback(feedback any)
	GetNetFeedback() any
	GetAllFeedback() []any
	DeleteFeedback(id string) bool
}
