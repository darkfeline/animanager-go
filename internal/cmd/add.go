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
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/subcommands"

	"go.felesatra.moe/animanager/internal/anidb"
	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/query"
)

type Add struct {
	addIncomplete bool
}

func (*Add) Name() string     { return "add" }
func (*Add) Synopsis() string { return "Add an anime." }
func (*Add) Usage() string {
	return `Usage: add aids...
       add -incomplete [aids...]
Add an anime.
`
}

func (a *Add) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&a.addIncomplete, "incomplete", false, "Re-add incomplete anime")
}

func (a *Add) Execute(ctx context.Context, f *flag.FlagSet, x ...interface{}) subcommands.ExitStatus {
	// Process arguments.
	if f.NArg() < 1 && !a.addIncomplete {
		fmt.Fprint(os.Stderr, a.Usage())
		return subcommands.ExitUsageError
	}
	aids := make([]int, f.NArg())
	for i, s := range f.Args() {
		aid, err := strconv.Atoi(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid AID: %s\n", err)
			return subcommands.ExitUsageError
		}
		aids[i] = aid
	}

	c := getConfig(x)
	db, err := database.Open(ctx, c.DBPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %s\n", err)
		return subcommands.ExitFailure
	}
	defer db.Close()
	if a.addIncomplete {
		as, err := getIncompleteAnime(db)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return subcommands.ExitFailure
		}
		aids = append(aids, as...)
	}
	for i, aid := range aids {
		if err := addAnime(db, aid); err != nil {
			fmt.Fprintf(os.Stderr, "Error adding anime: %s\n", err)
			return subcommands.ExitFailure
		}
		if i < len(aids)-1 {
			time.Sleep(2 * time.Second)
		}
	}
	return subcommands.ExitSuccess
}

func getIncompleteAnime(db *sql.DB) ([]int, error) {
	aids, err := query.GetAIDs(db)
	if err != nil {
		return nil, fmt.Errorf("get incomplete anime: %s", err)
	}
	var r []int
	for _, aid := range aids {
		ok, err := isIncomplete(db, aid)
		if err != nil {
			return nil, fmt.Errorf("get incomplete anime: %s", err)
		}
		if ok {
			r = append(r, aid)
		}
	}
	return r, nil
}

func isIncomplete(db *sql.DB, aid int) (bool, error) {
	a, err := query.GetAnime(db, aid)
	if err != nil {
		return false, fmt.Errorf("check %d completion: %s", aid, err)
	}
	eps, err := query.GetEpisodes(db, aid)
	if err != nil {
		return false, fmt.Errorf("check %d completion: %s", aid, err)
	}
	var rEps []query.Episode
	var unnamed int
	for _, e := range eps {
		if e.Type != query.EpRegular {
			continue
		}
		rEps = append(rEps, e)
		if isUnnamed(e) {
			unnamed += 1
		} else {
			// Unnamed episodes followed by named episodes
			// are probably just missing the episode title
			// entirely, so don't count them.
			unnamed = 0
		}
	}
	if len(rEps) < a.EpisodeCount {
		return true, nil
	}
	// This is just a heuristic, some shows only have titles for
	// first/last episode.
	if unnamed > 0 && unnamed < a.EpisodeCount-2 {
		return true, nil
	}
	return false, nil
}

func isUnnamed(e query.Episode) bool {
	return len(e.Title) > 8 && e.Title[:8] == "Episode "
}

func addAnime(db *sql.DB, aid int) error {
	Logger.Printf("Adding %d", aid)
	a, err := anidb.RequestAnime(aid)
	if err != nil {
		return err
	}
	if err := query.InsertAnime(db, a); err != nil {
		return err
	}
	return nil
}
