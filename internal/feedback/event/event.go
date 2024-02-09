package event

import "time"

type BaseModel struct {
	id          string
	timestamp   time.Time
	challengeID string // is this needed?
}

func (bm BaseModel) ID() string {
	return bm.id
}

func (bm BaseModel) Timestamp() time.Time {
	return bm.timestamp
}

func (bm BaseModel) ChallengeID() string {
	return bm.challengeID
}

type CompletionEvent struct {
	BaseModel
}
