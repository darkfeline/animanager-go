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
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/google/subcommands"

	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/query"
)

type SetDone struct {
	notDone bool
}

func (*SetDone) Name() string     { return "setdone" }
func (*SetDone) Synopsis() string { return "Set an episode's done status." }
func (*SetDone) Usage() string {
	return `Usage: add episodeID...
Set an episode's done status.
`
}

func (d *SetDone) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&d.notDone, "not", false, "Set status to not done")
}

func (d *SetDone) Execute(ctx context.Context, f *flag.FlagSet, x ...interface{}) subcommands.ExitStatus {
	if f.NArg() < 1 {
		fmt.Fprint(os.Stderr, d.Usage())
		return subcommands.ExitUsageError
	}
	ids := make([]int, f.NArg())
	for i, s := range f.Args() {
		id, err := strconv.Atoi(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid ID: %s\n", err)
			return subcommands.ExitUsageError
		}
		ids[i] = id
	}

	cfg := getConfig(x)
	db, err := database.Open(ctx, cfg.DBPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %s\n", err)
		return subcommands.ExitFailure
	}
	defer db.Close()
	for _, id := range ids {
		if err := query.UpdateEpisodeDone(db, id, !d.notDone); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating episode %d: %s\n", id, err)
			return subcommands.ExitFailure
		}
	}
	return subcommands.ExitSuccess
}
