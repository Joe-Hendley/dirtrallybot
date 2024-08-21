package event

import "time"

type baseEvent struct {
	id        string
	timestamp time.Time
}

func New(id string, timestamp int64) baseEvent {
	return baseEvent{
		id:        id,
		timestamp: time.Unix(timestamp, 0),
	}
}

func (b baseEvent) ID() string {
	return b.id
}

func (b baseEvent) Timestamp() time.Time {
	return b.timestamp
}

type Completion struct {
	baseEvent
	userID      string
	challengeID string
	duration    time.Duration
}

// Completion holds the details of a user-reported challenge completion
func (b baseEvent) AsCompletion(userID string, duration time.Duration) Completion {
	return Completion{
		baseEvent: b,
		userID:    userID,
		duration:  duration,
	}
}

func (c Completion) UserID() string {
	return c.userID
}

func (c Completion) Duration() time.Duration {
	return c.duration
}

func (c Completion) ChallengeID() string {
	return c.challengeID
}
