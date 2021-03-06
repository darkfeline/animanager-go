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

import (
	"testing"
	"time"
)

func TestString(t *testing.T) {
	t.Parallel()
	d := Date(978393600)
	s := d.String()
	exp := "2001-01-02"
	if s != exp {
		t.Errorf("Expected %#v, got %#v", exp, s)
	}
}

func TestParse(t *testing.T) {
	t.Parallel()
	s := "2001-01-02"
	d, err := Parse(s)
	if err != nil {
		t.Fatalf("Error making date: %s", err)
	}
	var exp int64 = 978393600
	if int64(d) != exp {
		t.Errorf("Expected %#v, got %#v", exp, d)
	}
}

func TestParse_invalid(t *testing.T) {
	t.Parallel()
	s := "foobar"
	_, err := Parse(s)
	if err == nil {
		t.Errorf("Got no error")
	}
}

func TestZero(t *testing.T) {
	t.Parallel()
	if Zero.String() != "0000-01-01" {
		t.Errorf("Zero is not 0000-01-01")
	}
}

func TestFromTime(t *testing.T) {
	t.Parallel()
	input := time.Date(2000, 12, 31, 23, 59, 59, 1_000_000, time.FixedZone("CET", 60*60))
	got := FromTime(input)
	want := Date(time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC).Unix())
	if got != want {
		t.Errorf("FromTime(%v) = %v; want %v", input, got, want)
	}
}
