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
	defer db.Close()
	a, err := getAnime(db, aid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return subcommands.ExitFailure
	}
	printAnime(os.Stdout, a)
	es, err := getEpisodes(db, aid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return subcommands.ExitFailure
	}
	for _, e := range es {
		printEpisode(os.Stdout, &e)
	}
	return subcommands.ExitSuccess
}

func getAnime(db *sql.DB, aid int) (*models.Anime, error) {
	t, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "failed to open transaction")
	}
	defer t.Rollback()
	r, err := t.Query(`SELECT aid, title, type, episodecount, startdate, enddate FROM anime WHERE aid=?`, aid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query anime")
	}
	defer r.Close()
	if !r.Next() {
		return nil, r.Err()
	}
	a := models.Anime{}
	if err := r.Scan(&a.AID, &a.Title, &a.Type, &a.EpisodeCount, &a.StartDate, &a.EndDate); err != nil {
		return nil, errors.Wrap(err, "failed to scan anime")
	}
	return &a, nil
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

func getEpisodes(db *sql.DB, aid int) ([]models.Episode, error) {
	r, err := db.Query(`
SELECT id, aid, type, number, title, length, user_watched
FROM episode WHERE aid=? ORDER BY type, number`, aid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query episode")
	}
	defer r.Close()
	var es []models.Episode
	for r.Next() {
		e := models.Episode{}
		if err := r.Scan(&e.ID, &e.AID, &e.Type, &e.Number, &e.Title, &e.Length, &e.UserWatched); err != nil {
			return nil, errors.Wrap(err, "failed to scan episode")
		}
		es = append(es, e)
	}
	if err := r.Err(); err != nil {
		return nil, err
	}
	return es, nil
}

func printEpisode(w io.Writer, e *models.Episode) error {
	bw := bufio.NewWriter(w)
	fmt.Fprintf(bw, "ID: %d\n", e.ID)
	fmt.Fprintf(bw, "AID: %d\n", e.AID)
	fmt.Fprintf(bw, "Type: %s\n", e.Type)
	fmt.Fprintf(bw, "Number: %d\n", e.Number)
	fmt.Fprintf(bw, "Title: %s\n", e.Title)
	fmt.Fprintf(bw, "Length: %d\n", e.Length)
	fmt.Fprintf(bw, "Watched: %t\n", e.UserWatched)
	return bw.Flush()
}
