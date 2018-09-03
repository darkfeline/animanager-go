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
	"os"
	"os/exec"
	"strconv"

	"github.com/google/subcommands"
	"github.com/pkg/errors"

	"go.felesatra.moe/animanager/internal/config"
	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/input"
	"go.felesatra.moe/animanager/internal/query"
)

type Watch struct {
	episode bool
}

func (*Watch) Name() string     { return "watch" }
func (*Watch) Synopsis() string { return "Watch anime." }
func (*Watch) Usage() string {
	return `Usage: watch [-episode] AID|episodeID
Watch anime.
`
}

func (w *Watch) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&w.episode, "episode", false, "Treat argument as episode ID")
}

func (w *Watch) Execute(ctx context.Context, f *flag.FlagSet, x ...interface{}) subcommands.ExitStatus {
	if f.NArg() != 1 {
		fmt.Fprint(os.Stderr, w.Usage())
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
	if w.episode {
		err = watchEpisode(c, db, id)
	} else {
		err = watchAnime(c, db, id)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

func watchEpisode(c config.Config, db *sql.DB, id int) error {
	e, err := query.GetEpisode(db, id)
	if err != nil {
		return errors.Wrap(err, "get episode")
	}
	printEpisode(os.Stdout, *e)
	fs, err := query.GetEpisodeFiles(db, id)
	if err != nil {
		return err
	}
	if len(fs) == 0 {
		return fmt.Errorf("no files for episode %d", id)
	}
	f := fs[0]
	fmt.Println(f.Path)
	if err := playFile(c, f.Path); err != nil {
		return err
	}
	if e.UserWatched {
		fmt.Println("Already watched")
		return nil
	}
	br := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Set done? [Y/n] ")
		ans, err := input.ReadYN(br, true)
		if err, ok := err.(temporary); ok && err.Temporary() {
			fmt.Println(err)
			continue
		}
		if err != nil {
			return err
		}
		if !ans {
			return nil
		}
		if err := query.UpdateEpisodeDone(db, id, true); err != nil {
			return err
		}
		return nil
	}
}

type temporary interface {
	Temporary() bool
}

func playFile(c config.Config, p string) error {
	cmd := exec.Command(c.Player[0], append(c.Player[1:], p)...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func watchAnime(c config.Config, db *sql.DB, aid int) error {
	a, err := query.GetAnime(db, aid)
	if err != nil {
		return errors.Wrap(err, "get anime")
	}
	printAnimeShort(os.Stdout, a)
	eps, err := query.GetEpisodes(db, aid)
	if err != nil {
		return err
	}
	for _, e := range eps {
		if e.UserWatched {
			continue
		}
		return watchEpisode(c, db, e.ID)
	}
	return errors.Errorf("no unwatched episodes for %d", aid)
}
