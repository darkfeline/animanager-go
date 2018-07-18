package database

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestOpen(t *testing.T) {
	d, err := Open("file::memory:?mode=memory&cache=shared")
	defer d.Close()
	if err != nil {
		t.Errorf("Error opening database: %s", err)
	}
}
