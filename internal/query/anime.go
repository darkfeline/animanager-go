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
	"database/sql"
	"fmt"

	"go.felesatra.moe/anidb"

	"go.felesatra.moe/animanager/internal/date"
)

// Anime values correspond to rows in the anime table.
type Anime struct {
	_table       struct{}  `anime`
	AID          int       `aid`
	Title        string    `title`
	Type         AnimeType `type`
	EpisodeCount int       `episodecount`
	// The following fields are nullable.  In most cases, use the
	// getter methods instead.
	NullStartDate sql.NullInt64 `startdate`
	NullEndDate   sql.NullInt64 `enddate`
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
func GetAnimeCount(db *sql.DB) (int, error) {
	r := db.QueryRow(`SELECT COUNT(*) FROM anime`)
	var n int
	err := r.Scan(&n)
	return n, err
}

// GetAIDs returns all AIDs.
func GetAIDs(db *sql.DB) ([]int, error) {
	r, err := db.Query(`SELECT aid FROM anime`)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var aids []int
	for r.Next() {
		var aid int
		if err := r.Scan(&aid); err != nil {
			return nil, err
		}
		aids = append(aids, aid)
	}
	if err := r.Err(); err != nil {
		return nil, err
	}
	return aids, nil
}

// GetAnime gets the anime from the database.
func GetAnime(db *sql.DB, aid int) (*Anime, error) {
	r := db.QueryRow(`
SELECT aid, title, type, episodecount, startdate, enddate
FROM anime WHERE aid=?`, aid)
	var a Anime
	if err := r.Scan(&a.AID, &a.Title, &a.Type, &a.EpisodeCount,
		&a.NullStartDate, &a.NullEndDate); err != nil {
		return nil, err
	}
	return &a, nil
}

// GetAllAnime returns all anime.
func GetAllAnime(db *sql.DB) ([]Anime, error) {
	r, err := db.Query(`
SELECT aid, title, type, episodecount, startdate, enddate FROM anime`)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var as []Anime
	for r.Next() {
		var a Anime
		if err := r.Scan(&a.AID, &a.Title, &a.Type, &a.EpisodeCount,
			&a.NullStartDate, &a.NullEndDate); err != nil {
			return nil, err
		}
		as = append(as, a)
	}
	if err := r.Err(); err != nil {
		return nil, err
	}
	return as, nil
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
title=?, type=?, episodecount=?, startdate=?, enddate=?
WHERE aid=?`,
		a.AID, title, a.Type, a.EpisodeCount, startDate, endDate,
		title, a.Type, a.EpisodeCount, startDate, endDate,
		a.AID,
	)
	if err != nil {
		return fmt.Errorf("failed to insert anime %d: %w", a.AID, err)
	}
	em, err := GetEpisodesMap(t, a.AID)
	if err != nil {
		return fmt.Errorf("failed to insert anime %d: %w", a.AID, err)
	}
	for _, e := range a.Episodes {
		k := EpisodeKey{AID: a.AID}
		k.Type, k.Number = parseEpNo(e.EpNo)
		if k.Type == EpUnknown {
			return fmt.Errorf("failed to insert anime %d: invalid epno %s", a.AID, e.EpNo)
		}
		if err := insertEpisode(t, k, e); err != nil {
			return fmt.Errorf("failed to insert episode %s for anime %d: %w",
				e.EpNo, a.AID, err)
		}
		delete(em, k)
	}
	for _, e := range em {
		if err := DeleteEpisode(t, e.ID); err != nil {
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

func insertEpisode(t *sql.Tx, k EpisodeKey, e anidb.Episode) error {
	title := mainEpTitle(e.Titles)
	_, err := t.Exec(`
INSERT INTO episode (eid, aid, type, number, title, length)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT (aid, type, number) DO UPDATE SET
eid=?, title=?, length=?
WHERE aid=? AND type=? AND number=?`,
		e.EID, k.AID, k.Type, k.Number, title, e.Length,
		e.EID, title, e.Length,
		k.AID, k.Type, k.Number,
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
