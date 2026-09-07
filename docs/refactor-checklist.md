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
- [x] Extract the repeated select-menu construction
      (`random` option, option loop, `hasDefault` fallback) into `buildSelectMenu`
      over a `menuEntry` slice. Each builder now just maps its domain list to
      entries; the weather one-option case stays bespoke
      (`internal/bot/handler/challenge/challengebuilder.go`)
- [x] Merge `HandleNewDR2Challenge` and `HandleNewWRCChallenge` into
      `HandleNewChallenge`, which derives the game from the command name. Fixes
      the WRC builder header, which read "Dirt Rally 2"
      (`internal/bot/handler/challenge/challenge.go`)
- [x] Delete the `// TODO: rewrite the entire custom ID system` comment once done
      (`internal/bot/handler/challenge/challengebuilder.go`)

## 2. Game abstraction

Every model package switched on `game.Model` into hardcoded twin functions
(`listDR2`/`listWRC`, `inClassDR2`/`inClassWRC`, ...). `stage.AtLocation` was a
single ~475-line switch. Adding a game meant editing every package.

- [x] Replace the `switch g { xDR2 / xWRC }` twin-function dispatch with
      package-level data tables: `location.byGame`, `stage.byLocation`,
      `class.byGame`, `class.byGameDrivetrain`, `car.namesByGameClass`,
      `drivetrain.byGame`. Lookups are one map access; adding a game is a data
      edit. Catalogue-integrity tests added (every listed class has cars, every
      listed location has stages). Moving the tables to embedded data files is
      tracked separately in `docs/catalogue-data-migration.md`
- [x] Finish the WRC car catalogue and remove `// TODO - this`. The data was in
      fact complete for all 18 WRC classes; the stale TODO is gone and a test
      now enforces completeness
      (`internal/model/car/car.go`)
- [x] Consolidate the game slug mapping. `game.Model.ID()` and `game.FromID()`
      now own the `dr2`/`wrc` <-> `Model` mapping; `gameIDString` and `gameFromID`
      in the handler are gone, and a test pins the handler's `DR2ID`/`WRCID`
      fragments to the game package
      (`internal/model/game/game.go`,
      `internal/bot/handler/challenge/challengebuilder.go`)

## 3. Hexagonal boundaries

- [x] Move Discord presentation (emoji, flags, Markdown) out of the domain into
      a new `internal/bot/render` package. Gone from the model: every
      `FancyString`, `car.Emoji`, `stage.Distance.Emoji`, `drivetrain.Emoji`,
      `weather.Emoji`, `location.Flag`, `location.WeatherStrings` and the
      `challenge` package's emoji/`RandomFancyString` vars plus the unused
      `Config.String`. Plain-text `String()` / `DetailedString()` / `Name()`
      stay - they're identifiers, not decoration. The vestigial `\x1f`
      `EmojiDelimiter` (nothing parsed it after §1) is dropped from the output
      (`internal/bot/render/`, `internal/model/*`, the handlers, `debug`)
- [x] Replace `config`'s `init()` and `os.Exit(1)` with an explicit
      `config.Load() (Config, error)` called from `main`. Importing the package
      no longer has a side effect; a missing `.env` is a returned error, not a
      process exit. Package-level mutable `var`s and the dead `NOTSET` defaults
      are gone
      (`internal/config/config.go`, `cmd/main.go`)
- [x] Make command registration and cleanup return errors instead of calling
      `os.Exit(1)` from library code. `createCommands` / `cleanupGuildCommands` /
      `cleanupGlobalCommands` are now unexported, return errors (cleanup joins
      them), and `bot.New` propagates the failure. Dead `cmdIDs` map removed;
      `bot.cfg` is now populated (see §5)
      (`internal/bot/commands.go`, `internal/bot/bot.go`, `cmd/main.go`)
- [x] Move the store port out of the misnamed `internal/model` package (which
      held only the `Store` interface) to `internal/store/port` as `port.Store`.
      `internal/model/` is now just a namespace directory for the `challenge`,
      `car`, `stage` ... packages. Reducing the count of types named `Model` and
      the `challengeModel` alias is left for later
      (`internal/store/port/port.go`)
- [x] Thread `context.Context` through the store methods and drive shutdown
      from a cancelled context. `main` uses `signal.NotifyContext`; the context
      flows `bot` -> interaction handlers -> `port.Store`. Both stores refuse a
      cancelled context at entry (bbolt has no mid-transaction cancellation).
      Passing `discordgo.WithContext` into the Discord calls is deferred - it
      needs the `discord.Session` mock assertions reworked
      (`internal/store/port/port.go`, `internal/store/...`, `internal/bot/...`,
      `cmd/main.go`)

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
- [x] `bot.New` never sets the `cfg` field, so `Shutdown` and the
      `!cars` / `!stages` test-server gate run against a zero `Config`. Fixed
      alongside §3 (`bot.New` now sets `cfg`; cleanup no longer needs it)
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
