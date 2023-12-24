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
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFilterFiles_multiple_match(t *testing.T) {
	t.Parallel()
	rs := []*regexp.Regexp{
		regexp.MustCompile(`mika`),
		regexp.MustCompile(`cute`),
		regexp.MustCompile(`princess`),
	}
	p := []string{"mika-is-a-cute-princess.mp4"}
	got := FilterFiles(rs, p)
	want := p
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("FilterFiles() mismatch (-want +got):\n%s", diff)
	}
}
