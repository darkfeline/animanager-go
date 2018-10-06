package query

import (
	"database/sql"

	"github.com/pkg/errors"
)

type Watching struct {
	AID    int
	Regexp string
	Offset int
}

// InsertWatching inserts or updates a watching entry into the database.
func InsertWatching(db *sql.DB, w Watching) error {
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
		return errors.Wrapf(err, "failed to insert watching %d", w.AID)
	}
	return t.Commit()
}

// GetWatching gets the watching entry for an anime from the
// database.  This function's error implements Error.
func GetWatching(db *sql.DB, aid int) (Watching, error) {
	r := db.QueryRow(`SELECT aid, regexp, offset FROM watching WHERE aid=?`, aid)
	var w Watching
	if err := r.Scan(&w.AID, &w.Regexp, &w.Offset); err != nil {
		return w, errors.Wrapf(err, "GetWatching %d", aid)
	}
	return w, nil
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
func DeleteWatching(db *sql.DB, aid int) error {
	_, err := db.Exec(`DELETE FROM watching WHERE aid=?`, aid)
	return err
}
