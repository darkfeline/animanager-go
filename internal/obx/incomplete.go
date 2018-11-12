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
	"fmt"

	"go.felesatra.moe/animanager/internal/query"
)

// GetIncompleteAnime returns the AIDs for incomplete anime.
func GetIncompleteAnime(db *sql.DB) ([]int, error) {
	aids, err := query.GetAIDs(db)
	if err != nil {
		return nil, fmt.Errorf("get incomplete anime: %s", err)
	}
	var r []int
	for _, aid := range aids {
		ok, err := isIncomplete(db, aid)
		if err != nil {
			return nil, fmt.Errorf("get incomplete anime: %s", err)
		}
		if ok {
			r = append(r, aid)
		}
	}
	return r, nil
}

func isIncomplete(db *sql.DB, aid int) (bool, error) {
	a, err := query.GetAnime(db, aid)
	if err != nil {
		return false, fmt.Errorf("check %d completion: %s", aid, err)
	}
	eps, err := query.GetEpisodes(db, aid)
	if err != nil {
		return false, fmt.Errorf("check %d completion: %s", aid, err)
	}
	var rEps []query.Episode
	var unnamed int
	for _, e := range eps {
		if e.Type != query.EpRegular {
			continue
		}
		rEps = append(rEps, e)
		if isUnnamed(e) {
			unnamed += 1
		} else {
			// Unnamed episodes followed by named episodes
			// are probably just missing the episode title
			// entirely, so don't count them.
			unnamed = 0
		}
	}
	if len(rEps) < a.EpisodeCount {
		return true, nil
	}
	// This is just a heuristic, some shows only have titles for
	// first/last episode.
	if unnamed > 0 && unnamed < a.EpisodeCount-2 {
		return true, nil
	}
	return false, nil
}

func isUnnamed(e query.Episode) bool {
	return len(e.Title) > 8 && e.Title[:8] == "Episode "
}
