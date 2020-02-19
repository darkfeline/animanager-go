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

// Package migrate implements migrations for the Animanager SQLite
// database.
package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// Migrate migrates the database to the newest version.
func Migrate(ctx context.Context, d *sql.DB) error {
	v, err := getUserVersion(d)
	if err != nil {
		return fmt.Errorf("get user version: %w", err)
	}
	for _, m := range migrations {
		if v != m.From {
			continue
		}
		log.Printf("Migrating database from %d to %d", m.From, m.To)
		if err := m.Func(ctx, d); err != nil {
			return fmt.Errorf("migrate from %d to %d: %w", m.From, m.To, err)
		}
		if err := setUserVersion(d, m.To); err != nil {
			return err
		}
		v = m.To
	}
	return nil
}

// IsCurrentVersion returns true if the database is the current
// version.
func IsCurrentVersion(d *sql.DB) (bool, error) {
	v, err := getUserVersion(d)
	if err != nil {
		return false, fmt.Errorf("get user version: %w", err)
	}
	return v == currentVersion(), nil
}

func currentVersion() int {
	return migrations[len(migrations)-1].To
}

type spec struct {
	From int
	To   int
	Func migrateFunc
}

var migrations = []spec{
	{0, 3, migrate3},
	{3, 4, migrate4},
	{4, 5, migrate5},
}

type migrateFunc func(context.Context, *sql.DB) error
