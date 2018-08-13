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
	"io"
	"os"
	"strconv"

	"github.com/google/subcommands"

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

func (s *Show) Execute(ctx context.Context, f *flag.FlagSet, x ...interface{}) subcommands.ExitStatus {
	if f.NArg() != 1 {
		fmt.Fprint(os.Stderr, s.Usage())
		return subcommands.ExitUsageError
	}
	aid, err := strconv.Atoi(f.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid AID: %s\n", err)
		return subcommands.ExitUsageError
	}
	c := getConfig(x)
	db, err := database.Open(ctx, c.DBPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %s\n", err)
		return subcommands.ExitFailure
	}
	defer db.Close()
	a, err := query.GetAnime(db, aid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return subcommands.ExitFailure
	}
	bw := bufio.NewWriter(os.Stdout)
	printAnime(bw, a)
	switch p, err := query.GetWatching(db, aid); err {
	case query.ErrMissing:
		io.WriteString(bw, "Not registered\n")
	case nil:
		fmt.Fprintf(bw, "Registered: %#v\n", p)
	default:
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return subcommands.ExitFailure
	}
	es, err := query.GetEpisodes(db, aid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return subcommands.ExitFailure
	}
	for _, e := range es {
		printEpisode(bw, &e)
	}
	bw.Flush()
	return subcommands.ExitSuccess
}

func printAnime(w io.Writer, a *query.Anime) {
	fmt.Fprintf(w, "AID: %d\n", a.AID)
	fmt.Fprintf(w, "Title: %s\n", a.Title)
	fmt.Fprintf(w, "Type: %s\n", a.Type)
	fmt.Fprintf(w, "Episodes: %d\n", a.EpisodeCount)
	fmt.Fprintf(w, "Start date: %s\n", a.StartDate)
	fmt.Fprintf(w, "End date: %s\n", a.EndDate)
}

func printEpisode(w io.Writer, e *query.Episode) {
	fmt.Fprintf(w, "%d: ", e.ID)
	fmt.Fprintf(w, "%s ", e.Type)
	fmt.Fprintf(w, "%d ", e.Number)
	fmt.Fprintf(w, "%s ", e.Title)
	fmt.Fprintf(w, "(%d min)", e.Length)
	if e.UserWatched {
		fmt.Fprintf(w, " (done)")
	}
	fmt.Fprintf(w, "\n")
}
