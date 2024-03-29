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
	"fmt"

	"go.felesatra.moe/animanager/internal/sqlc"
)

// GetAnimeFiles gets the episode files for all of the anime's episodes.
func GetAnimeFiles(db sqlc.DBTX, aid sqlc.AID) ([]EpisodeFiles, error) {
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
func DeleteAnimeFiles(db sqlc.DBTX, aid sqlc.AID) error {
	ctx := context.Background()
	err := sqlc.New(db).DeleteAnimeFiles(ctx, aid)
	if err != nil {
		return fmt.Errorf("DeleteAnimeFiles %d: %s", aid, err)
	}
	return nil
}
