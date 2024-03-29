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
	"errors"
	"fmt"

	"go.felesatra.moe/animanager/cmd/animanager/vars"
	"go.felesatra.moe/animanager/internal/query"
	"go.felesatra.moe/animanager/internal/sqlc"
)

var unregisterCmd = command{
	usageLine: "unregister [-finished] [aids]",
	shortDesc: "unregister anime",
	longDesc: `Unregister anime.
`,
	run: func(h *handle, args []string) error {
		f := h.flagSet()
		cfgv := vars.Config(f)
		finished := f.Bool("finished", false, "Unregister finished anime.")
		if err := f.Parse(args); err != nil {
			return err
		}

		if f.NArg() < 1 && !*finished {
			return errors.New("no anime specified")
		}
		aids, err := sqlc.ParseIDs[sqlc.AID](f.Args())
		if err != nil {
			return err
		}

		db, err := cfgv.OpenDB()
		if err != nil {
			return err
		}
		defer db.Close()
		if *finished {
			fas, err := query.GetFinishedWatchingAIDs(db)
			if err != nil {
				return err
			}
			aids = append(aids, fas...)
		}
		ctx := context.Background()
		q := sqlc.New(db)
		if err := q.PrepareDeleteWatching(ctx); err != nil {
			return err
		}
		defer q.Close()
		for _, aid := range aids {
			fmt.Println(aid)
			if err := q.DeleteWatching(ctx, aid); err != nil {
				return err
			}
		}
		return nil
	},
}
