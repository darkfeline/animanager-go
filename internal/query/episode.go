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

	"golang.org/x/xerrors"
)

type Episode struct {
	_table      struct{}    `episode`
	ID          int         `id`
	AID         int         `aid`
	Type        EpisodeType `type`
	Number      int         `number`
	Title       string      `title`
	Length      int         `length`
	UserWatched bool        `user_watched`
}

func (e Episode) Key() EpisodeKey {
	return EpisodeKey{
		AID:    e.AID,
		Type:   e.Type,
		Number: e.Number,
	}
}

type EpisodeKey struct {
	AID    int
	Type   EpisodeType
	Number int
}

// GetEpisodeCount returns the number of episode rows.
func GetEpisodeCount(db *sql.DB) (int, error) {
	r := db.QueryRow(`SELECT COUNT(*) FROM episode`)
	var n int
	err := r.Scan(&n)
	return n, err
}

// GetWatchedEpisodeCount returns the number of watched episodes.
func GetWatchedEpisodeCount(db *sql.DB) (int, error) {
	r := db.QueryRow(`SELECT COUNT(*) FROM episode where user_watched=1`)
	var n int
	err := r.Scan(&n)
	return n, err
}

// GetWatchedMinutes returns the number of minutes watched.
func GetWatchedMinutes(db *sql.DB) (int, error) {
	r := db.QueryRow(`SELECT SUM(length) FROM episode where user_watched=1`)
	var n int
	err := r.Scan(&n)
	return n, err
}

// GetEpisode gets the episode from the database.
func GetEpisode(db *sql.DB, id int) (*Episode, error) {
	r := db.QueryRow(`
SELECT id, aid, type, number, title, length, user_watched
FROM episode WHERE id=?`, id)
	var e Episode
	if err := r.Scan(&e.ID, &e.AID, &e.Type, &e.Number,
		&e.Title, &e.Length, &e.UserWatched); err != nil {
		return nil, err
	}
	return &e, nil
}

// GetEpisodes gets the episodes for an anime from the database.
func GetEpisodes(db *sql.DB, aid int) ([]Episode, error) {
	r, err := db.Query(`
SELECT id, aid, type, number, title, length, user_watched
FROM episode WHERE aid=? ORDER BY type, number`, aid)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var es []Episode
	for r.Next() {
		var e Episode
		if err := r.Scan(&e.ID, &e.AID, &e.Type, &e.Number,
			&e.Title, &e.Length, &e.UserWatched); err != nil {
			return nil, err
		}
		es = append(es, e)
	}
	if err := r.Err(); err != nil {
		return nil, err
	}
	return es, nil
}

// GetEpisodesMap returns a map of the episodes for an anime.
func GetEpisodesMap(db *sql.DB, aid int) (map[EpisodeKey]*Episode, error) {
	es, err := GetEpisodes(db, aid)
	if err != nil {
		return nil, xerrors.Errorf("get episodes map %v: %w", aid, err)
	}
	m := make(map[EpisodeKey]*Episode, len(es))
	for i, e := range es {
		m[e.Key()] = &es[i]
	}
	return m, nil
}

// GetAllEpisodes gets all episodes from the database.
func GetAllEpisodes(db *sql.DB) ([]Episode, error) {
	r, err := db.Query(`
SELECT id, aid, type, number, title, length, user_watched
FROM episode`)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var es []Episode
	for r.Next() {
		var e Episode
		if err := r.Scan(&e.ID, &e.AID, &e.Type, &e.Number,
			&e.Title, &e.Length, &e.UserWatched); err != nil {
			return nil, err
		}
		es = append(es, e)
	}
	if err := r.Err(); err != nil {
		return nil, err
	}
	return es, nil
}

// UpdateEpisodeDone updates the episode's done status.
func UpdateEpisodeDone(db *sql.DB, id int, done bool) error {
	var watched uint8
	if done {
		watched = 1
	} else {
		watched = 0
	}
	t, err := db.Begin()
	if err != nil {
		return err
	}
	defer t.Rollback()
	_, err = t.Exec(`UPDATE episode SET user_watched=? WHERE id=?`,
		watched, id)
	if err != nil {
		return fmt.Errorf("update episode %d done: %s", id, err)
	}
	if err := t.Commit(); err != nil {
		return fmt.Errorf("update episode %d done: %s", id, err)
	}
	return nil
}
