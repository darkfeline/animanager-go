package query

import (
	"database/sql"

	"github.com/pkg/errors"
)

type Watching struct {
	AID    int
	Regexp string
}

// InsertWatching inserts or updates a watching entry into the database.
func InsertWatching(db *sql.DB, aid int, regexp string) error {
	t, err := db.Begin()
	if err != nil {
		return err
	}
	defer t.Rollback()
	_, err = t.Exec(`
INSERT INTO watching (aid, regexp) VALUES (?, ?)
ON CONFLICT (aid) DO UPDATE SET regexp=? WHERE aid=?`,
		aid, regexp,
		regexp, aid,
	)
	if err != nil {
		return errors.Wrapf(err, "failed to insert watching %d", aid)
	}
	return t.Commit()
}

// GetWatching gets the watching entry for an anime from the
// database. ErrMissing is returned if the anime doesn't exist.
func GetWatching(db *sql.DB, aid int) (Watching, error) {
	var w Watching
	r, err := db.Query(`SELECT aid, regexp FROM watching WHERE aid=?`, aid)
	if err != nil {
		return w, errors.Wrap(err, "failed to query watching")
	}
	defer r.Close()
	if !r.Next() {
		if r.Err() != nil {
			return w, r.Err()
		}
		return w, ErrMissing
	}
	if err := r.Scan(&w.AID, &w.Regexp); err != nil {
		return w, errors.Wrap(err, "failed to scan episode")
	}
	return w, nil
}

// GetAllWatching gets all watching entries.
func GetAllWatching(db *sql.DB) ([]Watching, error) {
	r, err := db.Query(`SELECT aid, regexp FROM watching`)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var result []Watching
	for r.Next() {
		var w Watching
		if err := r.Scan(&w.AID, &w.Regexp); err != nil {
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
// database.  This function returns ErrMissing if no such entry
// exists.
func DeleteWatching(db *sql.DB, aid int) error {
	t, err := db.Begin()
	if err != nil {
		return err
	}
	defer t.Rollback()
	r, err := t.Exec(`DELETE FROM watching WHERE aid=?`, aid)
	if err != nil {
		return err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrMissing
	}
	return t.Commit()
}
