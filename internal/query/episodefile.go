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
	"log/slog"
)

type EpisodeFile struct {
	_table struct{} `sql:"episode_file"`
	EID    EID      `sql:"eid"`
	Path   string   `sql:"path"`
}

// InsertEpisodeFile inserts episode files into the database.
func InsertEpisodeFiles(db *sql.DB, l *slog.Logger, efs []EpisodeFile) error {
	s, err := db.Prepare(`INSERT INTO episode_file (eid, path) VALUES (?, ?)`)
	if err != nil {
		return fmt.Errorf("insert episode files: %w", err)
	}
	defer s.Close()
	for _, ef := range efs {
		if _, err = s.Exec(ef.EID, ef.Path); err != nil {
			// This is most likely due to EID foreign key error.
			l.Debug("error inserting episode file", "EpisodeFile", ef, "error", err)
		}
	}
	return nil
}

// GetEpisodeFiles returns the EpisodeFiles for the episode.
func GetEpisodeFiles(db *sql.DB, eid EID) (es []EpisodeFile, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("get episode %d files: %w", eid, err)
		}
	}()
	r, err := db.Query(`
SELECT eid, path
FROM episode_file WHERE eid=?`, eid)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	for r.Next() {
		var e EpisodeFile
		if err := r.Scan(&e.EID, &e.Path); err != nil {
			return nil, err
		}
		es = append(es, e)
	}
	if r.Err() != nil {
		return nil, r.Err()
	}
	return es, nil
}

// DeleteAllEpisodeFiles deletes all episode files.
func DeleteAllEpisodeFiles(db Executor) error {
	_, err := db.Exec(`DELETE FROM episode_file`)
	if err != nil {
		return err
	}
	return nil
}
