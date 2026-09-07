# dirtrallybot

A Discord bot that generates random rally challenges — a stage, weather and car —
for [DiRT Rally 2.0](https://dirtrally2.dirtgame.com/) and
[EA Sports WRC](https://www.ea.com/games/ea-sports-wrc), and tracks the times
people set against each challenge.

## Commands

- `/newstage-dr2` — open the challenge builder for DiRT Rally 2.0
- `/newstage-wrc` — open the challenge builder for WRC

The builder lets you pin any of location, distance, stage, weather, drivetrain,
class or car, or leave each on "random". Submitting posts the challenge to the
channel with buttons to enter a time (⏱️) and to list everyone's times (📋).

## Configuration

Configuration is read from a `.env` file in the working directory (the process
environment takes precedence). A missing `.env` is a fatal error.

| Key          | Description                                                        |
| ------------ | ----------------------------------------------------------------- |
| `token`      | the bot token                                                     |
| `app`        | the application ID, used to register the slash commands           |
| `testserver` | guild ID where the `!cars` / `!stages` debug commands are allowed |

Challenges are stored in a bbolt file (`rallybot.db`) by default.

## Running

```
make run      # go run ./cmd/dirtrallybot
make build    # go build ./...
make test     # go test ./...
```

## Layout

- `cmd/dirtrallybot` — entry point
- `internal/bot` — Discord session wiring and interaction routing
- `internal/bot/handler` — the slash command, component and modal handlers
- `internal/bot/render` — turns domain values into Discord message content
- `internal/model` — the domain: games, locations, stages, classes, cars, …
- `internal/randomiser` — challenge generation (deterministic, seeded)
- `internal/store` — persistence (`port` interface, `boltstore`, `memorystore`)

See `docs/` for the data model and the refactor checklist.
