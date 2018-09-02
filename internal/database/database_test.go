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

package database

import (
	"context"
	"testing"
)

func TestOpenMem(t *testing.T) {
	d, err := OpenMem(context.Background())
	defer d.Close()
	if err != nil {
		t.Errorf("Error opening database: %s", err)
	}
}

func TestGetPath(t *testing.T) {
	t.Parallel()
	cases := []struct {
		src  string
		want string
	}{
		{"some/path", "some/path"},
		{"file:some/path", "some/path"},
		{"some/path?arg=true", "some/path"},
		{"file:some/path?arg=true", "some/path"},
	}
	for _, c := range cases {
		t.Run(c.src, func(t *testing.T) {
			t.Parallel()
			got := getPath(c.src)
			if got != c.want {
				t.Errorf("getPath(%#v) = %#v; want %#v", c.src, got, c.want)
			}
		})
	}
}
