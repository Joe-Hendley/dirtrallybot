package event

import "time"

type base struct {
	id          string // the message ID might not be unique
	timestamp   time.Time
	challengeID string // is this needed?
}

func NewEvent(id string) base {
	return base{
		id:        id,
		timestamp: time.Now().UTC(),
	}
}

func (b base) ID() string {
	return b.id
}

func (b base) Timestamp() time.Time {
	return b.timestamp
}

func (b base) ChallengeID() string {
	return b.challengeID
}

type completion struct {
	base
	userID   string
	duration time.Duration
}

// Completion holds the details of a user-reported challenge completion
func (b base) Completion(userID string, duration time.Duration) completion {
	return completion{
		base:     b,
		userID:   userID,
		duration: duration,
	}
}

func (c completion) UserID() string {
	return c.userID
}

func (c completion) Duration() time.Duration {
	return c.duration
}
