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

	"github.com/google/subcommands"

	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/query"
)

type Watchable struct {
	all bool
}

func (*Watchable) Name() string     { return "watchable" }
func (*Watchable) Synopsis() string { return "Show watchable anime." }
func (*Watchable) Usage() string {
	return `Usage: watchable [-all]
Show watchable anime.
`
}

func (w *Watchable) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&w.all, "all", false, "Show all files")
}

func (w *Watchable) Execute(ctx context.Context, f *flag.FlagSet, x ...interface{}) subcommands.ExitStatus {
	if f.NArg() != 0 {
		fmt.Fprint(os.Stderr, w.Usage())
		return subcommands.ExitUsageError
	}

	c := getConfig(x)
	db, err := database.Open(ctx, c.DBPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %s\n", err)
		return subcommands.ExitFailure
	}
	defer db.Close()
	if err := showWatchable(db, w.all); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

const watchableEpsPrintLimit = 3

func showWatchable(db *sql.DB, all bool) error {
	bw := bufio.NewWriter(os.Stdout)
	defer bw.Flush()
	ws, err := query.GetAllWatching(db)
	if err != nil {
		return err
	}
	for _, w := range ws {
		var printed int
		afs, err := getAnimeFiles(db, w.AID)
		if err != nil {
			return err
		}
		for i, e := range afs.Episodes {
			// Skip if done.
			if e.UserWatched && !all {
				continue
			}
			fs := afs.Files[i]
			// Skip if no files.
			if len(fs) == 0 {
				continue
			}
			// If we have already printed enough episodes,
			// stop looping and just print that there are
			// more.
			if !all && printed >= watchableEpsPrintLimit {
				fmt.Fprint(bw, "MORE\t...\n")
				break
			}
			// Print anime and previous episode if we are
			// printing the first episode for an anime.
			if printed == 0 {
				a, err := query.GetAnime(db, w.AID)
				if err != nil {
					return err
				}
				printAnimeShort(bw, a)
				if i > 0 {
					e := afs.Episodes[i-1]
					printEpisode(bw, e)
				}
			}
			printEpisode(bw, e)
			printed++
			for _, f := range fs {
				fmt.Fprintf(bw, "\t\t  %s\n", f.Path)
			}
		}
		if printed > 0 {
			fmt.Fprintln(bw)
		}
	}
	return nil
}

func getAnimeFiles(db *sql.DB, aid int) (afs animeFiles, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("get anime %d files: %s", aid, err)
		}
	}()
	eps, err := query.GetEpisodes(db, aid)
	if err != nil {
		return afs, err
	}
	for _, e := range eps {
		afs.Episodes = append(afs.Episodes, e)
		fs, err := query.GetEpisodeFiles(db, e.ID)
		if err != nil {
			return afs, err
		}
		afs.Files = append(afs.Files, fs)
	}
	return afs, nil
}

type animeFiles struct {
	Episodes []query.Episode
	Files    [][]query.EpisodeFile
}

func printAnimeShort(w io.Writer, a *query.Anime) {
	fmt.Fprintf(w, "%d\t%s\t%d eps\n", a.AID, a.Title, a.EpisodeCount)
}