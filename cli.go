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

package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"path/filepath"

	"github.com/google/subcommands"
	_ "github.com/mattn/go-sqlite3"
	"go.felesatra.moe/animanager/internal/migrate"
)

type cliCmd struct {
}

func (*cliCmd) Name() string     { return "cli" }
func (*cliCmd) Synopsis() string { return "Start CLI." }
func (*cliCmd) Usage() string {
	return `cli:
  Start CLI.
`
}

func (c *cliCmd) SetFlags(f *flag.FlagSet) {
}

var (
	home = os.Getenv("HOME")
	// dbPath     = filepath.Join(home, ".animanager/database.db")
	dbPath     = filepath.Join(home, ".animanager/tmp.db")
	anidbCache = filepath.Join(home, ".animanager/anidb")
	watchDir   = filepath.Join(home, "anime")
)

const player = "mpv"

func (c *cliCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	d, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		eprintf("Error opening database: %s\n", err)
		return subcommands.ExitFailure
	}
	if err := migrate.Migrate(d); err != nil {
		eprintf("Error migrating database: %s\n", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
