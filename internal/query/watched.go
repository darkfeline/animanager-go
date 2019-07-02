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
)

// GetWatchedAnime returns watched anime (all episodes watched).
func GetWatchedAnime(db *sql.DB) ([]Anime, error) {
	return getAnimeByWatched(db, true)
}

// GetUnwatchedAnime returns unwatched anime (not all episodes watched).
func GetUnwatchedAnime(db *sql.DB) ([]Anime, error) {
	return getAnimeByWatched(db, false)
}

// getAnimeByWatched returns anime by episode user_watched status.
func getAnimeByWatched(db *sql.DB, watched bool) ([]Anime, error) {
	as, err := GetAllAnime(db)
	if err != nil {
		return nil, err
	}
	es, err := GetAllEpisodes(db)
	if err != nil {
		return nil, err
	}
	am := make(map[int]*Anime)
	counts := make(map[int]int)
	for i, a := range as {
		counts[a.AID] = a.EpisodeCount
		am[a.AID] = &as[i]
	}
	for _, e := range es {
		if e.UserWatched && e.Type == EpRegular {
			counts[e.AID]--
		}
	}
	var res []Anime
	var test func(int) bool
	if watched {
		test = func(count int) bool { return count <= 0 }
	} else {
		test = func(count int) bool { return count > 0 }
	}
	for aid, count := range counts {
		if test(count) {
			res = append(res, *am[aid])

		}
	}
	return res, nil
}
