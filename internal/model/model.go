package model

import (
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
)

type Store interface {
	PutChallenge(challengeID string, challenge challenge.Model) error
	GetChallenge(challengeID string) (challenge.Model, error)
	DeleteChallenge(challengeID string) error
	RegisterCompletion(challengeID string, completion challenge.Completion) error
}
