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
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.felesatra.moe/anidb"
	"go.felesatra.moe/animanager/internal/config"
	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/query"
	"golang.org/x/time/rate"
)

var addCmd = command{
	usageLine: "add [-incomplete] [-all] [aids]",
	shortDesc: "add an anime",
	longDesc: `Add an anime.
`,
	run: func(cmd *command, cfg *config.Config, args []string) error {
		f := cmd.flagSet()
		addAll := f.Bool("all", false, "Re-add all anime (expensive).")
		addIncomplete := f.Bool("incomplete", false, "Re-add incomplete anime.")
		if err := f.Parse(args); err != nil {
			return err
		}

		if f.NArg() < 1 && !*addIncomplete {
			return errors.New("no AIDs given")
		}
		aids, err := parseIDs(f.Args())
		if err != nil {
			return err
		}

		db, err := openDB(cfg)
		if err != nil {
			return err
		}
		defer db.Close()
		if *addAll {
			as, err := query.GetAIDs(db)
			if err != nil {
				return err
			}
			aids = append(aids, as...)
		} else if *addIncomplete {
			as, err := query.GetIncompleteAnime(db)
			if err != nil {
				return err
			}
			aids = append(aids, as...)
		}
		for _, aid := range aids {
			fmt.Println(aid)
			if err := addAnime(db, aid); err != nil {
				return err
			}
		}
		return nil
	},
}

var client = &anidb.Client{
	Name:    "kfanimanager",
	Version: 2,
	Limiter: rate.NewLimiter(rate.Every(2*time.Second), 1),
}

func addAnime(db *sql.DB, aid int) error {
	log.Printf("Adding %d", aid)
	c, err := client.RequestAnime(aid)
	if err != nil {
		return fmt.Errorf("add anime %v: %w", aid, err)
	}
	if err := query.InsertAnime(db, c); err != nil {
		return fmt.Errorf("add anime %v: %w", aid, err)
	}
	return nil
}

func openDB(cfg *config.Config) (*sql.DB, error) {
	return database.Open(context.Background(), cfg.DBPath)
}

func parseIDs(args []string) ([]int, error) {
	ids := make([]int, len(args))
	for i, s := range args {
		id, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid ID %v: %v", s, err)
		}
		ids[i] = id
	}
	return ids, nil
}
