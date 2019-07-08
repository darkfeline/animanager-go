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
	"os"

	"go.felesatra.moe/animanager/internal/afmt"
	"go.felesatra.moe/animanager/internal/config"
	"go.felesatra.moe/animanager/internal/database"
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

func (c *Watchable) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&c.all, "all", false, "Show all files")
	f.BoolVar(&c.missing, "missing", false, "Show next episodes missing files")
}

func (c *Watchable) Run(ctx context.Context, f *flag.FlagSet, cfg config.Config) error {
	if f.NArg() != 0 {
		return usageError{"no arguments allowed"}
	}

	db, err := database.Open(ctx, cfg.DBPath)
	if err != nil {
		return err
	}
	defer db.Close()
	o := afmt.PrintWatchableOption{
		IncludeWatched:      c.all,
		IncludeMissingFiles: c.missing,
	}
	if c.all {
		o.NumWatchable = -1
	}
	if err := showWatchable(db, o); err != nil {
		return err
	}
	return nil
}

func showWatchable(db *sql.DB, o afmt.PrintWatchableOption) error {
	bw := bufio.NewWriter(os.Stdout)
	defer bw.Flush()
	ws, err := query.GetAllWatching(db)
	if err != nil {
		return err
	}
	for _, c := range ws {
		if err := showWatchableSingle(db, bw, c.AID, o); err != nil {
			return err
		}
	}
	return nil
}

func showWatchableSingle(db *sql.DB, bw *bufio.Writer, aid int, o afmt.PrintWatchableOption) error {
	a, err := query.GetAnime(db, aid)
	if err != nil {
		return err
	}
	efs, err := query.GetAnimeFiles(db, aid)
	if err != nil {
		return err
	}
	return afmt.PrintWatchable(os.Stdout, a, efs, o)
}
