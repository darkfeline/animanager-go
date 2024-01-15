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
