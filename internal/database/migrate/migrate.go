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
	"go.felesatra.moe/database/sql/sqlite3/migrate"
)

var migrationSet = migrate.NewMigrationSet([]migrate.Migration{
	{From: 0, To: 3, Func: migrate3},
	{From: 3, To: 4, Func: migrate4},
	{From: 4, To: 5, Func: migrate5},
	{From: 5, To: 6, Func: migrate6},
	{From: 6, To: 7, Func: migrate7},
	{From: 7, To: 8, Func: migrate8},
	{From: 8, To: 9, Func: migrate9},
	{From: 9, To: 10, Func: migrate10},
	{From: 10, To: 11, Func: migrate11},
	{From: 11, To: 12, Func: migrate12},
})

var (
	Migrate      = migrationSet.Migrate
	NeedsMigrate = migrationSet.NeedsMigrate
)
