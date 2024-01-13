package query

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"

	"go.felesatra.moe/animanager/internal/sqlc"
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
	ctx := context.Background()
	p := sqlc.InsertWatchingParams{
		Aid:    nullint64(w.AID),
		Regexp: w.Regexp,
		Offset: int64(w.Offset),
	}
	if err := sqlc.New(db).InsertWatching(ctx, p); err != nil {
		return fmt.Errorf("insert watching %d: %w", w.AID, err)
	}
	return nil
}

// GetWatching gets the watching entry for an anime from the
// database.
func GetWatching(db sqlc.DBTX, aid AID) (Watching, error) {
	ctx := context.Background()
	w, err := sqlc.New(db).GetWatching(ctx, nullint64(aid))
	if err != nil {
		return Watching{}, fmt.Errorf("GetWatching %d: %w", aid, err)
	}
	return convertWatching(w), nil
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

func convertWatching(w sqlc.Watching) Watching {
	return Watching{
		AID:    AID(w.Aid.Int64),
		Regexp: w.Regexp,
		Offset: int(w.Offset),
	}
}
