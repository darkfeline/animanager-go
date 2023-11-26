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

package query

import (
	"fmt"
	"strconv"
)

// An AID is an ID for [Anime].
type AID int

// Scan implements [database/sql.Scanner].
func (t *AID) Scan(src any) error {
	return scanID(t, src)
}

// An EpID is an ID for [Episode].
type EpID int

// Scan implements [database/sql.Scanner].
func (t *EpID) Scan(src any) error {
	return scanID(t, src)
}

// An EID is an ID for [Episode].
type EID int

// Scan implements [database/sql.Scanner].
func (t *EID) Scan(src any) error {
	return scanID(t, src)
}

// scanID is a helper for implementing [database/sql.Scanner] for
// custom int types.
func scanID[T ~int](t *T, src any) error {
	v, ok := src.(int64)
	if !ok {
		return fmt.Errorf("wrong type %T for %T", src, *t)
	}
	v2 := T(v)
	if int64(v2) != v {
		return fmt.Errorf("value does not fit in %T: %v", *t, src)
	}
	*t = v2
	return nil
}

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
