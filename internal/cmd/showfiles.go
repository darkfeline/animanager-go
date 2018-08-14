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
	"flag"
	"fmt"
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
	var efs []query.EpisodeFile
	if sf.anime {
		efs, err = query.GetAnimeFiles(db, id)
	} else {
		efs, err = query.GetEpisodeFiles(db, id)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return subcommands.ExitFailure
	}
	bw := bufio.NewWriter(os.Stdout)
	if sf.anime {

	} else {
		for _, ef := range efs {
			fmt.Fprintf(bw, "%s\n", ef.Path)
		}
	}
	bw.Flush()
	return subcommands.ExitSuccess
}
