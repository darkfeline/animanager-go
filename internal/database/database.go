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
	"fmt"
	"io"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"go.felesatra.moe/animanager/internal/database/migrate"
)

// Open opens and returns the SQLite database.  The database is
// migrated to the newest version.
func Open(ctx context.Context, dataSrc string) (db *sql.DB, err error) {
	dataSrc = addParam(dataSrc, "_fk", "1")
	db, err = sql.Open("sqlite3", dataSrc)
	if err != nil {
		return nil, fmt.Errorf("open database %v: %w", dataSrc, err)
	}
	defer func(db *sql.DB) {
		if err != nil {
			db.Close()
		}
	}(db)
	current, err := migrate.IsCurrentVersion(db)
	if err != nil {
		return nil, fmt.Errorf("open database %v: %w", dataSrc, err)
	}
	if current {
		return db, nil
	}
	if !isMemorySource(dataSrc) {
		if err := backup(ctx, db, sourcePath(dataSrc)); err != nil {
			return nil, fmt.Errorf("open database %v: backup: %w", dataSrc, err)
		}
	}
	if err := migrate.Migrate(ctx, db); err != nil {
		return nil, fmt.Errorf("open database %v: %w", dataSrc, err)
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

// withLock calls the provided function with a write transaction lock
// on the database.
func withLock(ctx context.Context, db *sql.DB, f func() error) error {
	c, err := db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("with lock: %w", err)
	}
	defer c.Close()
	if _, err = c.ExecContext(ctx, "BEGIN IMMEDIATE"); err != nil {
		return fmt.Errorf("backup: %w", err)
	}
	defer c.ExecContext(ctx, "ROLLBACK")
	return f()
}

// backup the database file.
func backup(ctx context.Context, db *sql.DB, src string) error {
	dst := src + ".bak"
	f := func() error {
		return copyFile(src, dst)
	}
	if err := withLock(ctx, db, f); err != nil {
		return fmt.Errorf("backup: %w", err)
	}
	return nil
}

func copyFile(src, dst string) error {
	sf, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("copy file %s to %s: %w", src, dst, err)
	}
	defer sf.Close()
	df, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return fmt.Errorf("copy file %s to %s: %w", src, dst, err)
	}
	defer df.Close()
	if _, err := io.Copy(df, sf); err != nil {
		return fmt.Errorf("copy file %s to %s: %w", src, dst, err)
	}

	if err := df.Close(); err != nil {
		return fmt.Errorf("copy file %s to %s: %w", src, dst, err)
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
