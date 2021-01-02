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

package main

import (
	"reflect"
	"regexp"
	"testing"

	"go.felesatra.moe/anidb"
)

func TestTitleMatchesTrue(t *testing.T) {
	at := []anidb.Title{
		{Name: "sophie plachta"},
		{Name: "hikari homura"},
	}
	r := regexp.MustCompile("homu")
	got := titleMatches(r, at)
	if !got {
		t.Errorf("%#v did not match %#v", r, at)
	}
}

func TestTitleMatchesFalse(t *testing.T) {
	at := []anidb.Title{
		{Name: "sophie plachta"},
		{Name: "hikari homura"},
	}
	r := regexp.MustCompile("mayu")
	got := titleMatches(r, at)
	if got {
		t.Errorf("%#v matched %#v", r, at)
	}
}

func TestFilterTitles(t *testing.T) {
	at1 := anidb.AnimeT{
		AID: 1,
		Titles: []anidb.Title{
			{Name: "sophie plachta"},
			{Name: "hikari homura"},
		},
	}
	at2 := anidb.AnimeT{
		AID: 2,
		Titles: []anidb.Title{
			{Name: "ramius riche"},
			{Name: "gurigura kath"},
		},
	}

	r := regexp.MustCompile("homu")
	got := filterTitles(r, []anidb.AnimeT{at1, at2})
	exp := []anidb.AnimeT{at1}
	if !reflect.DeepEqual(got, exp) {
		t.Errorf("Expected %#v, got %#v", exp, got)
	}
}
