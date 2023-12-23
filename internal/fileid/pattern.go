// Copyright (C) 2023  Allen Li
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

package fileid

import (
	"database/sql"
	"fmt"
	"log/slog"
	"path/filepath"
	"regexp"
	"strconv"

	"go.felesatra.moe/animanager/internal/query"
)

// RefreshFiles updates all episode files in the database using
// [query.Watching] regexp patterns and the given video file paths.
func RefreshFiles(db *sql.DB, files []string) error {
	l := slog.Default().With("method", "pattern")
	ws, err := query.GetAllWatching(db)
	if err != nil {
		return fmt.Errorf("refresh files: %w", err)
	}
	var efs []query.EpisodeFile
	l.Info("start match files")
	for _, w := range ws {
		if err := query.DeleteAnimeFiles(db, w.AID); err != nil {
			return fmt.Errorf("refresh files: %w", err)
		}
		l.Info("start match files for anime", "aid", w.AID)
		eps, err := query.GetEpisodes(db, w.AID)
		if err != nil {
			return fmt.Errorf("refresh files: %w", err)
		}
		efs2, err := filterFiles(w, eps, files)
		if err != nil {
			return fmt.Errorf("refresh files: %w", err)
		}
		l.Info("matched files for anime", "aid", w.AID, "files", efs2)
		efs = append(efs, efs2...)
	}
	l.Info("start insert files")
	if err = query.InsertEpisodeFiles(db, slog.Default(), efs); err != nil {
		return fmt.Errorf("refresh files: %w", err)
	}
	return nil
}

// filterFiles returns files that match the [query.Watching] entry.
func filterFiles(w query.Watching, eps []query.Episode, files []string) ([]query.EpisodeFile, error) {
	var result []query.EpisodeFile
	r, err := regexp.Compile("(?i)" + w.Regexp)
	if err != nil {
		return nil, fmt.Errorf("filter files for %d: %w", w.AID, err)
	}
	regEps := reindexEpisodeByNumber(eps)
	for _, f := range files {
		ms := r.FindStringSubmatch(filepath.Base(f))
		if ms == nil {
			continue
		}
		if len(ms) < 2 {
			return nil, fmt.Errorf("filter files for %d: regexp %#v has no submatch",
				w.AID, w.Regexp)

		}
		n, err := strconv.Atoi(ms[1])
		if err != nil {
			return nil, fmt.Errorf("filter files for %d: regexp %#v submatch not a number",
				w.AID, w.Regexp)
		}
		n += w.Offset
		if n >= len(regEps) || n < 1 {
			continue
		}
		result = append(result, query.EpisodeFile{
			EID:  regEps[n].EID,
			Path: f,
		})
	}
	return result, nil
}

// reindexEpisodeByNumber returns a slice where each index maps to the regular
// episode with the same number.  The zero index will be empty since
// episodes cannot be number zero.
func reindexEpisodeByNumber(eps []query.Episode) []query.Episode {
	l := slog.Default().With("func", "reindexEpisodeByNumber")
	m := make([]query.Episode, maxEpisodeNumber(eps)+1)
	for _, e := range eps {
		if e.Type == query.EpRegular {
			if old := m[e.Number]; old.EID != 0 {
				l.Warn("duplicate regular episode", "old", old, "new", e)
			}
			m[e.Number] = e
		}
	}
	return m
}

func maxEpisodeNumber(eps []query.Episode) int {
	maxEp := 0
	for _, e := range eps {
		if e.Type == query.EpRegular && e.Number > maxEp {
			maxEp = e.Number
		}
	}
	return maxEp
}
