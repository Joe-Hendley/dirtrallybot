package event

import "time"

type base struct {
	id          string // the message ID might not be unique
	timestamp   time.Time
	challengeID string // is this needed?
}

func New(id string, timestamp int64) base {
	return base{
		id:        id,
		timestamp: time.Unix(timestamp, 0),
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

type Completion struct {
	base
	userID          string
	userDisplayName string
	duration        time.Duration
}

// Completion holds the details of a user-reported challenge completion
func (b base) AsCompletion(userID, userDisplayName string, duration time.Duration) Completion {
	return Completion{
		base:            b,
		userID:          userID,
		userDisplayName: userDisplayName,
		duration:        duration,
	}
}

func (c Completion) UserID() string {
	return c.userID
}

func (c Completion) UserDisplayName() string {
	return c.userDisplayName
}

func (c Completion) Duration() time.Duration {
	return c.duration
}
