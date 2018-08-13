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

package migrate

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestMigrate(t *testing.T) {
	d, err := sql.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
	defer d.Close()
	if err != nil {
		t.Fatalf("Error opening database: %s", err)
	}
	if err := Migrate(context.Background(), d); err != nil {
		t.Errorf("Error migrating database: %s", err)
	}
}

func TestUserVersion(t *testing.T) {
	d, err := sql.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
	defer d.Close()
	if err != nil {
		t.Fatalf("Error opening database: %s", err)
	}
	v, err := getUserVersion(d)
	if err != nil {
		t.Fatalf("Error getting version: %s", err)
	}
	if v != 0 {
		t.Errorf("Expected 0, got %d", v)
	}
	err = setUserVersion(d, 1)
	if err != nil {
		t.Fatalf("Error setting version: %s", err)
	}
	v, err = getUserVersion(d)
	if err != nil {
		t.Fatalf("Error getting version: %s", err)
	}
	if v != 1 {
		t.Errorf("Expected 1, got %d", v)
	}
}
