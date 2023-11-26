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

package fileid

import (
	"reflect"
	"testing"

	"go.felesatra.moe/animanager/internal/query"
)

func TestFilterFiles(t *testing.T) {
	w := query.Watching{
		Regexp: "lacia([0-9]+)",
	}
	eps := []query.Episode{
		{EID: 111, Type: query.EpRegular, Number: 1},
		{EID: 112, Type: query.EpRegular, Number: 3},
		{EID: 113, Type: query.EpRegular, Number: 5},
		{EID: 114, Type: query.EpOther, Number: 13},
	}
	files := []string{
		"/foo/lacia1",
		"/foo/lacia1v2",
		"/foo/lacia2",
		"/foo/foobar",
		"/foo/lacia5",
		"/foo/lacia13",
	}
	got, err := filterFiles(w, eps, files)
	if err != nil {
		t.Errorf("filterFiles returned error: %#v", err)
	}
	want := []query.EpisodeFile{
		{EID: 111, Path: "/foo/lacia1"},
		{EID: 111, Path: "/foo/lacia1v2"},
		{EID: 0, Path: "/foo/lacia2"},
		{EID: 113, Path: "/foo/lacia5"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("filterFiles() = %#v; want %#v", got, want)
	}
}

func TestMaxEpisodeNumber(t *testing.T) {
	eps := []query.Episode{
		{Type: query.EpRegular, Number: 1},
		{Type: query.EpRegular, Number: 3},
		{Type: query.EpRegular, Number: 5},
		{Type: query.EpOther, Number: 13},
	}
	got := maxEpisodeNumber(eps)
	want := 5
	if got != want {
		t.Errorf("maxEpisodeNumber(%#v) = %#v; want %#v", eps, got, want)
	}
}
