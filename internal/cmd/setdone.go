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
	"flag"
	"fmt"
	"strconv"

	"go.felesatra.moe/animanager/internal/config"
	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/query"
)

type SetDone struct {
	notDone bool
}

func (*SetDone) Name() string     { return "setdone" }
func (*SetDone) Synopsis() string { return "Set an episode's done status." }
func (*SetDone) Usage() string {
	return `Usage: add episodeID...
Set an episode's done status.
`
}

func (c *SetDone) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&c.notDone, "not", false, "Set status to not done")
}

func (c *SetDone) Run(ctx context.Context, f *flag.FlagSet, cfg config.Config) error {
	if f.NArg() < 1 {
		return usageError{"no arguments provided"}
	}
	ids := make([]int, f.NArg())
	for i, s := range f.Args() {
		id, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("invalid ID %v: %v", id, err)
		}
		ids[i] = id
	}

	db, err := database.Open(ctx, cfg.DBPath)
	if err != nil {
		return err
	}
	defer db.Close()
	for _, id := range ids {
		if err := query.UpdateEpisodeDone(db, id, !c.notDone); err != nil {
			return err
		}
	}
	return nil
}
