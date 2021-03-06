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
	"fmt"
	"io"
	"os"
	"strconv"

	"go.felesatra.moe/animanager/internal/afmt"
	"go.felesatra.moe/animanager/internal/config"
	"go.felesatra.moe/animanager/internal/query"
)

var showFilesCmd = command{
	usageLine: "showfiles [-episode] [AIDs | episodeIDs]",
	shortDesc: "show episode files",
	longDesc: `Show episode files.
`,
	run: func(c *command, cfg *config.Config, args []string) error {
		f := c.flagSet()
		episode := f.Bool("episode", false, "Show files for episode.")
		if err := f.Parse(args); err != nil {
			return err
		}

		if f.NArg() != 1 {
			return errors.New("must pass exactly one argument")
		}
		id, err := strconv.Atoi(f.Arg(0))
		if err != nil {
			return fmt.Errorf("invalid ID %v: %v", id, err)
		}

		db, err := openDB(cfg)
		if err != nil {
			return err
		}
		defer db.Close()
		bw := bufio.NewWriter(os.Stdout)
		if *episode {
			err = showEpisodeFiles(bw, db, id)
		} else {
			err = showAnimeFiles(bw, db, id)
		}
		bw.Flush()
		return err
	},
}

func showAnimeFiles(w io.Writer, db *sql.DB, aid int) error {
	eps, err := query.GetEpisodes(db, aid)
	if err != nil {
		return err
	}
	for _, e := range eps {
		afmt.PrintEpisode(w, e)
		efs, err := query.GetEpisodeFiles(db, e.ID)
		if err != nil {
			return err
		}
		for _, ef := range efs {
			fmt.Fprintf(w, "\t\t  %s\n", ef.Path)
		}
	}
	return nil
}

func showEpisodeFiles(w io.Writer, db *sql.DB, id int) error {
	efs, err := query.GetEpisodeFiles(db, id)
	if err != nil {
		return err
	}
	for _, ef := range efs {
		fmt.Fprintf(w, "%s\n", ef.Path)
	}
	return nil
}
