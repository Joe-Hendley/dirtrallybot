// Package popularity holds the running thumbs-up/down tally for each domain item
// a challenge is made of, and turns that tally into a generation weight. A
// biased randomiser reads a Snapshot; the store maintains it as votes arrive.
package popularity

import (
	"math"
	"strconv"
	"strings"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/drivetrain"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
)

// Kind is the sort of domain item a Key identifies.
type Kind uint8

const (
	KindStage Kind = iota + 1
	KindLocation
	KindDistance
	KindWeather
	KindCar
	KindClass
	KindDrivetrain
)

// String is the singular display name for the Kind.
func (k Kind) String() string {
	switch k {
	case KindStage:
		return "Stage"
	case KindLocation:
		return "Location"
	case KindDistance:
		return "Distance"
	case KindWeather:
		return "Weather"
	case KindCar:
		return "Car"
	case KindClass:
		return "Class"
	case KindDrivetrain:
		return "Drivetrain"
	}
	return "invalid kind"
}

// keySeparator joins a name with its parent enum value in a stage or car Key,
// keeping names that repeat across locations or classes distinct.
const keySeparator = "\x1f"

// Key identifies one domain item. ID is a stable string within a Kind: enum
// kinds use the integer value, stages and cars qualify their name with their
// parent so names that repeat across locations/classes stay distinct.
type Key struct {
	Kind Kind
	ID   string
}

// Tally is the raw feedback recorded against a Key.
type Tally struct {
	Up   int
	Down int
}

// Net is thumbs up minus thumbs down.
func (t Tally) Net() int {
	return t.Up - t.Down
}

// Snapshot is the whole popularity picture at a point in time. A missing Key is
// treated as an empty Tally.
type Snapshot map[Key]Tally

// curveScale (k) sets how many net votes meaningfully move a weight.
const curveScale = 5.0

// weightFloor and weightCeil bound the generation weight, so a disliked item is
// rare but never impossible and a loved one never crowds everything else out.
const (
	weightFloor = 0.1
	weightCeil  = 0.9
)

// Weight is the generation weight for a Key: a logistic curve of its net score,
// bounded by weightFloor and weightCeil, sitting at the midpoint when there is
// no feedback.
func (s Snapshot) Weight(k Key) float64 {
	net := float64(s[k].Net())
	return weightFloor + (weightCeil-weightFloor)/(1+math.Exp(-net/curveScale))
}

// LocationKey identifies a location.
func LocationKey(l location.Model) Key {
	return Key{Kind: KindLocation, ID: strconv.Itoa(int(l))}
}

// DistanceKey identifies a stage distance.
func DistanceKey(d stage.Distance) Key {
	return Key{Kind: KindDistance, ID: strconv.Itoa(int(d))}
}

// WeatherKey identifies a weather type.
func WeatherKey(w weather.Model) Key {
	return Key{Kind: KindWeather, ID: strconv.Itoa(int(w))}
}

// ClassKey identifies a car class.
func ClassKey(c class.Model) Key {
	return Key{Kind: KindClass, ID: strconv.Itoa(int(c))}
}

// DrivetrainKey identifies a drivetrain.
func DrivetrainKey(d drivetrain.Model) Key {
	return Key{Kind: KindDrivetrain, ID: strconv.Itoa(int(d))}
}

// StageKey identifies a stage, qualified by its location.
func StageKey(s stage.Model) Key {
	return Key{Kind: KindStage, ID: strconv.Itoa(int(s.Location())) + keySeparator + s.Name()}
}

// CarKey identifies a car, qualified by its class.
func CarKey(c car.Model) Key {
	return Key{Kind: KindCar, ID: strconv.Itoa(int(c.Class())) + keySeparator + c.Name()}
}

// Label is a human-readable description of the item a Key identifies. It
// reverses the encoding applied by the *Key constructors, falling back to the
// raw ID for a Kind or ID it cannot decode.
func (k Key) Label() string {
	switch k.Kind {
	case KindStage:
		if loc, name, ok := splitQualified(k.ID); ok {
			return location.Model(loc).String() + " » " + name
		}
	case KindCar:
		if cls, name, ok := splitQualified(k.ID); ok {
			return name + " (" + class.Model(cls).String() + ")"
		}
	case KindLocation:
		if n, err := strconv.Atoi(k.ID); err == nil {
			return location.Model(n).String()
		}
	case KindDistance:
		if n, err := strconv.Atoi(k.ID); err == nil {
			return stage.Distance(n).String()
		}
	case KindWeather:
		if n, err := strconv.Atoi(k.ID); err == nil {
			return weather.Model(n).String()
		}
	case KindClass:
		if n, err := strconv.Atoi(k.ID); err == nil {
			return class.Model(n).String()
		}
	case KindDrivetrain:
		if n, err := strconv.Atoi(k.ID); err == nil {
			return drivetrain.Model(n).String()
		}
	}
	return k.ID
}

// splitQualified separates a name from the parent enum value prefixed to it by
// StageKey or CarKey.
func splitQualified(id string) (parent int, name string, ok bool) {
	before, after, found := strings.Cut(id, keySeparator)
	if !found {
		return 0, "", false
	}
	parent, err := strconv.Atoi(before)
	if err != nil {
		return 0, "", false
	}
	return parent, after, true
}

// Delta is a change to apply to one Key's stored Tally.
type Delta struct {
	Key  Key
	Up   int
	Down int
}

// RegisterVote applies vote to c and reports the tally changes it implies.
//
// A first vote is recorded; voting the same way again withdraws it; voting the
// other way flips it. The returned deltas cover every constituent item of the
// challenge (see KeysFor) and are empty when nothing changed.
func RegisterVote(c *challenge.Model, vote challenge.Vote) []Delta {
	existing, voted := c.VoteFor(vote.UserID())
	keys := KeysFor(*c)

	var up, down int
	switch {
	case !voted:
		c.SetVote(vote)
		up, down = contribution(vote.Sentiment(), 1)
	case existing.Sentiment() == vote.Sentiment():
		c.WithdrawVote(vote.UserID())
		up, down = contribution(existing.Sentiment(), -1)
	default:
		c.SetVote(vote)
		outUp, outDown := contribution(existing.Sentiment(), -1)
		inUp, inDown := contribution(vote.Sentiment(), 1)
		up, down = outUp+inUp, outDown+inDown
	}

	if up == 0 && down == 0 {
		return nil
	}

	deltas := make([]Delta, 0, len(keys))
	for _, key := range keys {
		deltas = append(deltas, Delta{Key: key, Up: up, Down: down})
	}
	return deltas
}

// contribution is the tally change one vote of the given sentiment makes, scaled
// by sign (+1 to add the vote, -1 to remove it).
func contribution(s challenge.Sentiment, sign int) (up, down int) {
	switch s {
	case challenge.Up:
		return sign, 0
	case challenge.Down:
		return 0, sign
	}
	return 0, 0
}

// KeysFor decomposes a challenge into every domain item a vote on it should
// count towards.
func KeysFor(c challenge.Model) []Key {
	s := c.Stage()
	vehicle := c.Car()

	return []Key{
		StageKey(s),
		LocationKey(s.Location()),
		DistanceKey(s.Distance()),
		WeatherKey(c.Weather()),
		CarKey(vehicle),
		ClassKey(vehicle.Class()),
		DrivetrainKey(vehicle.Class().Drivetrain()),
	}
}
