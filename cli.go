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

	"github.com/google/subcommands"
)

type cliCmd struct {
	capitalize bool
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

func (c *cliCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	return subcommands.ExitSuccess

	d, err := sql.Open("sqlite3", p)
	if err != nil {

	}
}
