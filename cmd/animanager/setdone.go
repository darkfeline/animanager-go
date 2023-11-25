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

	"go.felesatra.moe/animanager/internal/query"
)

var setDoneCmd = command{
	usageLine: "setdone [episodeIDs]",
	shortDesc: "set an episode's done status",
	longDesc: `Set an episode's done status.
`,
	run: func(cmd *command, args []string) error {
		f := cmd.flagSet()
		notDone := f.Bool("not", false, "Set status to not done.")
		if err := f.Parse(args); err != nil {
			return err
		}
		cfg, err := cmd.loadConfig()
		if err != nil {
			return err
		}

		if f.NArg() < 1 {
			return errors.New("no arguments provided")
		}

		ids, err := parseIDs[query.EpID](f.Args())
		if err != nil {
			return err
		}

		db, err := openDB(cfg)
		if err != nil {
			return err
		}
		defer db.Close()
		for _, id := range ids {
			if err := query.UpdateEpisodeDone(db, id, !*notDone); err != nil {
				return err
			}
		}
		return nil
	},
}
