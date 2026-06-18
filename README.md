# dice-collector

CLI tool for recording craps dice sessions to a SQLite database. (Originally a Python 2 / MySQL tool; see [Legacy Python version](#legacy-python-version) below.)

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

SQLite. Tables: `shooter`, `location`, `session`, `turn`, `throw`, `result` — same shape as the original MySQL schema. `turn.id` and `throw.sequence` are scoped per `(session, shooter)` / `(session, shooter, turn)`, matching the original MyISAM grouped `AUTO_INCREMENT` behavior.

The schema and seed data (original locations, shooters, result codes, and all historical sessions/turns/throws from `PyTom/Data/dice.sql.gz`) are embedded in the binary (`internal/dice/schema.sql`, `internal/dice/seed.sql`) and applied automatically to a fresh database file on first run.

By default the database file is `dice.db` in the working directory; override with `DICE_DB_PATH`.

## Structure

```
main.go            # entry point / CLI loop
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
go run .
# or
go build -o dice-collector . && ./dice-collector
```

## Legacy Python version

The original Python 2 / MySQL implementation is still present under `PyTom/` and `crunchdicesesh.py` for reference. It requires Python 2, `MySQLdb`, and a MySQL server seeded from `PyTom/Data/dice.sql.gz`; it is no longer maintained.
