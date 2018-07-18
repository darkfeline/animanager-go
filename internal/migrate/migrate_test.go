package migrate

import (
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
	if err := Migrate(d); err != nil {
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
