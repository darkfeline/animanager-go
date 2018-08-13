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
	"testing"
)

func TestParseEpNo(t *testing.T) {
	t.Parallel()
	cases := []struct {
		EpNo   string
		Type   EpisodeType
		Number int
	}{
		{"S1", EpSpecial, 1},
		{"T2", EpTrailer, 2},
		{"15", EpRegular, 15},
		{"Clarion", EpInvalid, 0},
	}
	for _, c := range cases {
		t.Run(c.EpNo, func(t *testing.T) {
			t.Parallel()
			eptype, n := parseEpNo(c.EpNo)
			if eptype != c.Type || n != c.Number {
				t.Errorf("ParseEpNo(%s) = %s, %d (expected %s, %d)",
					c.EpNo, eptype, n, c.Type, c.Number)
			}
		})
	}
}

func TestEpisodeType_Prefix(t *testing.T) {
	t.Parallel()
	cases := []struct {
		Type   EpisodeType
		Prefix string
	}{
		{EpRegular, ""},
		{EpSpecial, "S"},
		{EpTrailer, "T"},
	}
	for _, c := range cases {
		t.Run(c.Type.String(), func(t *testing.T) {
			t.Parallel()
			got := c.Type.Prefix()
			if got != c.Prefix {
				t.Errorf("%s.Prefix() = %#v (expected %#v)", c.Type, got, c.Prefix)
			}
		})
	}
}
