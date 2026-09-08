# Data Model

## What the bot stores today

The bot is not fully stateless: button interactions on a posted challenge need
the challenge's details to still be available when someone later submits a time.

The `port.Store` interface (`internal/store/port`) is the whole persistence
surface:

- `PutChallenge` — save a generated challenge
- `GetChallenge` — load one
- `DeleteChallenge` — remove one
- `RegisterCompletion` — append a completion time to a stored challenge
- `RegisterVote` — apply a 👍 / 👎 to a stored challenge and the popularity tally
- `Popularity` — the current tally for every domain item that has been voted on

A challenge is keyed by the Discord message ID of the message the bot posts for
it. That ID is a snowflake, so it also encodes a creation timestamp.

### Stored shape

Challenges are stored as a versioned DTO (`internal/store/dto`), not as the
domain type directly, so the on-disk format can evolve independently of the
model. Each DTO carries a `Version` field.

A stored challenge holds:

- the stage (name, location, distance)
- the weather
- the car (name, class)
- the list of completions, each `{userID, duration}`
- the list of votes, each `{userID, sentiment}` — at most one per user

`RegisterCompletion` is read-modify-write: load the challenge, append the
completion, write it back.

### Feedback and popularity

`RegisterVote` records one user's 👍 / 👎 on a challenge. Re-voting the same way
withdraws the vote; voting the other way flips it.

Every vote counts towards **all** of the challenge's constituent domain items —
its stage, location, distance, weather, car, class and drivetrain
(`popularity.KeysFor`). The store keeps a running `{up, down}` tally per item,
adjusted on each vote and its reversal, so no full scan of challenges is needed.
An item whose tally nets back to zero is dropped.

`Popularity` returns the whole tally as a `popularity.Snapshot`. The `Biased`
randomiser (the default; `randomiser` config key) turns each item's net score
`s` into a generation weight on a logistic curve,
`weight(s) = 0.1 + 0.8 / (1 + exp(-s/5))`: an unrated item sits at `0.5`, a
loved one approaches `0.9`, a disliked one `0.1`. Nothing is ever excluded, and
with no feedback yet every weight is `0.5`, so generation starts out uniform.

### Implementations

- `boltstore` — a bbolt file, a `challenges` bucket and a `popularity` bucket
  (one gob-encoded `Tally` per item key), gob-encoded DTOs. This is the default.
- `memorystore` — a map of the same challenge DTOs plus a `map[Key]Tally`, for
  tests and local runs. It holds DTOs rather than domain values so that reads
  and writes copy, matching bolt.

The challenge builder's in-progress state is separate again: it lives only in
memory (`internal/store/buildersession`) with a 15-minute TTL and never reaches
`port.Store`.

## Not built yet

- **Granular feedback.** A vote is a single 👍 / 👎 on the whole challenge; there
  is no way to say which part was liked or disliked.
- **Vote decay.** Old votes count as much as fresh ones; popularity never ages.

## A possible future direction: an event log

An earlier design treated the store as an append-only event log rather than a
mutable blob:

- `Creation` — the challenge details
- `Completion` — `{userID, displayName, challengeID, duration}`
- `Feedback` — `{userID, challengeID, feedback}`

Each event would be keyed by the ID of the interaction that produced it. On
restart the bot would replay events in snowflake order to rebuild state; order
only really matters for keeping a challenge's later events after its creation.

Feedback and the popularity tally are instead built on the existing
read-modify-write blob (see above). An event log would keep the raw votes
addressable — useful for decay or recomputing the weighting — at the cost of
replay logic and log compaction. It is not what the code does today.
