// Package port defines the storage interface the bot depends on. Implementations
// live under internal/store.
package port

import (
	"context"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
)

// Store persists challenges and the completions recorded against them.
type Store interface {
	PutChallenge(ctx context.Context, challengeID string, challenge challenge.Model) error
	GetChallenge(ctx context.Context, challengeID string) (challenge.Model, error)
	DeleteChallenge(ctx context.Context, challengeID string) error
	RegisterCompletion(ctx context.Context, challengeID string, completion challenge.Completion) error
}
