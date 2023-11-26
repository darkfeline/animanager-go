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

// GetAnimeFiles gets the episode files for all of the anime's episodes.
func GetAnimeFiles(db *sql.DB, aid AID) ([]EpisodeFiles, error) {
	eps, err := GetEpisodes(db, aid)
	if err != nil {
		return nil, fmt.Errorf("get anime %d files: %w", aid, err)
	}
	var efs []EpisodeFiles
	for _, e := range eps {
		ef := EpisodeFiles{
			Episode: e,
		}
		fs, err := GetEpisodeFiles(db, e.EID)
		if err != nil {
			return nil, fmt.Errorf("get anime %d files: %w", aid, err)
		}
		ef.Files = fs
		efs = append(efs, ef)
	}
	return efs, nil
}

type EpisodeFiles struct {
	Episode Episode
	Files   []EpisodeFile
}

// DeleteAnimeFiles deletes episode files for the given anime.
func DeleteAnimeFiles(db Executor, aid AID) error {
	_, err := db.Exec(`DELETE FROM episode_file
WHERE ROWID IN (
    SELECT episode_file.ROWID FROM episode_file
    JOIN episode ON (episode_file.eid = episode.eid)
    WHERE episode.aid=?
)`, aid)
	if err != nil {
		return err
	}
	return nil
}
