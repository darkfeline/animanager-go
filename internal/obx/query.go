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

package obx

import (
	"database/sql"

	"go.felesatra.moe/go2/errors"

	"go.felesatra.moe/animanager/internal/query"
)

// GetAnimeFiles gets the episode files for the anime's episodes.
func GetAnimeFiles(db *sql.DB, aid int) ([]EpisodeFiles, error) {
	eps, err := query.GetEpisodes(db, aid)
	if err != nil {
		return nil, errors.Wrapf(err, "get anime %d files", aid)
	}
	var efs []EpisodeFiles
	for _, e := range eps {
		ef := EpisodeFiles{
			Episode: e,
		}
		fs, err := query.GetEpisodeFiles(db, e.ID)
		if err != nil {
			return nil, errors.Wrapf(err, "get anime %d files", aid)
		}
		ef.Files = fs
		efs = append(efs, ef)
	}
	return efs, nil
}

type EpisodeFiles struct {
	Episode query.Episode
	Files   []query.EpisodeFile
}

// GetCompletedAnime returns completed anime.
func GetCompletedAnime(db *sql.DB, incomplete bool) ([]query.Anime, error) {
	as, err := query.GetAllAnime(db)
	if err != nil {
		return nil, err
	}
	es, err := query.GetAllEpisodes(db)
	if err != nil {
		return nil, err
	}
	am := make(map[int]*query.Anime)
	counts := make(map[int]int)
	for i, a := range as {
		counts[a.AID] = a.EpisodeCount
		am[a.AID] = &as[i]
	}
	for _, e := range es {
		if e.UserWatched && e.Type == query.EpRegular {
			counts[e.AID]--
		}
	}
	var res []query.Anime
	var test func(int) bool
	if incomplete {
		test = func(count int) bool { return count > 0 }
	} else {
		test = func(count int) bool { return count <= 0 }
	}
	for aid, count := range counts {
		if test(count) {
			res = append(res, *am[aid])

		}
	}
	return res, nil
}
