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
	"go.felesatra.moe/animanager/internal/clientid"
	"go.felesatra.moe/animanager/internal/config"
	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/query"
	"golang.org/x/time/rate"
)

var addCmd = command{
	usageLine: "add [-incomplete] [-no-eid] [aids]",
	shortDesc: "add an anime",
	longDesc: `Add an anime.
`,
	run: func(cmd *command, args []string) error {
		f := cmd.flagSet()
		addNoEID := f.Bool("no-eid", false, "Add anime missing EIDs.")
		addIncomplete := f.Bool("incomplete", false, "Re-add incomplete anime.")
		if err := f.Parse(args); err != nil {
			return err
		}
		cfg, err := cmd.loadConfig()
		if err != nil {
			return err
		}

		if f.NArg() < 1 && !(*addIncomplete || *addNoEID) {
			return errors.New("no AIDs given")
		}
		aids, err := parseIDs[query.AID](f.Args())
		if err != nil {
			return err
		}

		db, err := openDB(cfg)
		if err != nil {
			return err
		}
		defer db.Close()
		if *addNoEID {
			as, err := query.GetAIDsMissingEIDs(db)
			if err != nil {
				return err
			}
			log.Printf("%d anime missing EIDs", len(as))
			// Limit entries to not get banned.
			if len(as) > 50 {
				as = as[:50]
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
			if err := addAnime(db, aid); err != nil {
				return err
			}
		}
		return nil
	},
}

var client = &anidb.Client{
	Name:    clientid.HTTPName,
	Version: clientid.HTTPVersion,
	Limiter: rate.NewLimiter(rate.Every(2*time.Second), 1),
}

func addAnime(db *sql.DB, aid query.AID) error {
	log.Printf("Adding %d", aid)
	c, err := client.RequestAnime(int(aid))
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

func parseIDs[T ~int](args []string) ([]T, error) {
	ids := make([]T, len(args))
	for i, s := range args {
		id, err := parseID[T](s)
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}
	return ids, nil
}

func parseID[T ~int](s string) (T, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %q into %T: %s", s, T(0), err)
	}
	return T(id), nil
}
