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
	"bufio"
	"database/sql"
	"errors"
	"os"

	"go.felesatra.moe/animanager/internal/afmt"
	"go.felesatra.moe/animanager/internal/query"
)

var watchableCmd = command{
	usageLine: "watchable [-all] [-missing]",
	shortDesc: "show watchable anime",
	longDesc: `Show watchable anime.
`,
	run: func(cmd *command, args []string) error {
		f := cmd.flagSet()
		all := f.Bool("all", false, "Show all files.")
		missing := f.Bool("missing", false, "Show next episodes missing files.")
		if err := f.Parse(args); err != nil {
			return err
		}
		cfg, err := cmd.loadConfig()
		if err != nil {
			return err
		}

		if f.NArg() != 0 {
			return errors.New("no arguments allowed")
		}

		db, err := openDB(cfg)
		if err != nil {
			return err
		}
		defer db.Close()
		o := afmt.PrintWatchableOption{
			IncludeWatched:      *all,
			IncludeMissingFiles: *missing,
		}
		if *all {
			o.NumWatchable = -1
		}
		if err := showWatchable(db, o); err != nil {
			return err
		}
		return nil
	},
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
