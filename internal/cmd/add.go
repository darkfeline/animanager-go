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
	"log"
	"strconv"
	"time"

	"go.felesatra.moe/animanager/internal/anidb"
	"go.felesatra.moe/animanager/internal/config"
	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/query"
	"golang.org/x/xerrors"
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

func (c *Add) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&c.addIncomplete, "incomplete", false, "Re-add incomplete anime")
}

func (c *Add) Run(ctx context.Context, f *flag.FlagSet, cfg config.Config) error {
	// Process arguments.
	if f.NArg() < 1 && !c.addIncomplete {
		return usageError{"no AIDs given"}
	}
	aids := make([]int, f.NArg())
	for i, s := range f.Args() {
		aid, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("invalid AID %v: %v", aid, err)
		}
		aids[i] = aid
	}

	db, err := database.Open(ctx, cfg.DBPath)
	if err != nil {
		return err
	}
	defer db.Close()
	if c.addIncomplete {
		as, err := query.GetIncompleteAnime(db)
		if err != nil {
			return err
		}
		aids = append(aids, as...)
	}
	for i, aid := range aids {
		fmt.Println(aid)
		if err := addAnime(db, aid); err != nil {
			return err
		}
		if i < len(aids)-1 {
			time.Sleep(2 * time.Second)
		}
	}
	return nil
}

func addAnime(db *sql.DB, aid int) error {
	log.Printf("Adding %d", aid)
	c, err := anidb.RequestAnime(aid)
	if err != nil {
		return xerrors.Errorf("add anime %v: %w", aid, err)
	}
	if err := query.InsertAnime(db, c); err != nil {
		return xerrors.Errorf("add anime %v: %w", aid, err)
	}
	return nil
}
