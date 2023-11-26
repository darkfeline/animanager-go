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

	"go.felesatra.moe/animanager/cmd/animanager/vars"
	"go.felesatra.moe/animanager/internal/afmt"
	"go.felesatra.moe/animanager/internal/query"
)

var showFilesCmd = command{
	usageLine: "showfiles [-episode] [AIDs | episodeIDs]",
	shortDesc: "show episode files",
	longDesc: `Show episode files.
`,
	run: func(cmd *command, args []string) error {
		f := cmd.flagSet()
		cfgv := vars.Config(f)
		episode := f.Bool("episode", false, "Show files for episode.")
		if err := f.Parse(args); err != nil {
			return err
		}

		if f.NArg() != 1 {
			return errors.New("must pass exactly one argument")
		}
		id, err := parseID[int](f.Arg(0))
		if err != nil {
			return fmt.Errorf("invalid ID %v: %v", id, err)
		}

		db, err := cfgv.OpenDB()
		if err != nil {
			return err
		}
		defer db.Close()
		bw := bufio.NewWriter(os.Stdout)
		if *episode {
			err = showEpisodeFiles(bw, db, query.EpID(id))
		} else {
			err = showAnimeFiles(bw, db, query.AID(id))
		}
		bw.Flush()
		return err
	},
}

func showAnimeFiles(w io.Writer, db *sql.DB, aid query.AID) error {
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

func showEpisodeFiles(w io.Writer, db *sql.DB, id query.EpID) error {
	efs, err := query.GetEpisodeFiles(db, id)
	if err != nil {
		return err
	}
	for _, ef := range efs {
		fmt.Fprintf(w, "%s\n", ef.Path)
	}
	return nil
}
