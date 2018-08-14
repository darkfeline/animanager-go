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

	"github.com/google/subcommands"

	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/query"
)

type ShowFiles struct {
	anime bool
}

func (*ShowFiles) Name() string     { return "showfiles" }
func (*ShowFiles) Synopsis() string { return "Show episode files." }
func (*ShowFiles) Usage() string {
	return `Usage: findfiles [-anime] id|aid
Show episode files.
`
}

func (sf *ShowFiles) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&sf.anime, "anime", false, "Show files for anime")
}

func (sf *ShowFiles) Execute(ctx context.Context, f *flag.FlagSet, x ...interface{}) subcommands.ExitStatus {
	if f.NArg() != 1 {
		fmt.Fprint(os.Stderr, sf.Usage())
		return subcommands.ExitUsageError
	}
	id, err := strconv.Atoi(f.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid ID: %s\n", err)
		return subcommands.ExitUsageError
	}

	c := getConfig(x)
	db, err := database.Open(ctx, c.DBPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %s\n", err)
		return subcommands.ExitFailure
	}
	defer db.Close()
	bw := bufio.NewWriter(os.Stdout)
	if sf.anime {
		err = showAnimeFiles(bw, db, id)
	} else {
		err = showEpisodeFiles(bw, db, id)
	}
	bw.Flush()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

func showAnimeFiles(w io.Writer, db *sql.DB, aid int) error {
	eps, err := query.GetEpisodes(db, aid)
	if err != nil {
		return err
	}
	for _, e := range eps {
		printEpisode(w, &e)
		efs, err := query.GetEpisodeFiles(db, e.ID)
		if err != nil {
			return err
		}
		for _, ef := range efs {
			fmt.Fprintf(w, "\t\t%s\n", ef.Path)
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
