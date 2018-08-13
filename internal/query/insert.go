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

	"github.com/pkg/errors"
	"go.felesatra.moe/anidb"
	"go.felesatra.moe/animanager/internal/date"
)

// InsertAnime inserts or updates an anime into the database.
func InsertAnime(db *sql.DB, a *anidb.Anime) error {
	t, err := db.Begin()
	defer t.Rollback()
	startDate, err := date.NewString(a.StartDate)
	if err != nil {
		return errors.Wrapf(err, "failed to insert anime %d", a.AID)
	}
	endDate, err := date.NewString(a.EndDate)
	if err != nil {
		return errors.Wrapf(err, "failed to insert anime %d", a.AID)
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

// InsertWatching inserts or updates a watching entry into the database.
func InsertWatching(db *sql.DB, aid int, regexp string) error {
	t, err := db.Begin()
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

// InsertEpisodeFile inserts a file for an episode into the database.
func InsertEpisodeFile(db *sql.DB, id int, path string) error {
	t, err := db.Begin()
	defer t.Rollback()
	_, err = t.Exec(`INSERT INTO episode_file (episode_id, path) VALUES (?, ?)`, id, path)
	if err != nil {
		return fmt.Errorf("insert episode %d file: %s", id, err)
	}
	return t.Commit()
}
