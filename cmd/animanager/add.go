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
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"go.felesatra.moe/anidb"
	"go.felesatra.moe/animanager/cmd/animanager/vars"
	"go.felesatra.moe/animanager/internal/clientid"
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
		cfgv := vars.Config(f)
		addIncomplete := f.Bool("incomplete", false, "Re-add incomplete anime.")
		if err := f.Parse(args); err != nil {
			return err
		}

		if f.NArg() < 1 && !(*addIncomplete) {
			return errors.New("no AIDs given")
		}
		aids, err := query.ParseIDs[query.AID](f.Args())
		if err != nil {
			return err
		}

		db, err := cfgv.OpenDB()
		if err != nil {
			return err
		}
		defer db.Close()
		if *addIncomplete {
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
