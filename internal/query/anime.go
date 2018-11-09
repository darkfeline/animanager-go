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
	"go.felesatra.moe/go2/errors"

	"go.felesatra.moe/animanager/internal/date"
)

type Anime struct {
	AID          int
	Title        string
	Type         AnimeType
	EpisodeCount int
	// The following fields are nullable.  In most cases, use the
	// getter methods instead.
	NullStartDate sql.NullInt64
	NullEndDate   sql.NullInt64
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

// InsertAnime inserts or updates an anime into the database.
func InsertAnime(db *sql.DB, a *anidb.Anime) error {
	t, err := db.Begin()
	if err != nil {
		return err
	}
	defer t.Rollback()
	var startDate interface{}
	var endDate interface{}
	startDate, err = date.NewString(a.StartDate)
	if err != nil {
		startDate = nil
	}
	endDate, err = date.NewString(a.EndDate)
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
		return errors.Wrapf(err, "failed to insert anime %d", a.AID)
	}
	for _, e := range a.Episodes {
		if err := insertEpisode(t, a.AID, e); err != nil {
			return errors.Wrapf(err, "failed to insert episode %s for anime %d",
				e.EpNo, a.AID)
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

func insertEpisode(t *sql.Tx, aid int, e anidb.Episode) error {
	eptype, n := parseEpNo(e.EpNo)
	if eptype == EpInvalid {
		return fmt.Errorf("invalid epno %s", e.EpNo)
	}
	title := mainEpTitle(e.Titles)
	_, err := t.Exec(`
INSERT INTO episode (aid, type, number, title, length)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT (aid, type, number) DO UPDATE SET
title=?, length=?
WHERE aid=? AND type=? AND number=?`,
		aid, eptype, n, title, e.Length,
		title, e.Length,
		aid, eptype, n,
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
