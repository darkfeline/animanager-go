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

// Package cmd implements subcommands.
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
	"go.felesatra.moe/animanager/internal/config"
	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/models"
	"go.felesatra.moe/animanager/internal/query"
)

type Show struct {
}

func (*Show) Name() string     { return "show" }
func (*Show) Synopsis() string { return "Show information about a series." }
func (*Show) Usage() string {
	return `show AID:
  Show information about a series.
`
}

func (s *Show) SetFlags(f *flag.FlagSet) {
}

func (s *Show) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s", s.Usage())
		return subcommands.ExitUsageError
	}
	aid, err := strconv.Atoi(f.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid AID: %s\n", err)
		return subcommands.ExitUsageError
	}
	c := config.New()
	db, err := database.Open(c.DBPath())
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
	printAnime(os.Stdout, a)
	es, err := query.GetEpisodes(db, aid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return subcommands.ExitFailure
	}
	for _, e := range es {
		printEpisode(os.Stdout, &e)
	}
	return subcommands.ExitSuccess
}

func printAnime(w io.Writer, a *models.Anime) error {
	bw := bufio.NewWriter(w)
	fmt.Fprintf(bw, "AID: %d\n", a.AID)
	fmt.Fprintf(bw, "Title: %s\n", a.Title)
	fmt.Fprintf(bw, "Type: %s\n", a.Type)
	fmt.Fprintf(bw, "Episodes: %d\n", a.EpisodeCount)
	fmt.Fprintf(bw, "Start date: %s\n", a.StartDate)
	fmt.Fprintf(bw, "End date: %s\n", a.EndDate)
	return bw.Flush()
}

func printEpisode(w io.Writer, e *models.Episode) error {
	bw := bufio.NewWriter(w)
	fmt.Fprintf(bw, "%d: ", e.ID)
	fmt.Fprintf(bw, "%d ", e.AID)
	fmt.Fprintf(bw, "%s ", e.Type)
	fmt.Fprintf(bw, "%d ", e.Number)
	fmt.Fprintf(bw, "%s ", e.Title)
	fmt.Fprintf(bw, "(%d min)", e.Length)
	if e.UserWatched {
		fmt.Fprintf(bw, " (done)")
	}
	fmt.Fprintf(bw, "\n")
	return bw.Flush()
}
