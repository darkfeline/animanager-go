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
	"fmt"
	"strconv"
	"strings"

	"go.felesatra.moe/animanager/internal/sqlc"
)

type EpisodeType = sqlc.EpisodeType

const (
	EpUnknown EpisodeType = iota
	EpRegular
	EpSpecial
	EpCredit
	EpTrailer
	EpParody
	EpOther
)

type prefixPair struct {
	Prefix string
	Type   EpisodeType
}

var epnoPrefixes = []prefixPair{
	{"S", EpSpecial},
	{"C", EpCredit},
	{"T", EpTrailer},
	{"P", EpParody},
	{"O", EpOther},
	{"", EpRegular},
}

// parseEpNo parses episode number information from the AniDB format.
// If parse fails, EpUnknown is returned for the episode type.
func parseEpNo(epno string) (EpisodeType, int) {
	for _, p := range epnoPrefixes {
		if strings.HasPrefix(epno, p.Prefix) {
			n, err := strconv.Atoi(epno[len(p.Prefix):])
			if err != nil {
				return EpUnknown, 0
			}
			return p.Type, n
		}
	}
	panic(fmt.Sprintf("ParseEpNo %s unreachable code", epno))
}
