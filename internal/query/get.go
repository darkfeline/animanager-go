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

	"github.com/pkg/errors"
)

// ErrMissing is returned from queries that fetch a unique item of the
// keyed item is missing.
var ErrMissing = errors.New("row missing")

// GetAnime gets the anime from the database.  ErrMissing is returned
// if the anime doesn't exist.
func GetAnime(db *sql.DB, aid int) (*Anime, error) {
	r, err := db.Query(`
SELECT aid, title, type, episodecount, startdate, enddate
FROM anime WHERE aid=?`, aid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query anime")
	}
	defer r.Close()
	if !r.Next() {
		if r.Err() != nil {
			return nil, r.Err()
		}
		return nil, ErrMissing
	}
	a := Anime{}
	if err := r.Scan(&a.AID, &a.Title, &a.Type,
		&a.EpisodeCount, &a.StartDate, &a.EndDate); err != nil {
		return nil, errors.Wrap(err, "failed to scan anime")
	}
	return &a, nil
}

// GetEpisodes gets the episodes for an anime from the database.
func GetEpisodes(db *sql.DB, aid int) ([]Episode, error) {
	r, err := db.Query(`
SELECT id, aid, type, number, title, length, user_watched
FROM episode WHERE aid=? ORDER BY type, number`, aid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query episode")
	}
	defer r.Close()
	var es []Episode
	for r.Next() {
		e := Episode{}
		if err := r.Scan(&e.ID, &e.AID, &e.Type, &e.Number,
			&e.Title, &e.Length, &e.UserWatched); err != nil {
			return nil, errors.Wrap(err, "failed to scan episode")
		}
		es = append(es, e)
	}
	if err := r.Err(); err != nil {
		return nil, err
	}
	return es, nil
}

// GetWatching gets the watching entry for an anime from the
// database. ErrMissing is returned if the anime doesn't exist.
func GetWatching(db *sql.DB, aid int) (regexp string, err error) {
	r, err := db.Query(`SELECT regexp FROM watching WHERE aid=?`, aid)
	if err != nil {
		return "", errors.Wrap(err, "failed to query watching")
	}
	defer r.Close()
	if !r.Next() {
		if r.Err() != nil {
			return "", r.Err()
		}
		return "", ErrMissing
	}
	if err := r.Scan(&regexp); err != nil {
		return "", errors.Wrap(err, "failed to scan episode")
	}
	return regexp, nil
}
