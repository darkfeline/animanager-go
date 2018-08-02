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
	"database/sql"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/google/subcommands"
	"github.com/pkg/errors"
	"go.felesatra.moe/animanager/internal/config"
	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/models"
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
	if err := printAnime(os.Stdout, db, aid); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

func printAnime(w io.Writer, db *sql.DB, aid int) error {
	r, err := db.Query(`select aid, title, type, episodecount, startdate, enddate from anime where aid=?`, aid)
	if err != nil {
		return errors.Wrap(err, "failed to query anime")
	}
	defer r.Close()
	if !r.Next() {
		return r.Err()
	}
	a := models.Anime{}
	if err := r.Scan(&a.AID, &a.Title, &a.Type, &a.EpisodeCount, &a.StartDate, &a.EndDate); err != nil {
		return errors.Wrap(err, "failed to scan anime")
	}
	bw := bufio.NewWriter(w)
	fmt.Fprintf(bw, "AID: %d\n", a.AID)
	fmt.Fprintf(bw, "Title: %s\n", a.Title)
	fmt.Fprintf(bw, "Type: %s\n", a.Type)
	fmt.Fprintf(bw, "Episodes: %d\n", a.EpisodeCount)
	fmt.Fprintf(bw, "Start date: %s\n", a.StartDate)
	fmt.Fprintf(bw, "End date: %s\n", a.EndDate)
	return bw.Flush()
}
