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
	"os"
	"sort"

	"go.felesatra.moe/animanager/internal/afmt"
	"go.felesatra.moe/animanager/internal/config"
	"go.felesatra.moe/animanager/internal/query"
)

var unfinishedCmd = command{
	usageLine: "unfinished",
	shortDesc: "print unfinished anime",
	longDesc: `Print unfinished anime.
`,
	run: func(cmd *command, cfg *config.Config, args []string) error {
		f := cmd.flagSet()
		if err := f.Parse(args); err != nil {
			return err
		}

		db, err := openDB(cfg)
		if err != nil {
			return err
		}
		defer db.Close()
		as, err := query.GetUnfinishedAnime(db)
		bw := bufio.NewWriter(os.Stdout)
		sort.Slice(as, func(i, j int) bool { return as[i].AID < as[j].AID })
		for _, a := range as {
			afmt.PrintAnimeShort(bw, a)
		}
		bw.Flush()
		return nil
	},
}
