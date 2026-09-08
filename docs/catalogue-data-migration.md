# Catalogue Data Migration (future work)

## Where we are

The game catalogue — which locations, stages, car classes, cars and drivetrains
exist for each game — lives in package-level Go maps under `internal/model/`:

- `location.byGame` — locations per game
- `stage.byLocation` — stages per location
- `class.byGame`, `class.byGameDrivetrain` — classes per game, and per drivetrain
- `car.byGameClass` — cars per class, per game
- `drivetrain.byGame` — drivetrains per game

This replaced an earlier `switch game { listDR2() / listWRC() }` twin-function
pattern. Lookups are now a single map access instead of a dispatch, and adding a
game is a data edit rather than a new function in every package.

## Why move further

The data is still Go source. That means:

- only a developer with a toolchain can change a stage name or add a car
- a data change needs a recompile and redeploy
- the data and the code that reads it are in the same files, so the packages
  are large and dominated by literals

The next step is to move the catalogue into embedded data files
(`//go:embed` JSON or TOML) loaded through a small registry at startup, so the
data can be edited on its own and validated on load.

## What that involves

- a schema per entity (location, stage, class, car, drivetrain), keyed by game
- data files under something like `internal/model/data/`
- a loader that reads them once, validates referential integrity (every car's
  class exists, every stage's location exists, every class's drivetrain exists)
  and fails fast on a bad file
- the `List` / `AtLocation` / `InClass` / `WithDrivetrain` functions read from
  the loaded registry instead of the package maps
- the enum types (`location.Model`, `class.Model`, ...) either stay as the
  identifiers used in the data files, or are replaced by string IDs

## Trade-offs to weigh first

- compile-time exhaustiveness goes away; a typo in a data file is a runtime
  error caught only by the loader's validation, so that validation has to be
  thorough and well tested
- `String()`, `Emoji()`, `Flag()` and similar presentation methods are
  genuinely code, not data — decide whether they move into the data files or
  stay as Go
- the randomiser and the challenge builder must keep working against whatever
  the registry exposes
