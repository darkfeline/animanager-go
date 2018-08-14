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

package cmd

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/google/subcommands"

	"go.felesatra.moe/animanager/internal/database"
)

type Watch struct {
}

func (*Watch) Name() string     { return "watch" }
func (*Watch) Synopsis() string { return "Watch anime." }
func (*Watch) Usage() string {
	return `Usage: watch [episodeID]
Watch anime.
`
}

func (*Watch) SetFlags(f *flag.FlagSet) {
}

func (w *Watch) Execute(ctx context.Context, f *flag.FlagSet, x ...interface{}) subcommands.ExitStatus {
	var id int
	switch f.NArg() {
	case 0:
	case 1:
		var err error
		id, err = strconv.Atoi(f.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid ID: %s\n", err)
			return subcommands.ExitUsageError
		}
	default:
		fmt.Fprint(os.Stderr, w.Usage())
		return subcommands.ExitUsageError
	}

	c := getConfig(x)
	db, err := database.Open(ctx, c.DBPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %s\n", err)
		return subcommands.ExitFailure
	}
	defer db.Close()
	if id == 0 {
		err = showWatchable(db)
	} else {
		err = watchEpisode(db, id)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

func showWatchable(db *sql.DB) error {
	return nil
}

func watchEpisode(db *sql.DB, id int) error {
	return nil
}
