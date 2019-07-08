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

	"golang.org/x/xerrors"

	"go.felesatra.moe/animanager/internal/afmt"
	"go.felesatra.moe/animanager/internal/config"
	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/query"
)

type Show struct {
}

func (*Show) Name() string     { return "show" }
func (*Show) Synopsis() string { return "Show information about a series." }
func (*Show) Usage() string {
	return `Usage: show aid
Show information about a series.
`
}

func (*Show) SetFlags(f *flag.FlagSet) {
}

func (*Show) Run(ctx context.Context, f *flag.FlagSet, cfg config.Config) error {
	if f.NArg() != 1 {
		return usageError{"must pass exactly one argument"}
	}
	aid, err := strconv.Atoi(f.Arg(0))
	if err != nil {
		return fmt.Errorf("invalid AID %v: %v", aid, err)
	}
	db, err := database.Open(ctx, cfg.DBPath)
	if err != nil {
		return err
	}
	defer db.Close()
	a, err := query.GetAnime(db, aid)
	if err != nil {
		return err
	}
	bw := bufio.NewWriter(os.Stdout)
	afmt.PrintAnime(bw, a)
	w, err := query.GetWatching(db, aid)
	switch {
	case err == nil:
		fmt.Fprintf(bw, "Registered: %#v (offset %d)\n", w.Regexp, w.Offset)
	case xerrors.Is(err, sql.ErrNoRows):
		io.WriteString(bw, "Not registered\n")
	default:
		return err
	}
	es, err := query.GetEpisodes(db, aid)
	if err != nil {
		return err
	}
	for _, e := range es {
		afmt.PrintEpisode(bw, e)
	}
	bw.Flush()
	return nil
}
