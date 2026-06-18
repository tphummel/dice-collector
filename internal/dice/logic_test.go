package dice

import (
	"reflect"
	"testing"
)

func TestDecodeThrowChar(t *testing.T) {
	cases := []struct {
		ch         rune
		wantValue  int
		wantResult int
		wantErr    bool
	}{
		{'T', 10, 0, false},
		{'E', 11, 0, false},
		{'B', 12, 0, false},
		{'O', 1, ResultOffTable, false},
		{'M', 1, ResultMisc, false},
		{'5', 5, 0, false},
		{'x', 0, 0, true},
	}

	for _, c := range cases {
		value, result, err := decodeThrowChar(c.ch)
		if c.wantErr {
			if err == nil {
				t.Errorf("decodeThrowChar(%q): want error, got none", c.ch)
			}
			continue
		}
		if err != nil {
			t.Errorf("decodeThrowChar(%q): unexpected error: %v", c.ch, err)
		}
		if value != c.wantValue || result != c.wantResult {
			t.Errorf("decodeThrowChar(%q) = (%d, %d), want (%d, %d)", c.ch, value, result, c.wantValue, c.wantResult)
		}
	}
}

func TestProcessTurnThrows(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []Throw
	}{
		{
			name:  "comeout win on 7",
			input: "7",
			want:  []Throw{{Value: 7, Result: ResultComeoutWin, PropComeout: true}},
		},
		{
			name:  "comeout win on 11",
			input: "E",
			want:  []Throw{{Value: 11, Result: ResultComeoutWin, PropComeout: true}},
		},
		{
			name:  "comeout loss on 2",
			input: "2",
			want:  []Throw{{Value: 2, Result: ResultComeoutLoss, PropComeout: true}},
		},
		{
			name:  "comeout loss on 12",
			input: "B",
			want:  []Throw{{Value: 12, Result: ResultComeoutLoss, PropComeout: true}},
		},
		{
			name:  "point established then made",
			input: "55",
			want: []Throw{
				{Value: 5, Result: ResultPoint, PropComeout: true},
				{Value: 5, Result: ResultPointWin, PropComeout: false},
			},
		},
		{
			name:  "point established, no-consequence rolls, seven-out, new point",
			input: "456T74",
			want: []Throw{
				{Value: 4, Result: ResultPoint, PropComeout: true},
				{Value: 5, Result: ResultNoConsequence, PropComeout: false},
				{Value: 6, Result: ResultNoConsequence, PropComeout: false},
				{Value: 10, Result: ResultNoConsequence, PropComeout: false},
				{Value: 7, Result: ResultSevenOut, PropComeout: false},
				{Value: 4, Result: ResultPoint, PropComeout: true},
			},
		},
		{
			name:  "off-table rolls don't disturb comeout state",
			input: "7O7",
			want: []Throw{
				{Value: 7, Result: ResultComeoutWin, PropComeout: true},
				{Value: 1, Result: ResultOffTable, PropComeout: true},
				{Value: 7, Result: ResultComeoutWin, PropComeout: true},
			},
		},
		{
			name:  "misc rolls don't disturb an established point",
			input: "5M5",
			want: []Throw{
				{Value: 5, Result: ResultPoint, PropComeout: true},
				{Value: 1, Result: ResultMisc, PropComeout: false},
				{Value: 5, Result: ResultPointWin, PropComeout: false},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := ProcessTurnThrows(c.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("ProcessTurnThrows(%q) = %+v, want %+v", c.input, got, c.want)
			}
		})
	}
}

func TestProcessTurnThrows_InvalidCharacter(t *testing.T) {
	if _, err := ProcessTurnThrows("5x9"); err == nil {
		t.Fatal("want error for invalid throw character, got none")
	}
}
