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

package query

import (
	"strconv"

	"go.felesatra.moe/animanager/internal/sqlc"
)

type EpisodeType = sqlc.EpisodeType

const (
	EpRegular EpisodeType = 1
	EpSpecial EpisodeType = 2
	EpCredit  EpisodeType = 3
	EpTrailer EpisodeType = 4
	EpParody  EpisodeType = 5
	EpOther   EpisodeType = 6
)

// parseEpNo parses episode number information from the AniDB format.
// If parse fails, returns an invalid EpisodeType.
func parseEpNo(epno string) (EpisodeType, int) {
	t := sqlc.ParseEpisodeType(epno)
	if !t.Valid() {
		return 0, 0
	}
	n, err := strconv.Atoi(epno[len(t.Prefix()):])
	if err != nil {
		return 0, 0
	}
	return t, n
}
