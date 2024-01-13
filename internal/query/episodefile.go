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
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"go.felesatra.moe/animanager/internal/sqlc"
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
func GetEpisodeFiles(db sqlc.DBTX, eid EID) ([]EpisodeFile, error) {
	ctx := context.Background()
	es, err := sqlc.New(db).GetEpisodeFiles(ctx, int64(eid))
	if err != nil {
		return nil, fmt.Errorf("GetEpisodeFiles %d: %s", eid, err)
	}
	return smap(es, convertEpisodeFile), nil
}

// DeleteAllEpisodeFiles deletes all episode files.
func DeleteAllEpisodeFiles(db Executor) error {
	_, err := db.Exec(`DELETE FROM episode_file`)
	if err != nil {
		return err
	}
	return nil
}

func convertEpisodeFile(e sqlc.EpisodeFile) EpisodeFile {
	return EpisodeFile{
		EID:  EID(e.Eid),
		Path: e.Path,
	}
}
