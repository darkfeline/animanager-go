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

// GetFinishedAnime returns finished anime.
// Finished anime have all EpRegular episodes marked as user_watched.
func GetFinishedAnime(db *sql.DB) ([]*Anime, error) {
	anime, err := GetAnimeFinished(db)
	if err != nil {
		return nil, err
	}
	res := make([]*Anime, 0, len(anime))
	for _, v := range anime {
		if v.Value {
			res = append(res, v.Anime)
		}
	}
	return res, nil
}

// GetUnfinishedAnime returns unfinished anime.
// Unfinished anime either are incomplete or don't have all EpRegular
// episodes marked as user_watched.
func GetUnfinishedAnime(db *sql.DB) ([]*Anime, error) {
	anime, err := GetAnimeFinished(db)
	if err != nil {
		return nil, err
	}
	res := make([]*Anime, 0, len(anime))
	for _, v := range anime {
		if !v.Value {
			res = append(res, v.Anime)
		}
	}
	return res, nil
}

// An AnimeBool is just an anime with a bool field extension.
// Used as a result from functions that partition anime into two buckets.
type AnimeBool struct {
	*Anime
	Value bool
}

func (b AnimeBool) GoString() string {
	return fmt.Sprintf("query.AnimeBool{Anime:%#v, Value:%t}", b.Anime, b.Value)
}

// GetAnimeFinished returns all anime annotated with whether they are finished.
// Finished anime have all EpRegular episodes marked as user_watched.
func GetAnimeFinished(db *sql.DB) ([]AnimeBool, error) {
	anime, err := GetAllAnime(db)
	if err != nil {
		return nil, fmt.Errorf("get anime finished: %s", err)
	}
	res := make([]AnimeBool, 0, len(anime))
	for i := range anime {
		eps, err := GetEpisodes(db, anime[i].AID)
		if err != nil {
			return nil, fmt.Errorf("get anime finished: %s", err)
		}
		res = append(res, AnimeBool{
			Anime: &anime[i],
			Value: isAnimeFinished(&anime[i], eps),
		})
	}
	return res, nil
}

// Returns whether anime is finished.
// Finished anime have all EpRegular episodes marked as user_watched.
func isAnimeFinished(a *Anime, eps []Episode) bool {
	if isIncomplete2(a, eps) {
		return false
	}
	watched := 0
	for _, e := range eps {
		if e.UserWatched && e.Type == EpRegular {
			watched++
		}
	}
	if watched < a.EpisodeCount {
		return false
	}
	return true
}
