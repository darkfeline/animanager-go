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
	_, err = t.Exec(`
INSERT INTO anime (aid, title, type, episodecount, startdate, enddate)
VALUES (?, ?, ?, ?, ?, ?)`, a.AID, mainTitle(a.Titles), a.Type, a.EpisodeCount,
		startDate, endDate)
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
	_, err := t.Exec(`
INSERT INTO episode (aid, type, number, title, length)
VALUES (?, ?, ?, ?, ?)`, aid, eptype, n, mainEpTitle(e.Titles), e.Length)
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
