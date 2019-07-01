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

import (
	"testing"

	"go.felesatra.moe/animanager/internal/query"
)

func TestIsUnnamed(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		title string
		want  bool
	}{
		{"unnamed", "Episode 1", true},
		{"named", "挑・発", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			e := query.Episode{Title: c.title}
			got := isUnnamed(e)
			if got != c.want {
				t.Errorf("isUnnamed(%#v) = %#v; want %#v", e, got, c.want)
			}
		})
	}
}
