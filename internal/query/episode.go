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

type Episode struct {
	ID          int
	AID         int
	Type        EpisodeType
	Number      int
	Title       string
	Length      int
	UserWatched bool
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
