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
	"io"
	"os"
	"strconv"

	"go.felesatra.moe/animanager/internal/afmt"
	"go.felesatra.moe/animanager/internal/config"
	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/query"
)

type ShowFiles struct {
	episode bool
}

func (*ShowFiles) Name() string     { return "showfiles" }
func (*ShowFiles) Synopsis() string { return "Show episode files." }
func (*ShowFiles) Usage() string {
	return `Usage: showfiles [-episode] AID|episodeID
Show episode files.
`
}

func (c *ShowFiles) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&c.episode, "episode", false, "Show files for episode")
}

func (c *ShowFiles) Run(ctx context.Context, f *flag.FlagSet, cfg config.Config) error {
	if f.NArg() != 1 {
		return usageError{"must pass exactly one argument"}
	}
	id, err := strconv.Atoi(f.Arg(0))
	if err != nil {
		return fmt.Errorf("invalid ID %v: %v", id, err)
	}

	db, err := database.Open(ctx, cfg.DBPath)
	if err != nil {
		return err
	}
	defer db.Close()
	bw := bufio.NewWriter(os.Stdout)
	if c.episode {
		err = showEpisodeFiles(bw, db, id)
	} else {
		err = showAnimeFiles(bw, db, id)
	}
	bw.Flush()
	return err
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
