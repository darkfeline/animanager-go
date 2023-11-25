package query

import (
	"database/sql"
	"fmt"
	"regexp"
)

type Watching struct {
	_table struct{} `sql:"watching"`
	AID    AID      `sql:"aid"`
	Regexp string   `sql:"regexp"`
	Offset int      `sql:"offset"`
}

// InsertWatching inserts or updates a watching entry into the database.
func InsertWatching(db *sql.DB, w Watching) error {
	if _, err := regexp.Compile(w.Regexp); err != nil {
		return fmt.Errorf("insert watching %d: %w", w.AID, err)
	}
	t, err := db.Begin()
	if err != nil {
		return err
	}
	defer t.Rollback()
	_, err = t.Exec(`
INSERT INTO watching (aid, regexp, offset) VALUES (?, ?, ?)
ON CONFLICT (aid) DO UPDATE SET regexp=?, offset=? WHERE aid=?`,
		w.AID, w.Regexp, w.Offset,
		w.Regexp, w.Offset, w.AID,
	)
	if err != nil {
		return fmt.Errorf("insert watching %d: %w", w.AID, err)
	}
	return t.Commit()
}

// GetWatching gets the watching entry for an anime from the
// database.
func GetWatching(db *sql.DB, aid AID) (Watching, error) {
	r := db.QueryRow(`SELECT aid, regexp, offset FROM watching WHERE aid=?`, aid)
	var w Watching
	if err := r.Scan(&w.AID, &w.Regexp, &w.Offset); err != nil {
		return w, fmt.Errorf("GetWatching %d: %w", aid, err)
	}
	return w, nil
}

// GetWatchingCount returns the number of watching rows.
func GetWatchingCount(db *sql.DB) (int, error) {
	r := db.QueryRow(`SELECT COUNT(*) FROM watching`)
	var n int
	err := r.Scan(&n)
	return n, err
}

// GetAllWatching gets all watching entries.
func GetAllWatching(db *sql.DB) ([]Watching, error) {
	r, err := db.Query(`SELECT aid, regexp, offset FROM watching`)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var result []Watching
	for r.Next() {
		var w Watching
		if err := r.Scan(&w.AID, &w.Regexp, &w.Offset); err != nil {
			return nil, err
		}
		result = append(result, w)
	}
	if r.Err() != nil {
		return nil, r.Err()
	}
	return result, nil
}

// DeleteWatching deletes the watching entry for an anime from the
// database.
func DeleteWatching(db *sql.DB, aid AID) error {
	_, err := db.Exec(`DELETE FROM watching WHERE aid=?`, aid)
	return err
}
