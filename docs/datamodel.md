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

`RegisterCompletion` is read-modify-write: load the challenge, append the
completion, write it back.

### Implementations

- `boltstore` — a bbolt file, one `challenges` bucket, gob-encoded DTOs. This is
  the default.
- `memorystore` — a map of the same DTOs, for tests and local runs. It holds
  DTOs rather than domain values so that reads and writes copy, matching bolt.

The challenge builder's in-progress state is separate again: it lives only in
memory (`internal/store/buildersession`) with a 15-minute TTL and never reaches
`port.Store`.

## Not built yet

- **Feedback.** The 👍 / 👎 buttons on a posted challenge are disabled. The
  intent was to record whether a car / stage / weather combination was good or
  bad and feed that back into generation.
- **Generation modifiers.** There is no persisted or in-memory set of weightings
  that biases future challenge generation.

## A possible future direction: an event log

An earlier design treated the store as an append-only event log rather than a
mutable blob:

- `Creation` — the challenge details
- `Completion` — `{userID, displayName, challengeID, duration}`
- `Feedback` — `{userID, challengeID, feedback}`

Each event would be keyed by the ID of the interaction that produced it. On
restart the bot would replay events in snowflake order to rebuild state; order
only really matters for keeping a challenge's later events after its creation.

This would make feedback and modifiers natural to add, at the cost of replay
logic and log compaction. It is not what the code does today.
