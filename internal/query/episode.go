// Copyright (C) 2018  Allen Li
//
// This file is part of Animanager.
//
// Animanager is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Animanager is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Animanager.  If not, see <http://www.gnu.org/licenses/>.

package query

import (
	"context"
	"fmt"

	"go.felesatra.moe/animanager/internal/sqlc"
)

type Episode struct {
	_table      struct{}    `sql:"episode"`
	EID         sqlc.EID    `sql:"eid"`
	AID         sqlc.AID    `sql:"aid"`
	Type        EpisodeType `sql:"type"`
	Number      int         `sql:"number"`
	Title       string      `sql:"title"`
	Length      int         `sql:"length"`
	UserWatched bool        `sql:"user_watched"`
}

func (e Episode) Key() EpisodeKey {
	return EpisodeKey{
		AID:    e.AID,
		Type:   e.Type,
		Number: e.Number,
	}
}

// EpisodeKey represents the unique key for an Episode.  This is
// separate from ID because SQLite treats numeric row IDs specially.
type EpisodeKey struct {
	AID    sqlc.AID
	Type   EpisodeType
	Number int
}

// GetEpisodeCount returns the number of episode rows.
func GetEpisodeCount(db sqlc.DBTX) (int, error) {
	ctx := context.Background()
	r, err := sqlc.New(db).GetEpisodeCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("GetEpisodeCount: %s", err)
	}
	return int(r), nil
}

// GetWatchedEpisodeCount returns the number of watched episodes.
func GetWatchedEpisodeCount(db sqlc.DBTX) (int, error) {
	ctx := context.Background()
	r, err := sqlc.New(db).GetWatchedEpisodeCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("GetWatchedEpisodeCount: %s", err)
	}
	return int(r), nil
}

// GetWatchedMinutes returns the number of minutes watched.
func GetWatchedMinutes(db sqlc.DBTX) (int, error) {
	ctx := context.Background()
	r, err := sqlc.New(db).GetWatchedMinutes(ctx)
	if err != nil {
		return 0, fmt.Errorf("GetWatchedMinutes: %s", err)
	}
	// BUG(kf): Watched minutes truncated to 32-bit, in case you
	// have watched more than 8000 years of anime.
	return int(r), err
}

// GetEpisode gets the episode from the database.
func GetEpisode(db sqlc.DBTX, eid sqlc.EID) (*Episode, error) {
	ctx := context.Background()
	e, err := sqlc.New(db).GetEpisode(ctx, eid)
	if err != nil {
		return nil, fmt.Errorf("GetEpisode %d: %s", eid, err)
	}
	e2 := convertEpisode(e)
	return &e2, nil
}

// DeleteEpisode deletes the episode from the database.
func DeleteEpisode(db sqlc.DBTX, eid sqlc.EID) error {
	ctx := context.Background()
	err := sqlc.New(db).DeleteEpisode(ctx, eid)
	if err != nil {
		return fmt.Errorf("DeleteEpisode %v: %w", eid, err)
	}
	return nil
}

// GetEpisodes gets the episodes for an anime from the database.
func GetEpisodes(db sqlc.DBTX, aid sqlc.AID) ([]Episode, error) {
	ctx := context.Background()
	e, err := sqlc.New(db).GetEpisodes(ctx, aid)
	if err != nil {
		return nil, fmt.Errorf("GetEpisodes %d: %s", aid, err)
	}
	e2 := smap(e, convertEpisode)
	return e2, nil
}

// GetEpisodesMap returns a map of the episodes for an anime.
func GetEpisodesMap(db sqlc.DBTX, aid sqlc.AID) (map[sqlc.EID]*Episode, error) {
	es, err := GetEpisodes(db, aid)
	if err != nil {
		return nil, fmt.Errorf("GetEpisodesMap %v: %w", aid, err)
	}
	m := make(map[sqlc.EID]*Episode, len(es))
	for i, e := range es {
		m[e.EID] = &es[i]
	}
	return m, nil
}

// GetAllEpisodes gets all episodes from the database.
func GetAllEpisodes(db sqlc.DBTX) ([]Episode, error) {
	ctx := context.Background()
	e, err := sqlc.New(db).GetAllEpisodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetAllEpisodes: %s", err)
	}
	e2 := smap(e, convertEpisode)
	return e2, nil
}

// UpdateEpisodeDone updates the episode's done status.
func UpdateEpisodeDone(db sqlc.DBTX, eid sqlc.EID, done bool) error {
	p := sqlc.UpdateEpisodeDoneParams{
		Eid:         eid,
		UserWatched: done,
	}

	ctx := context.Background()
	if err := sqlc.New(db).UpdateEpisodeDone(ctx, p); err != nil {
		return fmt.Errorf("UpdateEpisodeDone: %s", err)
	}
	return nil
}

func convertEpisode(e sqlc.Episode) Episode {
	e2 := Episode{
		EID:         e.Eid,
		AID:         e.Aid,
		Type:        e.Type,
		Number:      e.Number,
		Title:       e.Title,
		Length:      e.Length,
		UserWatched: e.UserWatched,
	}
	return e2
}

func smap[T, T2 any](v []T, f func(T) T2) []T2 {
	if len(v) == 0 {
		return nil
	}
	v2 := make([]T2, len(v))
	for i, v := range v {
		v2[i] = f(v)
	}
	return v2
}
