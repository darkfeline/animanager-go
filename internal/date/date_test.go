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

package date

import "testing"

func TestString(t *testing.T) {
	t.Parallel()
	d := Date(978393600)
	s := d.String()
	exp := "2001-01-02"
	if s != exp {
		t.Errorf("Expected %#v, got %#v", exp, s)
	}
}

func TestNewString(t *testing.T) {
	t.Parallel()
	s := "2001-01-02"
	d, err := NewString(s)
	if err != nil {
		t.Fatalf("Error making date: %s", err)
	}
	var exp int64 = 978393600
	if int64(d) != exp {
		t.Errorf("Expected %#v, got %#v", exp, d)
	}
}
