// Copyright (C) 2020  Allen Li
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

package cmd

import (
	"regexp"
	"testing"

	"go.felesatra.moe/animanager/internal/query"
)

func TestAnimeDefaultRegexp(t *testing.T) {
	t.Parallel()
	a := &query.Anime{
		Title: "Keit-ai 2",
	}
	got := animeDefaultRegexp(a)
	p, err := regexp.Compile(got)
	if err != nil {
		t.Fatalf("Could not compile regexp %#v: %s", got, err)
	}
	cases := []struct {
		String  string
		Version string
	}{
		{"[Meme-raws] Keit-ai 2nd Season - 13 END [720p].mkv", "13"},
	}
	for _, c := range cases {
		t.Run(c.String, func(t *testing.T) {
			t.Parallel()
			got := p.FindStringSubmatch(c.String)
			if got[1] != c.Version {
				t.Errorf("p.FindStringSubmatch(%#v) = %#v (expected ANY, %#v)",
					c.String, got, c.Version)
			}
		})
	}
}
