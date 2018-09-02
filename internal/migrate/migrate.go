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
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
)

// Logger is used by this package for logging.
var Logger = log.New(ioutil.Discard, "migrate: ", log.LstdFlags)

// Migrate migrates the database to the newest version.
func Migrate(ctx context.Context, d *sql.DB) error {
	v, err := getUserVersion(d)
	if err != nil {
		return errors.Wrap(err, "get user version")
	}
	for _, m := range migrations {
		if v != m.From {
			continue
		}
		Logger.Printf("Migrating from %d to %d", m.From, m.To)
		if err := m.Func(ctx, d); err != nil {
			return errors.Wrapf(err, "migrate from %d to %d", m.From, m.To)
		}
		if err := setUserVersion(d, m.To); err != nil {
			return err
		}
		v = m.To
	}
	return nil
}

var migrations = []struct {
	From int
	To   int
	Func migrateFunc
}{
	{0, 3, migrate3},
	{3, 4, migrate4},
	{4, 5, migrate5},
}

type migrateFunc func(context.Context, *sql.DB) error
