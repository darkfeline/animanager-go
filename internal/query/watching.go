package query

import (
	"context"
	"fmt"
	"regexp"

	"go.felesatra.moe/animanager/internal/sqlc"
)

type Watching = sqlc.Watching

// InsertWatching inserts or updates a watching entry into the database.
func InsertWatching(db sqlc.DBTX, w Watching) error {
	if _, err := regexp.Compile(w.Regexp); err != nil {
		return fmt.Errorf("insert watching %d: %w", w.Aid, err)
	}
	ctx := context.Background()
	p := sqlc.InsertWatchingParams{
		Aid:    w.Aid,
		Regexp: w.Regexp,
		Offset: w.Offset,
	}
	if err := sqlc.New(db).InsertWatching(ctx, p); err != nil {
		return fmt.Errorf("insert watching %d: %w", w.Aid, err)
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
	return w, nil
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
	return w, nil
}

// GetFinishedWatchingAIDs returns the AIDs for finished anime with
// watching entries.
func GetFinishedWatchingAIDs(db sqlc.DBTX) ([]sqlc.AID, error) {
	watching, err := GetAllWatching(db)
	if err != nil {
		return nil, fmt.Errorf("GetFinishedWatchingAIDs: %s", err)
	}
	watchingMap := make(map[sqlc.AID]bool)
	for _, w := range watching {
		watchingMap[w.Aid] = true
	}

	finished, err := GetFinishedAnime(db)
	if err != nil {
		return nil, fmt.Errorf("GetFinishedWatchingAIDs: %s", err)
	}
	var aids []sqlc.AID
	for _, a := range finished {
		if watchingMap[a.AID] {
			aids = append(aids, a.AID)
		}
	}
	return aids, nil
}
