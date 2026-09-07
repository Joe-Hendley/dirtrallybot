# Refactor Checklist

A prioritised list of refactoring work for the bot. Ordered so that earlier
items unblock later ones (in particular, the builder rework in section 1 makes
the domain and presentation split in section 3 far easier to test).

Tick items off as they land. Each item notes the primary files involved.

## 1. Challenge builder state

The builder reconstructs its in-progress state from Discord component custom IDs
*and* by parsing its own rendered message text. Display-string changes silently
break state reconstruction, and the `-` delimiter is never escaped.

- [x] Persist the in-progress `challenge.Config` server-side, keyed by the
      builder message ID, rather than encoding it into component custom IDs and
      message content. New in-memory `buildersession.Store` with a 15-minute TTL
      (`internal/store/buildersession/store.go`,
      `internal/bot/handler/challenge/challengebuilder.go`)
- [x] Remove the message-content parsing in `buildCarConfigFromInteraction`,
      including the stage-distance reverse-engineering that splits the stage
      name on spaces
      (`internal/bot/handler/challenge/challengebuilder.go`)
- [x] Collapse `buildStageConfigFromInteraction` and
      `buildCarConfigFromInteraction` into a single `configFromInteraction`
      load-apply-save path. Cascade-clearing added to the `applyX` helpers so a
      stale downstream selection (e.g. a stage from a different location) is
      dropped, and "random" now clears its own field
      (`internal/bot/handler/challenge/challengebuilder.go`)
- [ ] Extract the repeated select-menu construction
      (`random` option, option loop, `hasDefault` fallback) shared by
      `buildLocationsMenu`, `buildDistanceMenu`, `buildStageMenu`,
      `buildWeatherMenu`, `buildDriveTrainMenu`, `buildClassMenu`, `buildCarMenu`
- [x] Merge `HandleNewDR2Challenge` and `HandleNewWRCChallenge` into
      `HandleNewChallenge`, which derives the game from the command name. Fixes
      the WRC builder header, which read "Dirt Rally 2"
      (`internal/bot/handler/challenge/challenge.go`)
- [x] Delete the `// TODO: rewrite the entire custom ID system` comment once done
      (`internal/bot/handler/challenge/challengebuilder.go`)

## 2. Game abstraction

Every model package switches on `game.Model` into hardcoded twin functions
(`listDR2`/`listWRC`, `inClassDR2`/`inClassWRC`, ...). `stage.AtLocation` is a
single ~475-line switch. Adding a game means editing every package.

- [ ] Move the catalogue data (locations, stages, classes, cars, drivetrains,
      weather) out of Go source into embedded data files loaded via `//go:embed`
      into a registry
- [ ] Reduce the per-package `switch g` dispatch to generic registry lookups
- [ ] Finish the WRC car catalogue and remove `// TODO - this`
      (`internal/model/car/car.go`)
- [ ] Consolidate the three parallel game-to-string mappings — `gameIDString`,
      `gameFromID`, and `game.Model.String()` — into one owned by the `game`
      package
      (`internal/bot/handler/challenge/challengebuilder.go`,
      `internal/model/game/game.go`)

## 3. Hexagonal boundaries

- [ ] Move Discord presentation (emoji, Markdown) out of the domain: have the
      domain expose a neutral view model and render it in a Discord adapter.
      Affects `challenge.Model.FancyString`, `Config.FancyStageString`,
      `Config.FancyCarString`, and the `FancyString`/`FancyString`-style methods
      on `car`, `stage`, `drivetrain`
- [ ] Replace `config`'s `init()` and `os.Exit(1)` with an explicit
      `config.Load() (Config, error)` called from `main`, so importers and tests
      no longer require a `.env` file
      (`internal/config/config.go`)
- [ ] Make command registration and cleanup return errors instead of calling
      `os.Exit(1)` from library code; `bot.New` should not be able to kill the
      process
      (`internal/bot/commands.go`, `internal/bot/bot.go`)
- [ ] Rename the `internal/model` package (it holds only the store port) to
      something like `internal/store/port`, or fold the interface into
      `internal/store`; reduce the number of types named `Model` and drop the
      `challengeModel` import alias
