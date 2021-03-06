// Copyright (C) 2019  Allen Li
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package query

import (
	"database/sql"
	"fmt"
)

type EpisodeFile struct {
	_table    struct{} `episode_file`
	EpisodeID int      `episode_id`
	Path      string   `path`
}

// InsertEpisodeFile inserts episode files into the database.
func InsertEpisodeFiles(db *sql.DB, efs []EpisodeFile) error {
	t, err := db.Begin()
	if err != nil {
		return fmt.Errorf("insert episode files: %w", err)
	}
	defer t.Rollback()
	s, err := t.Prepare(`INSERT INTO episode_file (episode_id, path) VALUES (?, ?)`)
	if err != nil {
		return fmt.Errorf("insert episode files: %w", err)
	}
	for _, ef := range efs {
		if _, err = s.Exec(ef.EpisodeID, ef.Path); err != nil {
			return fmt.Errorf("insert episode files: %w", err)
		}
	}
	if err := t.Commit(); err != nil {
		return fmt.Errorf("insert episode files: %w", err)
	}
	return nil
}

// GetEpisodeFiles returns the EpisodeFiles for the episode.
func GetEpisodeFiles(db *sql.DB, episodeID int) (es []EpisodeFile, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("get episode %d files: %w", episodeID, err)
		}
	}()
	r, err := db.Query(`
SELECT episode_id, path
FROM episode_file WHERE episode_id=?`, episodeID)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	for r.Next() {
		var e EpisodeFile
		if err := r.Scan(&e.EpisodeID, &e.Path); err != nil {
			return nil, err
		}
		es = append(es, e)
	}
	if r.Err() != nil {
		return nil, r.Err()
	}
	return es, nil
}

// DeleteEpisodeFiles deletes all episode files.
func DeleteEpisodeFiles(db *sql.DB) error {
	t, err := db.Begin()
	if err != nil {
		return err
	}
	defer t.Rollback()
	_, err = t.Exec(`DELETE FROM episode_file`)
	if err != nil {
		return err
	}
	return t.Commit()
}
