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
	"errors"
	"fmt"

	"go.felesatra.moe/animanager/internal/config"
	"go.felesatra.moe/animanager/internal/query"
)

var unregisterCmd = command{
	usageLine: "unregister [-watched] [aids]",
	shortDesc: "unregister anime",
	longDesc:  "Unregister anime.",
	run: func(c *command, cfg *config.Config, args []string) error {
		f := c.flagSet()
		watched := f.Bool("watched", false, "Unregister watched anime.")
		if err := f.Parse(args); err != nil {
			return err
		}

		if f.NArg() < 1 && !*watched {
			return errors.New("no anime specified")
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
		if *watched {
			watching, err := query.GetAllWatching(db)
			if err != nil {
				return err
			}
			watchingMap := make(map[int]bool)
			for _, w := range watching {
				watchingMap[w.AID] = true
			}

			watched, err := query.GetFinishedAnime(db)
			if err != nil {
				return err
			}
			for _, a := range watched {
				if watchingMap[a.AID] {
					aids = append(aids, a.AID)
				}
			}
		}
		for _, aid := range aids {
			fmt.Println(aid)
			if err := query.DeleteWatching(db, aid); err != nil {
				return err
			}
		}
		return nil
	},
}
