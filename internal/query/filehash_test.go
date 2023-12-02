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

package query

import (
	"reflect"
	"testing"

	"go.felesatra.moe/animanager/internal/database"
)

func TestInsertAndGetFileHash(t *testing.T) {
	db := database.OpenMem(t)
	fh := &FileHash{
		Size: 135,
		Hash: "shirasuazusa",
		EID:  555,
		AID:  444,
	}
	if err := InsertFileHash(db, fh); err != nil {
		t.Fatalf("Error inserting file hash: %s", err)
	}
	t.Run("get anime", func(t *testing.T) {
		got, err := GetFileHash(db, 135, "shirasuazusa")
		if err != nil {
			t.Fatalf("Error getting file hash: %s", err)
		}
		if !reflect.DeepEqual(got, fh) {
			t.Errorf("GetFileHash() = %#v; want %#v", got, fh)
		}
	})
}
