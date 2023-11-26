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
)

// An Executor executes queries.
// This interface is used for functions that can be used by both
// sql.DB and sql.Tx.
type Executor interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type Episode struct {
	_table      struct{}    `sql:"episode"`
	ID          EpID        `sql:"id"`
	EID         EID         `sql:"eid"`
	AID         AID         `sql:"aid"`
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

const selectEpisodeFields = `id, eid, aid, type, number, title, length, user_watched`

// A Scanner supports both sql.Row and sql.Rows.
type Scanner interface {
	Scan(dest ...interface{}) error
}

func scanEpisode(r Scanner, e *Episode) error {
	return r.Scan(&e.ID, &e.EID, &e.AID, &e.Type, &e.Number,
		&e.Title, &e.Length, &e.UserWatched)
}

// EpisodeKey represents the unique key for an Episode.  This is
// separate from ID because SQLite treats numeric row IDs specially.
type EpisodeKey struct {
	AID    AID
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
func GetEpisode(db *sql.DB, id EpID) (*Episode, error) {
	r := db.QueryRow(`
SELECT `+selectEpisodeFields+`
FROM episode WHERE id=?`, id)
	var e Episode
	if err := scanEpisode(r, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

// GetEpisodeByKey gets an episode from the database by unique key.
func GetEpisodeByKey(db *sql.DB, k EpisodeKey) (*Episode, error) {
	r := db.QueryRow(`
SELECT `+selectEpisodeFields+`
FROM episode WHERE aid=? AND type=? AND number=?`, k.AID, k.Type, k.Number)
	var e Episode
	if err := scanEpisode(r, &e); err != nil {
		return nil, err
	}
	return &e, nil

}

// DeleteEpisode deletes the episode from the database.
func DeleteEpisode(db Executor, id EpID) error {
	if _, err := db.Exec(`DELETE FROM episode WHERE id=?`, id); err != nil {
		return fmt.Errorf("delete episode %v: %w", id, err)
	}
	return nil
}

// GetEpisodes gets the episodes for an anime from the database.
func GetEpisodes(db Executor, aid AID) ([]Episode, error) {
	r, err := db.Query(`
SELECT `+selectEpisodeFields+`
FROM episode WHERE aid=? ORDER BY type, number`, aid)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var es []Episode
	for r.Next() {
		var e Episode
		if err := scanEpisode(r, &e); err != nil {
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
func GetEpisodesMap(db Executor, aid AID) (map[EpisodeKey]*Episode, error) {
	es, err := GetEpisodes(db, aid)
	if err != nil {
		return nil, fmt.Errorf("get episodes map %v: %w", aid, err)
	}
	m := make(map[EpisodeKey]*Episode, len(es))
	for i, e := range es {
		m[e.Key()] = &es[i]
	}
	return m, nil
}

// GetAllEpisodes gets all episodes from the database.
func GetAllEpisodes(db *sql.DB) ([]Episode, error) {
	r, err := db.Query(`SELECT ` + selectEpisodeFields + ` FROM episode`)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var es []Episode
	for r.Next() {
		var e Episode
		if err := scanEpisode(r, &e); err != nil {
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
func UpdateEpisodeDone(db *sql.DB, id EpID, done bool) error {
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
		return fmt.Errorf("update episode %d done: %w", id, err)
	}
	if err := t.Commit(); err != nil {
		return fmt.Errorf("update episode %d done: %w", id, err)
	}
	return nil
}

// GetAIDsMissingEIDs is a temporary function to support EID backfill.
func GetAIDsMissingEIDs(db *sql.DB) ([]AID, error) {
	r, err := db.Query(`SELECT DISTINCT aid FROM episode WHERE eid IS NULL`)
	if err != nil {
		return nil, err
	}
	var aids []AID
	for r.Next() {
		var aid AID
		if err := r.Scan(&aid); err != nil {
			return nil, err
		}
		aids = append(aids, aid)
	}
	return aids, nil
}
