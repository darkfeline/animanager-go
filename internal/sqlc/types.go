// Copyright (C) 2024  Allen Li
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

package sqlc

import (
	"fmt"
	"strconv"
)

// An AID is an ID for [Anime].
type AID int

// An EID is an ID for [Episode].
type EID int

// ParseIDs parses multiple IDs using [ParseID].
func ParseIDs[T ~int](args []string) ([]T, error) {
	ids := make([]T, len(args))
	for i, s := range args {
		id, err := ParseID[T](s)
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}
	return ids, nil
}

// ParseID parses an ID type like [AID].
func ParseID[T ~int](s string) (T, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %q into %T: %s", s, T(0), err)
	}
	return T(id), nil
}

// A Hash is an eD2k hash formatted as a hex string.
type Hash string

// An EpisodeType is the type for an [Episode].
type EpisodeType int

const (
	EpRegular EpisodeType = 1
	EpSpecial EpisodeType = 2
	EpCredit  EpisodeType = 3
	EpTrailer EpisodeType = 4
	EpParody  EpisodeType = 5
	EpOther   EpisodeType = 6
)

//go:generate stringer -type=EpisodeType

func (t EpisodeType) Prefix() string {
	switch t {
	case EpRegular:
		return ""
	case EpSpecial:
		return "S"
	case EpCredit:
		return "C"
	case EpTrailer:
		return "T"
	case EpParody:
		return "P"
	case EpOther:
		return "O"
	default:
		panic(fmt.Sprintf("invalid %T=%v", t, t))
	}
}

// Valid returns whether the [EpisodeType] is a valid value.
func (t EpisodeType) Valid() bool {
	return 1 <= t && t <= 6
}

// ParseEpisodeType parses an [EpisodeType] from the string prefix.
// The string is usually an episode number in AniDB format.
// If parse fails, returns an invalid [EpisodeType].
// Note that this function may return [EpRegular] even if the string
// is an invalid episode number, as this function does not assume
// AniDB format.
func ParseEpisodeType(epno string) EpisodeType {
	if len(epno) < 1 {
		return 0
	}
	switch epno[:1] {
	case EpSpecial.Prefix():
		return EpSpecial
	case EpCredit.Prefix():
		return EpCredit
	case EpTrailer.Prefix():
		return EpTrailer
	case EpParody.Prefix():
		return EpParody
	case EpOther.Prefix():
		return EpOther
	}
	return EpRegular
}
