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

## Challenge viewer

While the bot runs it also serves a read-only web page at `WEBADDR` (default
`http://localhost:8080`) listing every stored challenge — date, track, distance,
weather and car. The date comes from the challenge's Discord message ID, which
is a snowflake. Set `WEBADDR=off` to disable it.

Each challenge has a collapsible panel of its completions, sortable by lap time
or submission date. Sorting is driven by [htmx](https://htmx.org), vendored at
`internal/bot/web/static/htmx.min.js` and served by the bot — no external
requests. Completions recorded before the bot stored submitter names and
submission times show a dash in those columns.

## Configuration

Configuration is read from a `.env` file in the working directory (the process
environment takes precedence). A missing `.env` is a fatal error. Copy
[`default.env`](default.env) to `.env` and fill it in:

```
cp default.env .env
```

| Key          | Description                                                                       |
| ------------ | -------------------------------------------------------------------------------- |
| `TOKEN`      | the bot token                                                                     |
| `APP`        | the application ID, used to register the slash commands                           |
| `TESTSERVER` | guild ID where the `!cars` / `!stages` debug commands are allowed                 |
| `RANDOMISER` | `biased` (default), `random`, or `deterministic` — how challenges are filled in   |
| `WEBADDR`    | challenge viewer listen address (default `localhost:8080`); `off` disables it      |

Challenges are stored in a bbolt file (`rallybot.db`) by default.

`biased` leans generation towards stages, cars and weather that have had good
👍 / 👎 feedback; `random` picks uniformly; `deterministic` replays a fixed
sequence (for debugging).

## Running

```
make run      # go run ./cmd/dirtrallybot
make build    # go build ./...
make test     # go test ./...
make generate # regenerate the templ views after editing a .templ file
```

`make generate` needs the templ CLI:
`go install github.com/a-h/templ/cmd/templ@latest`. The generated
`*_templ.go` files are committed, so this is only needed after editing a
`.templ` file.

## Layout

- `cmd/dirtrallybot` — entry point
- `internal/bot` — Discord session wiring and interaction routing
- `internal/bot/handler` — the slash command, component and modal handlers
- `internal/bot/render` — turns domain values into Discord message content
- `internal/bot/web` — the local challenge viewer (templ views)
- `internal/model` — the domain: games, locations, stages, classes, cars, …
- `internal/randomiser` — challenge generation (deterministic, seeded)
- `internal/store` — persistence (`port` interface, `boltstore`, `memorystore`)

See `docs/` for the data model and outstanding work.
