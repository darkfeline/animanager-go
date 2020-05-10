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

package input

import (
	"bufio"
	"errors"
	"strings"
	"testing"
)

func TestReadYN(t *testing.T) {
	t.Parallel()
	t.Run("yes", func(t *testing.T) {
		t.Parallel()
		r := testReader("yes\n")
		got, err := ReadYN(r, false)
		if err != nil {
			t.Fatal(err)
		}
		want := true
		if got != want {
			t.Errorf("ReadYN() = %#v; want %#v", got, want)
		}
	})
	t.Run("no", func(t *testing.T) {
		t.Parallel()
		r := testReader("no\n")
		got, err := ReadYN(r, false)
		if err != nil {
			t.Fatal(err)
		}
		want := false
		if got != want {
			t.Errorf("ReadYN() = %#v; want %#v", got, want)
		}
	})
	t.Run("default", func(t *testing.T) {
		t.Parallel()
		r := testReader("\n")
		got, err := ReadYN(r, true)
		if err != nil {
			t.Fatal(err)
		}
		want := true
		if got != want {
			t.Errorf("ReadYN() = %#v; want %#v", got, want)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		t.Parallel()
		r := testReader("blah\n")
		_, err := ReadYN(r, true)
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, ErrInvalid) {
			t.Errorf("expected ErrInvalid")
		}
	})
}

func testReader(input string) Reader {
	return bufio.NewReader(strings.NewReader(input))
}
