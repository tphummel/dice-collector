package dice

import (
	"path/filepath"
	"testing"
)

func TestOpen_SeedsOnceAndIsIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")

	store, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	shooters, err := store.Shooters()
	if err != nil {
		t.Fatalf("Shooters: %v", err)
	}
	if len(shooters) != 4 {
		t.Fatalf("got %d seeded shooters, want 4", len(shooters))
	}
	store.Close()

	// Reopening an existing database must not re-run schema/seed SQL.
	store2, err := Open(path)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer store2.Close()

	shooters2, err := store2.Shooters()
	if err != nil {
		t.Fatalf("Shooters after reopen: %v", err)
	}
	if len(shooters2) != 4 {
		t.Fatalf("got %d shooters after reopen, want 4 (seed data was duplicated)", len(shooters2))
	}
}

func TestInsertTurn_ScopedPerSessionAndShooter(t *testing.T) {
	store, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer store.Close()

	sessionID, err := store.InsertSession(1)
	if err != nil {
		t.Fatalf("InsertSession: %v", err)
	}

	turn1, err := store.InsertTurn(sessionID, 1)
	if err != nil {
		t.Fatalf("InsertTurn: %v", err)
	}
	if turn1 != 1 {
		t.Errorf("first turn for shooter 1 = %d, want 1", turn1)
	}

	turn2, err := store.InsertTurn(sessionID, 1)
	if err != nil {
		t.Fatalf("InsertTurn: %v", err)
	}
	if turn2 != 2 {
		t.Errorf("second turn for shooter 1 = %d, want 2", turn2)
	}

	turn3, err := store.InsertTurn(sessionID, 2)
	if err != nil {
		t.Fatalf("InsertTurn: %v", err)
	}
	if turn3 != 1 {
		t.Errorf("first turn for shooter 2 = %d, want 1 (independent of shooter 1's count)", turn3)
	}
}

func TestInsertThrow_ScopedPerSessionShooterTurn(t *testing.T) {
	store, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer store.Close()

	sessionID, err := store.InsertSession(1)
	if err != nil {
		t.Fatalf("InsertSession: %v", err)
	}
	turnA, err := store.InsertTurn(sessionID, 1)
	if err != nil {
		t.Fatalf("InsertTurn: %v", err)
	}
	turnB, err := store.InsertTurn(sessionID, 1)
	if err != nil {
		t.Fatalf("InsertTurn: %v", err)
	}

	if err := store.InsertThrow(sessionID, 1, turnA, 5, ResultNoConsequence, false); err != nil {
		t.Fatalf("InsertThrow: %v", err)
	}
	if err := store.InsertThrow(sessionID, 1, turnA, 6, ResultNoConsequence, false); err != nil {
		t.Fatalf("InsertThrow: %v", err)
	}
	if err := store.InsertThrow(sessionID, 1, turnB, 7, ResultSevenOut, false); err != nil {
		t.Fatalf("InsertThrow: %v", err)
	}

	var seqA, countA, seqB int64
	if err := store.db.QueryRow(`SELECT MAX(sequence), COUNT(*) FROM throw WHERE sessionid = ? AND shooterid = 1 AND turnid = ?`, sessionID, turnA).Scan(&seqA, &countA); err != nil {
		t.Fatalf("query turnA throws: %v", err)
	}
	if seqA != 2 || countA != 2 {
		t.Errorf("turnA: got max sequence %d over %d rows, want 2 over 2 rows", seqA, countA)
	}

	if err := store.db.QueryRow(`SELECT MAX(sequence) FROM throw WHERE sessionid = ? AND shooterid = 1 AND turnid = ?`, sessionID, turnB).Scan(&seqB); err != nil {
		t.Fatalf("query turnB throws: %v", err)
	}
	if seqB != 1 {
		t.Errorf("turnB sequence = %d, want 1 (independent of turnA's count)", seqB)
	}
}
