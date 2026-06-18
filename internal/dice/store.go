package dice

import (
	"database/sql"
	_ "embed"
	"fmt"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaSQL string

//go:embed seed.sql
var seedSQL string

type Shooter struct {
	ID   int64
	Name string
}

type Location struct {
	ID   int64
	Name string
}

type Store struct {
	db *sql.DB
}

// Open connects to the SQLite database at path, creating and seeding the
// schema on first use.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		db.Close()
		return nil, err
	}

	var tableCount int
	err = db.QueryRow(`SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name = 'session'`).Scan(&tableCount)
	if err != nil {
		db.Close()
		return nil, err
	}

	if tableCount == 0 {
		if _, err := db.Exec(schemaSQL); err != nil {
			db.Close()
			return nil, fmt.Errorf("applying schema: %w", err)
		}
		if _, err := db.Exec(seedSQL); err != nil {
			db.Close()
			return nil, fmt.Errorf("applying seed data: %w", err)
		}
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Shooters() ([]Shooter, error) {
	rows, err := s.db.Query(`SELECT id, name FROM shooter ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shooters []Shooter
	for rows.Next() {
		var sh Shooter
		if err := rows.Scan(&sh.ID, &sh.Name); err != nil {
			return nil, err
		}
		shooters = append(shooters, sh)
	}
	return shooters, rows.Err()
}

func (s *Store) Locations() ([]Location, error) {
	rows, err := s.db.Query(`SELECT id, name FROM location ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []Location
	for rows.Next() {
		var loc Location
		if err := rows.Scan(&loc.ID, &loc.Name); err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}
	return locations, rows.Err()
}

func (s *Store) InsertSession(locationID int64) (int64, error) {
	res, err := s.db.Exec(`INSERT INTO session (locationid) VALUES (?)`, locationID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// InsertTurn records a new turn for shooterID within sessionID, with id
// scoped per (session, shooter) — see README's Database section.
func (s *Store) InsertTurn(sessionID, shooterID int64) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var nextID int64
	err = tx.QueryRow(
		`SELECT COALESCE(MAX(id), 0) + 1 FROM turn WHERE sessionid = ? AND shooterid = ?`,
		sessionID, shooterID,
	).Scan(&nextID)
	if err != nil {
		return 0, err
	}

	if _, err := tx.Exec(
		`INSERT INTO turn (sessionid, shooterid, id) VALUES (?, ?, ?)`,
		sessionID, shooterID, nextID,
	); err != nil {
		return 0, err
	}

	return nextID, tx.Commit()
}

// InsertThrow records one roll, with sequence scoped per (session, shooter, turn).
func (s *Store) InsertThrow(sessionID, shooterID, turnID int64, value, result int, propComeout bool) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var nextSeq int64
	err = tx.QueryRow(
		`SELECT COALESCE(MAX(sequence), 0) + 1 FROM throw WHERE sessionid = ? AND shooterid = ? AND turnid = ?`,
		sessionID, shooterID, turnID,
	).Scan(&nextSeq)
	if err != nil {
		return err
	}

	comeout := "N"
	if propComeout {
		comeout = "Y"
	}

	if _, err := tx.Exec(
		`INSERT INTO throw (sessionid, shooterid, turnid, sequence, value, result, event_prop_comeout) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		sessionID, shooterID, turnID, nextSeq, value, result, comeout,
	); err != nil {
		return err
	}

	return tx.Commit()
}
