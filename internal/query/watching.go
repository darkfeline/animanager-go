package query

import (
	"context"
	"fmt"
	"regexp"

	"go.felesatra.moe/animanager/internal/sqlc"
)

type Watching struct {
	_table struct{} `sql:"watching"`
	AID    sqlc.AID `sql:"aid"`
	Regexp string   `sql:"regexp"`
	Offset int      `sql:"offset"`
}

// InsertWatching inserts or updates a watching entry into the database.
func InsertWatching(db sqlc.DBTX, w Watching) error {
	if _, err := regexp.Compile(w.Regexp); err != nil {
		return fmt.Errorf("insert watching %d: %w", w.AID, err)
	}
	ctx := context.Background()
	p := sqlc.InsertWatchingParams{
		Aid:    w.AID,
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
func GetWatching(db sqlc.DBTX, aid sqlc.AID) (Watching, error) {
	ctx := context.Background()
	w, err := sqlc.New(db).GetWatching(ctx, aid)
	if err != nil {
		return Watching{}, fmt.Errorf("GetWatching %d: %w", aid, err)
	}
	return convertWatching(w), nil
}

// GetWatchingCount returns the number of watching rows.
func GetWatchingCount(db sqlc.DBTX) (int, error) {
	ctx := context.Background()
	r, err := sqlc.New(db).GetWatchingCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("GetWatchingCount: %s", err)
	}
	return int(r), nil
}

// GetAllWatching gets all watching entries.
func GetAllWatching(db sqlc.DBTX) ([]Watching, error) {
	ctx := context.Background()
	w, err := sqlc.New(db).GetAllWatching(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetAllWatching: %s", err)
	}
	return smap(w, convertWatching), nil
}

// DeleteWatching deletes the watching entry for an anime from the
// database.
func DeleteWatching(db sqlc.DBTX, aid sqlc.AID) error {
	ctx := context.Background()
	if err := sqlc.New(db).DeleteWatching(ctx, aid); err != nil {
		return fmt.Errorf("DeleteWatching %d: %s", aid, err)
	}
	return nil
}

func convertWatching(w sqlc.Watching) Watching {
	return Watching{
		AID:    sqlc.AID(w.Aid),
		Regexp: w.Regexp,
		Offset: int(w.Offset),
	}
}
