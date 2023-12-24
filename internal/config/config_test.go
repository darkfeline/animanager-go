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

package config

import "testing"

func TestConfig_FileRegexps_empty(t *testing.T) {
	t.Parallel()
	c := Config{}
	got, err := c.FileRegexps()
	if err != nil {
		t.Fatal(err)
	}
	if n := len(got); n != 0 {
		t.Errorf("len(got) = %v; want %v", n, 0)
	}
}

func TestConfig_FileRegexps(t *testing.T) {
	t.Parallel()
	c := Config{
		FilePatterns: []string{
			"azusa",
		},
	}
	got, err := c.FileRegexps()
	if err != nil {
		t.Fatal(err)
	}
	if n := len(got); n != 1 {
		t.Errorf("len(got) = %v; want %v", n, 1)
	}
}
