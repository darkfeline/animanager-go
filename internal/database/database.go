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
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	sqlite3 "github.com/mattn/go-sqlite3"
	"golang.org/x/xerrors"

	"go.felesatra.moe/animanager/internal/migrate"
)

// Logger is used by this package for logging.
var Logger = log.New(ioutil.Discard, "database: ", log.LstdFlags)

// Open opens and returns the SQLite database.  The database is
// migrated to the newest version.
func Open(ctx context.Context, src string) (db *sql.DB, err error) {
	logSQLiteVersion()
	src = addParam(src, "_fk", "1")
	db, err = sql.Open("sqlite3", src)
	if err != nil {
		return nil, xerrors.Errorf("open database %v: %w", src, err)
	}
	defer func(db *sql.DB) {
		if err != nil {
			db.Close()
		}
	}(db)
	if !isMemorySource(src) {
		ok, err := migrate.IsCurrentVersion(db)
		if err != nil {
			return nil, xerrors.Errorf("open database %v: %w", src, err)
		}
		if !ok {
			sp := sourcePath(src)
			dp := sourcePath(src) + ".bak"
			if err := backup(ctx, db, sp, dp); err != nil {
				return nil, xerrors.Errorf("open database %v: backup: %w", src, err)
			}
		}
	}
	if err := migrate.Migrate(ctx, db); err != nil {
		return nil, xerrors.Errorf("open database %v: %w", src, err)
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

func backup(ctx context.Context, db *sql.DB, src, dst string) error {
	c, err := db.Conn(ctx)
	if err != nil {
		return xerrors.Errorf("backup: %w", err)
	}
	defer c.Close()
	_, err = c.ExecContext(ctx, "BEGIN IMMEDIATE")
	if err != nil {
		return xerrors.Errorf("backup: %w", err)
	}
	defer c.ExecContext(ctx, "ROLLBACK")
	err = copyFile(src, dst)
	if err != nil {
		return xerrors.Errorf("backup: %w", err)
	}
	return nil
}

func copyFile(src, dst string) error {
	sf, err := os.Open(src)
	if err != nil {
		return xerrors.Errorf("copy file %s to %s: %w", src, dst, err)
	}
	defer sf.Close()
	df, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return xerrors.Errorf("copy file %s to %s: %w", src, dst, err)
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	if err != nil {
		return xerrors.Errorf("copy file %s to %s: %w", src, dst, err)
	}
	err = df.Close()
	if err != nil {
		return xerrors.Errorf("copy file %s to %s: %w", src, dst, err)
	}
	return nil
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
