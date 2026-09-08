// Package port defines the storage interface the bot depends on. Implementations
// live under internal/store.
package port

import (
	"context"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/popularity"
)

// Store persists challenges, the completions recorded against them, and the
// running popularity tally that votes feed.
type Store interface {
	PutChallenge(ctx context.Context, challengeID string, challenge challenge.Model) error
	GetChallenge(ctx context.Context, challengeID string) (challenge.Model, error)
	DeleteChallenge(ctx context.Context, challengeID string) error
	RegisterCompletion(ctx context.Context, challengeID string, completion challenge.Completion) error

	// RegisterVote applies a user's vote to a stored challenge and its
	// popularity tally. Re-voting the same way withdraws the vote; voting the
	// other way flips it.
	RegisterVote(ctx context.Context, challengeID string, vote challenge.Vote) error
	// Popularity returns the current tally for every domain item that has been
	// voted on.
	Popularity(ctx context.Context) (popularity.Snapshot, error)
}
