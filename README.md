# dice-collector

CLI tool for recording craps dice sessions to a SQLite database.

## What it does

Walks through an interactive session entry flow:

1. Pick a casino/location from the DB
2. Loop over shooter turns — choose shooter, enter all throws as a compact string
3. Each throw is parsed, craps logic is applied (comeout, point, 7-out, etc.), and the result is written to the DB
4. Prints a session summary on exit

## Throw input encoding

Throws are entered as a single string of characters, one char per roll:

| Char | Value |
|------|-------|
| `2`–`9` | Face value |
| `T` | 10 |
| `E` | 11 (yo-leven) |
| `B` | 12 (boxcars) |
| `O` | Off table |
| `M` | Misc / invalid roll |

Example: `456T74` → 4, 5, 6, 10, 7, 4

## Result codes (stored in `throw.result`)

| Code | Meaning |
|------|---------|
| 1 | No consequence |
| 2 | Comeout win (7 or 11) |
| 3 | Comeout loss (2, 3, or 12) |
| 4 | Point established |
| 5 | Point made (win) |
| 6 | Seven-out (loss) |
| 7 | Off table |
| 8 | Misc / invalid |

## Database

SQLite. Tables: `shooter`, `location`, `session`, `turn`, `throw`, `result`. `turn.id` is scoped per `(session, shooter)` and `throw.sequence` per `(session, shooter, turn)` — each counts up independently within its group rather than globally.

Schema and seed data (locations, shooters, result codes, historical sessions/turns/throws) are embedded in the binary (`internal/dice/schema.sql`, `internal/dice/seed.sql`) and applied automatically on first run. `seed.sql` only runs against a brand-new database file, so it has no effect on a `dice.db` you've already been using — shooters and locations on an existing database are managed via ad hoc queries against the SQLite file.

By default the database file is `dice.db` in the working directory; override with `DICE_DB_PATH`.

## Structure

```
cmd/dice-collector/
  main.go          # entry point / CLI loop
internal/dice/
  logic.go         # craps rules / throw processing
  store.go         # DB read/write helpers, schema bootstrap
  schema.sql       # table definitions (embedded)
  seed.sql         # lookup + historical data (embedded)
```

## Requirements

- Go 1.25+
- `modernc.org/sqlite` (pure-Go SQLite driver, no CGO/C toolchain needed) — fetched automatically via `go build`/`go run`

## Running

```
go run ./cmd/dice-collector
# or
go build -o dice-collector ./cmd/dice-collector && ./dice-collector
```

## Installing a prebuilt binary

Every push to `main` rebuilds cross-platform binaries (linux/darwin, amd64/arm64, plus windows/amd64) and publishes them to the [`latest` release](https://github.com/tphummel/dice-collector/releases/tag/latest), alongside a `checksums.txt`. Download the one matching your platform and run it directly — no Go toolchain required.

## CI/CD

`.github/workflows/ci.yml` runs `go test` and `go vet` on every push/PR to `main`, then (on `main` pushes only) builds and publishes the binaries described above.
