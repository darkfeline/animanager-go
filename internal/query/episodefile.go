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
	"log/slog"

	"go.felesatra.moe/animanager/internal/sqlc"
)

type EpisodeFile = sqlc.EpisodeFile

// InsertEpisodeFile inserts episode files into the database.
// The caller should pre-emptively prepare
// [sqlc.Queries.PrepareInsertEpisodeFile] for performance.
func InsertEpisodeFiles(ctx context.Context, q *sqlc.Queries, l *slog.Logger, efs []EpisodeFile) error {
	for _, ef := range efs {
		p := sqlc.InsertEpisodeFileParams{
			Eid:  ef.Eid,
			Path: ef.Path,
		}
		if err := q.InsertEpisodeFile(ctx, p); err != nil {
			// This is most likely due to EID foreign key error.
			l.Debug("error inserting episode file", "EpisodeFile", ef, "error", err)
		}
	}
	return nil
}

// GetEpisodeFiles returns the EpisodeFiles for the episode.
func GetEpisodeFiles(db sqlc.DBTX, eid sqlc.EID) ([]EpisodeFile, error) {
	ctx := context.Background()
	es, err := sqlc.New(db).GetEpisodeFiles(ctx, eid)
	if err != nil {
		return nil, fmt.Errorf("GetEpisodeFiles %d: %s", eid, err)
	}
	return es, nil
}

// DeleteAllEpisodeFiles deletes all episode files.
func DeleteAllEpisodeFiles(db sqlc.DBTX) error {
	ctx := context.Background()
	return sqlc.New(db).DeleteAllEpisodeFiles(ctx)
}
