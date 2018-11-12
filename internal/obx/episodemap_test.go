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