- [ ] Thread `context.Context` through the store methods and Discord calls, and
      drive shutdown from a cancelled context
      (`internal/model/model.go`, `internal/store/...`, `cmd/main.go`)

## 4. Store layer

- [ ] Route both stores through the versioned DTO so serialisation guarantees
      match; `memorystore` currently stores domain objects directly
      (`internal/store/memorystore/store.go`,
      `internal/store/boltstore/internal/dto/dto.go`)
- [ ] Reconcile `docs/datamodel.md` (an event log keyed by snowflake) with the
      implementation (a mutable challenge blob with read-modify-write): pick one
      model and align the other
- [ ] Implement challenge feedback (the good/bad events in `datamodel.md`); the
      feedback buttons are hardcoded `Disabled: true`
      (`internal/bot/handler/challenge/challenge.go`)
- [ ] Fix the `memorystore` aliasing: `challenge.Model` is copied by value but
      shares the `completions` slice backing array, so `RegisterCompletion`'s
      `append` can mutate the stored copy
      (`internal/store/memorystore/store.go`,
      `internal/model/challenge/challenge.go`)
- [ ] Use the `ChallengeBucketID` constant in `GetChallenge`, `DeleteChallenge`
      and `RegisterCompletion` instead of the `"challenges"` literal
      (`internal/store/boltstore/store.go`)

## 5. Correctness bugs

- [ ] `NewRandomChallenge` ignores `Config.Distance` when randomising a stage,
      so choosing a distance plus a random stage gives any length. Add
      `StageOfDistance` to the `Randomiser` interface (it already exists on
      `Simple`) and use it
      (`internal/model/challenge/challenge.go`,
      `internal/randomiser/simple.go`)
- [ ] WRC challenges use the DR2 randomiser — `DR2Randomiser` is the only one
      wired in
      (`internal/bot/handler/challenge/challenge.go`)
- [ ] `bot.New` never sets the `cfg` field, so `Shutdown` and the
      `!cars` / `!stages` test-server gate run against a zero `Config`
      (`internal/bot/bot.go`)
- [ ] `timestamp.Format` is not the inverse of `Parse`: milliseconds use `%d`
      (5ms renders as `.5`, not `.005`) and minutes are space-padded via `%2.f`.
      Add a round-trip test
      (`internal/model/timestamp/timestamp.go`)
- [ ] Standardise on `slog`; remove `log.Printf`, and fix printf directives
      passed to structured calls, e.g. `slog.Error("starting session: %w", ...)`
      and `slog.Error("editing challenge id: %s : %v\n", ...)`
      (`cmd/main.go`, `internal/bot/handler/completion/completion.go`)
- [ ] Surface failures to the user where handlers currently only log and return

## 6. Randomiser naming

The `Simple` randomiser seeds `rand.NewPCG(0, 0)` deliberately — challenge
sequences are meant to be deterministic and reproducible. The name does not
convey that.

- [ ] Rename `Simple` / `NewSimple` to something that signals determinism
      (e.g. `Deterministic`, `Seeded`, `Fixed`) and document why the seed is
      fixed
      (`internal/randomiser/simple.go`, callers in
      `internal/bot/handler/challenge/challenge.go`)

## 7. Project layout and tooling

- [ ] Move `cmd/main.go` (`package cmd`) to `cmd/dirtrallybot/main.go` so the
      binary is named `dirtrallybot`
- [ ] Fix the `build:` target indentation in the `Makefile` (spaces, not a tab)
- [ ] `bot.New` returns an unexported `*bot` from an exported function
      (`internal/bot/bot.go`)
- [ ] Expand the README; consider a conventional `DISCORD_TOKEN`-style env key
      rather than lowercase `token`

## 8. Test coverage

- [x] Cover the builder config-from-interaction parsing (the riskiest code)
      before reworking it, as a safety net
      (`internal/bot/handler/challenge/challengebuilder_test.go`)
- [ ] Add tests for the stores, the randomiser, and `config`
