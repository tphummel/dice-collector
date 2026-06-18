package dice

import (
	"fmt"
	"strconv"
)

// Result codes stored in throw.result.
const (
	ResultNoConsequence = 1
	ResultComeoutWin    = 2
	ResultComeoutLoss   = 3
	ResultPoint         = 4
	ResultPointWin      = 5
	ResultSevenOut      = 6
	ResultOffTable      = 7
	ResultMisc          = 8
)

// Throw is one parsed roll: its face value, the craps result it produced, and
// whether it occurred during a comeout.
type Throw struct {
	Value       int
	Result      int
	PropComeout bool
}

// ProcessTurnThrows applies craps rules to a string of throw characters,
// one rune per roll (see README for the character encoding).
func ProcessTurnThrows(throws string) ([]Throw, error) {
	isComeout := true
	point := 0

	out := make([]Throw, 0, len(throws))
	for _, ch := range throws {
		value, result, err := decodeThrowChar(ch)
		if err != nil {
			return nil, err
		}

		t := Throw{Value: value, Result: result, PropComeout: isComeout}

		if value != 1 {
			if isComeout {
				switch value {
				case 7, 11:
					t.Result = ResultComeoutWin
				case 2, 3, 12:
					t.Result = ResultComeoutLoss
				default:
					point = value
					t.Result = ResultPoint
					isComeout = false
				}
			} else {
				switch {
				case value == point:
					t.Result = ResultPointWin
					isComeout = true
				case value == 7:
					t.Result = ResultSevenOut
					isComeout = true
					point = 0
				default:
					t.Result = ResultNoConsequence
				}
			}
		}

		out = append(out, t)
	}

	return out, nil
}

func decodeThrowChar(ch rune) (value int, result int, err error) {
	switch ch {
	case 'T':
		return 10, 0, nil
	case 'E':
		return 11, 0, nil
	case 'B':
		return 12, 0, nil
	case 'O':
		return 1, ResultOffTable, nil
	case 'M':
		return 1, ResultMisc, nil
	default:
		n, err := strconv.Atoi(string(ch))
		if err != nil {
			return 0, 0, fmt.Errorf("invalid throw character %q", ch)
		}
		return n, 0, nil
	}
}
