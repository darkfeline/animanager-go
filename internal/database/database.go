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
	"context"
	"database/sql"
	"io/ioutil"
	"log"
	"strings"

	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"go.felesatra.moe/animanager/internal/migrate"
)

// Logger is used by this package for logging.
var Logger = log.New(ioutil.Discard, "database: ", log.LstdFlags)

// Open opens and returns the SQLite database.  The database is
// migrated to the newest version.
func Open(ctx context.Context, src string) (db *sql.DB, err error) {
	db, err = openDB(ctx, src)
	if err != nil {
		return nil, err
	}
	defer func(db *sql.DB) {
		if err != nil {
			db.Close()
		}
	}(db)
	if err := migrate.Migrate(ctx, db); err != nil {
		return nil, errors.Wrap(err, "migrate database")
	}
	return db, nil
}

// OpenMem opens and returns a SQLite database from memory.  The
// database is migrated to the newest version.
//
// The database is shared between all concurrent connections, so it
// must be closed out between tests.
func OpenMem(ctx context.Context) (*sql.DB, error) {
	return Open(ctx, "file::memory:?mode=memory&cache=shared")
}

func openDB(ctx context.Context, src string) (*sql.DB, error) {
	logSQLiteVersion()
	src = addParam(src, "_fk", "1")
	db, err := sql.Open("sqlite3", src)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// addParam adds the parameter to the SQL data source.
func addParam(src, param, value string) string {
	var b strings.Builder
	b.WriteString(src)
	if strings.IndexByte(src, '?') == -1 {
		b.WriteString("?")
	} else {
		b.WriteString("&")
	}
	b.WriteString(param)
	b.WriteString("=")
	b.WriteString(value)
	return b.String()
}

func isMemorySource(src string) bool {
	return sourcePath(src) == ":memory:"
}

// sourcePath returns the path of the SQL data source string.
func sourcePath(src string) string {
	// Remove file: prefix if it exists.
	if strings.HasPrefix(src, "file:") {
		src = src[len("file:"):]
	}
	// Remove params
	if i := strings.IndexByte(src, '?'); i != -1 {
		src = src[:i]
	}
	return src
}

func logSQLiteVersion() {
	v, vn, id := sqlite3.Version()
	Logger.Printf("SQLite version: %s %d %s", v, vn, id)
}
