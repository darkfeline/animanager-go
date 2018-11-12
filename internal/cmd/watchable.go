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
	"bufio"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"

	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/obx"
	"go.felesatra.moe/animanager/internal/query"
)

type Watchable struct {
	all     bool
	missing bool
}

func (*Watchable) Name() string     { return "watchable" }
func (*Watchable) Synopsis() string { return "Show watchable anime." }
func (*Watchable) Usage() string {
	return `Usage: watchable [-all] [-missing]
Show watchable anime.
`
}

func (w *Watchable) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&w.all, "all", false, "Show all files")
	f.BoolVar(&w.missing, "missing", false, "Show next episodes missing files")
}

func (w *Watchable) Execute(ctx context.Context, f *flag.FlagSet, x ...interface{}) subcommands.ExitStatus {
	if f.NArg() != 0 {
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
	o := obx.PrintWatchableOption{
		IncludeWatched:      w.all,
		IncludeMissingFiles: w.missing,
	}
	if w.all {
		o.NumWatchable = -1
	}
	if err := showWatchable(db, o); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

func showWatchable(db *sql.DB, o obx.PrintWatchableOption) error {
	bw := bufio.NewWriter(os.Stdout)
	defer bw.Flush()
	ws, err := query.GetAllWatching(db)
	if err != nil {
		return err
	}
	for _, w := range ws {
		if err := showWatchableSingle(db, bw, w.AID, o); err != nil {
			return err
		}
	}
	return nil
}

func showWatchableSingle(db *sql.DB, bw *bufio.Writer, aid int, o obx.PrintWatchableOption) error {
	a, err := query.GetAnime(db, aid)
	if err != nil {
		return err
	}
	efs, err := obx.GetAnimeFiles(db, aid)
	if err != nil {
		return err
	}
	return obx.PrintWatchable(os.Stdout, a, efs, o)
}
