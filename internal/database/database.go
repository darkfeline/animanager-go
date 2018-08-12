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

// Package database provides database functions for Animanager.
package database

import (
	"database/sql"
	"io/ioutil"
	"log"

	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"go.felesatra.moe/animanager/internal/migrate"
)

// Logger is used by this package for logging.
var Logger = log.New(ioutil.Discard, "database: ", log.LstdFlags)

// Open opens and returns the SQLite database.  The database is
// migrated to the newest version.
func Open(p string) (d *sql.DB, err error) {
	logSQLiteVersion()
	d, err = sql.Open("sqlite3", p)
	if err != nil {
		return nil, err
	}
	defer func(d *sql.DB) {
		if err != nil {
			d.Close()
		}
	}(d)
	if err := migrate.Migrate(d); err != nil {
		return nil, errors.Wrap(err, "migrate database")
	}
	return d, nil
}

// Open opens and returns a SQLite database from memory.  The database is
// migrated to the newest version.
//
// The database is between all connections, so it must be closed out
// between tests.
func OpenMem() (*sql.DB, error) {
	return Open("file::memory:?mode=memory&cache=shared")
}

func logSQLiteVersion() {
	v, vn, id := sqlite3.Version()
	Logger.Printf("SQLite version: %s %d %s", v, vn, id)
}
