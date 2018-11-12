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

import "go.felesatra.moe/animanager/internal/query"

// MakeEpisodeMap returns a slice where each index maps to the regular
// episode with the same number.  The zero index will be empty since
// episodes cannot be number zero.
func MakeEpisodeMap(eps []query.Episode) []query.Episode {
	m := make([]query.Episode, maxEpisodeNumber(eps)+1)
	for _, e := range eps {
		if e.Type == query.EpRegular {
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
