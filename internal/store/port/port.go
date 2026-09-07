// Package port defines the storage interface the bot depends on. Implementations
// live under internal/store.
package port

import (
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
)

// Store persists challenges and the completions recorded against them.
type Store interface {
	PutChallenge(challengeID string, challenge challenge.Model) error
	GetChallenge(challengeID string) (challenge.Model, error)
	DeleteChallenge(challengeID string) error
	RegisterCompletion(challengeID string, completion challenge.Completion) error
}
