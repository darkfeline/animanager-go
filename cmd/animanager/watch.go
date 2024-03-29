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

package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"go.felesatra.moe/animanager/cmd/animanager/vars"
	"go.felesatra.moe/animanager/internal/afmt"
	"go.felesatra.moe/animanager/internal/config"
	"go.felesatra.moe/animanager/internal/input"
	"go.felesatra.moe/animanager/internal/query"
	"go.felesatra.moe/animanager/internal/sqlc"
)

var watchCmd = command{
	usageLine: "watch [-episode] [aid | episodeID]",
	shortDesc: "watch anime",
	longDesc: `Watch anime.
`,
	run: func(h *handle, args []string) error {
		f := h.flagSet()
		cfgv := vars.Config(f)
		episode := f.Bool("episode", false, "Treat argument as episode ID")
		if err := f.Parse(args); err != nil {
			return err
		}
		cfg, err := cfgv.Load()
		if err != nil {
			return err
		}

		if f.NArg() != 1 {
			return errors.New("must pass exactly one argument")
		}
		id, err := sqlc.ParseID[int](f.Arg(0))
		if err != nil {
			return fmt.Errorf("invalid ID %v: %v", id, err)
		}

		db, err := cfgv.OpenDB()
		if err != nil {
			return err
		}
		defer db.Close()
		if *episode {
			err = watchEpisode(cfg, db, sqlc.EID(id))
		} else {
			err = watchAnime(cfg, db, sqlc.AID(id))
		}
		return err
	},
}

func watchEpisode(cfg *config.Config, db *sql.DB, eid sqlc.EID) error {
	e, err := query.GetEpisode(db, eid)
	if err != nil {
		return fmt.Errorf("get episode: %c", err)
	}
	afmt.PrintEpisode(os.Stdout, *e)
	fs, err := query.GetEpisodeFiles(db, eid)
	if err != nil {
		return err
	}
	if len(fs) == 0 {
		return fmt.Errorf("no files for episode %d", eid)
	}
	f := fs[0]
	fmt.Println(f.Path)
	if err := playFile(cfg, f.Path); err != nil {
		return err
	}
	if e.UserWatched {
		fmt.Println("Already watched")
		return nil
	}
	br := bufio.NewReader(os.Stdin)
readInput:
	for {
		fmt.Print("Set done? [Y/n] ")
		ans, err := input.ReadYN(br, true)
		if err != nil {
			if errors.Is(err, input.ErrInvalid) {
				fmt.Println(err)
				continue readInput
			}
			return err
		}
		if !ans {
			return nil
		}
		if err := query.UpdateEpisodeDone(db, eid, true); err != nil {
			return err
		}
		return nil
	}
}

func playFile(cfg *config.Config, p string) error {
	cmd := exec.Command(cfg.Player[0], append(cfg.Player[1:], p)...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func watchAnime(cfg *config.Config, db *sql.DB, aid sqlc.AID) error {
	a, err := query.GetAnime(db, aid)
	if err != nil {
		return fmt.Errorf("get anime: %c", err)
	}
	afmt.PrintAnimeShort(os.Stdout, a)
	eps, err := query.GetEpisodes(db, aid)
	if err != nil {
		return err
	}
	for _, e := range eps {
		if e.UserWatched {
			continue
		}
		return watchEpisode(cfg, db, e.EID)
	}
	return fmt.Errorf("no unwatched episodes for %d", aid)
}
