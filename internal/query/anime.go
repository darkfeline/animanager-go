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
	"database/sql"
	"fmt"
	"log/slog"

	"go.felesatra.moe/anidb"

	"go.felesatra.moe/animanager/internal/date"
	"go.felesatra.moe/animanager/internal/sqlc"
)

// Anime values correspond to rows in the anime table.
type Anime struct {
	_table       struct{}  `sql:"anime"`
	AID          AID       `sql:"aid"`
	Title        string    `sql:"title"`
	Type         AnimeType `sql:"type"`
	EpisodeCount int       `sql:"episodecount"`
	// The following fields are nullable.  In most cases, use the
	// getter methods instead.
	NullStartDate sql.NullInt64 `sql:"startdate"`
	NullEndDate   sql.NullInt64 `sql:"enddate"`
}

// StartDate returns the NullStartDate field as a Date.  If the
// field is invalid, returns date.Zero.
func (a Anime) StartDate() date.Date {
	v := a.NullStartDate
	if v.Valid {
		return date.Date(v.Int64)
	} else {
		return date.Zero
	}
}

// EndDate returns the NullEndDate field as a Date.  If the
// field is invalid, returns date.Zero.
func (a Anime) EndDate() date.Date {
	v := a.NullEndDate
	if v.Valid {
		return date.Date(v.Int64)
	} else {
		return date.Zero
	}
}

type AnimeType string

// GetAnimeCount returns the number of anime rows.
func GetAnimeCount(db sqlc.DBTX) (int, error) {
	ctx := context.Background()
	n, err := sqlc.New(db).GetAnimeCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("GetAnimeCount: %s", err)
	}
	return int(n), nil
}

// GetAIDs returns all AIDs.
func GetAIDs(db sqlc.DBTX) ([]AID, error) {
	ctx := context.Background()
	aids, err := sqlc.New(db).GetAIDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetAIDs: %s", err)
	}
	return convertMany(aids, func(v sql.NullInt64) AID { return AID(v.Int64) }), nil
}

// GetAnime gets the anime from the database.
func GetAnime(db sqlc.DBTX, aid AID) (*Anime, error) {
	ctx := context.Background()
	a, err := sqlc.New(db).GetAnime(ctx, nullint64(aid))
	if err != nil {
		return nil, fmt.Errorf("GetAnime %d: %s", aid, err)
	}
	a2 := convertAnime(a)
	return &a2, nil
}

// GetAllAnime returns all anime.
func GetAllAnime(db sqlc.DBTX) ([]Anime, error) {
	ctx := context.Background()
	a, err := sqlc.New(db).GetAllAnime(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetAllAnime: %s", err)
	}
	return convertMany(a, convertAnime), nil
}

// InsertAnime inserts or updates an anime into the database.
func InsertAnime(db *sql.DB, a *anidb.Anime) error {
	t, err := db.Begin()
	if err != nil {
		return err
	}
	defer t.Rollback()
	var startDate interface{}
	var endDate interface{}
	startDate, err = date.Parse(a.StartDate)
	if err != nil {
		startDate = nil
	}
	endDate, err = date.Parse(a.EndDate)
	if err != nil {
		endDate = nil
	}
	title := mainTitle(a.Titles)
	_, err = t.Exec(`
INSERT INTO anime (aid, title, type, episodecount, startdate, enddate)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT (aid) DO UPDATE SET
title=excluded.title, type=excluded.type, episodecount=excluded.episodecount,
startdate=excluded.startdate, enddate=excluded.enddate
WHERE aid=excluded.aid`,
		a.AID, title, a.Type, a.EpisodeCount, startDate, endDate,
	)
	if err != nil {
		return fmt.Errorf("failed to insert anime %d: %w", a.AID, err)
	}
	em, err := GetEpisodesMap(t, AID(a.AID))
	if err != nil {
		return fmt.Errorf("failed to insert anime %d: %w", a.AID, err)
	}
	for _, e := range a.Episodes {
		if err := insertEpisode(t, AID(a.AID), e); err != nil {
			return fmt.Errorf("failed to insert episode %s for anime %d: %w",
				e.EpNo, a.AID, err)
		}
		delete(em, EID(e.EID))
	}
	for eid := range em {
		if err := DeleteEpisode(t, eid); err != nil {
			return fmt.Errorf("failed to insert anime %d: %w", a.AID, err)
		}
	}
	return t.Commit()
}

// mainTitle returns the main title from a slice of titles.
func mainTitle(ts []anidb.Title) string {
	for _, t := range ts {
		if t.Type == "main" {
			return t.Name
		}
	}
	return ts[0].Name
}

func insertEpisode(t *sql.Tx, aid AID, e anidb.Episode) error {
	title := mainEpTitle(e.Titles)
	typ, num := parseEpNo(e.EpNo)
	slog.Debug("insert episode",
		"eid", e.EID,
		"aid", aid,
		"type", typ,
		"number", num,
		"title", title,
		"length", e.Length,
	)
	if typ == EpUnknown {
		return fmt.Errorf("failed to insert anime %d: invalid epno %s", aid, e.EpNo)
	}
	_, err := t.Exec(`
INSERT INTO episode (eid, aid, type, number, title, length)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT (eid) DO UPDATE SET
aid=excluded.aid, type=excluded.type, number=excluded.number,
title=excluded.title, length=excluded.length
WHERE eid=excluded.eid`,
		e.EID, aid, typ, num, title, e.Length,
	)
	if err != nil {
		return err
	}
	return nil
}

// mainEpTitle returns the title to use from a slice of episode titles.
func mainEpTitle(ts []anidb.EpTitle) string {
	for _, t := range ts {
		if t.Lang == "ja" {
			return t.Title
		}
	}
	return ts[0].Title
}

func convertAnime(a sqlc.Anime) Anime {
	return Anime{
		AID:           AID(a.Aid.Int64),
		Title:         a.Title,
		Type:          AnimeType(a.Type),
		EpisodeCount:  int(a.Episodecount),
		NullStartDate: a.Startdate,
		NullEndDate:   a.Enddate,
	}
}
