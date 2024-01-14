// Copyright (C) 2024  Allen Li
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
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.felesatra.moe/animanager/internal/query"
)

func TestGetFileKey(t *testing.T) {
	t.Parallel()
	d := t.TempDir()
	p := filepath.Join(d, "testfile")
	if err := os.WriteFile(p, []byte("message digest"), 0600); err != nil {
		t.Fatal(err)
	}
	got, err := getFileKey(p)
	if err != nil {
		t.Fatal(err)
	}
	want := fileKey{
		Size: 14,
		Hash: query.Hash("d9130a8164549fe818874806e1c7014b"),
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("getFileKey() mismatch (-want +got):\n%s", diff)
	}
}
