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

	"go.felesatra.moe/animanager/internal/date"
)

type Anime struct {
	AID          int
	Title        string
	Type         AnimeType
	EpisodeCount int
	// The following fields are nullable.  In most cases, use the
	// getter methods instead.
	NStartDate interface{}
	NEndDate   interface{}
}

// StartDate returns the nullable NStartDate field as a Date.  If the
// field is null, returns date.Zero.
func (a Anime) StartDate() date.Date {
	switch d := a.NStartDate.(type) {
	case int64:
		return date.Date(d)
	case nil:
		return date.Zero
	default:
		panic(fmt.Sprintf("bad field type %T %#v", d, d))
	}
}

// EndDate returns the nullable NEndDate field as a Date.  If the
// field is null, return the zero value.
func (a Anime) EndDate() date.Date {
	switch d := a.NEndDate.(type) {
	case int64:
		return date.Date(d)
	case nil:
		return date.Zero
	default:
		panic(fmt.Sprintf("bad field type %T %#v", d, d))
	}
}

type AnimeType string
