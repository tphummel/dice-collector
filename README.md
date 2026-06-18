# dice-collector

CLI tool for recording craps dice sessions to a MySQL database.

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

MySQL. Schema is in `PyTom/Data/dice.sql.gz`. Tables: `shooter`, `location`, `session`, `turn`, `throw`, `result`.

Connection is configured in `PyTom/Data/mysql.py`.

## Structure

```
crunchdicesesh.py       # entry point
PyTom/
  Dice/
    logic.py            # craps rules / throw processing
    data.py             # DB read/write helpers
    sql/                # MySQL stored functions
  Data/
    mysql.py            # DB connection
  Tcrypt/Utilities/
    io.py               # file I/O utility (unused)
```

## Requirements

- Python 3
- `mysqlclient` (`pip install mysqlclient`)
