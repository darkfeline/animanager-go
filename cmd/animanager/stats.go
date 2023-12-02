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
	"fmt"

	"go.felesatra.moe/animanager/cmd/animanager/vars"
	"go.felesatra.moe/animanager/internal/query"
)

var statsCmd = command{
	usageLine: "stats",
	shortDesc: "print various stats",
	longDesc: `Print various stats.
`,
	run: func(h *handle, args []string) error {
		f := h.flagSet()
		cfgv := vars.Config(f)
		if err := f.Parse(args); err != nil {
			return err
		}

		db, err := cfgv.OpenDB()
		if err != nil {
			return err
		}
		defer db.Close()
		n, err := query.GetAnimeCount(db)
		if err != nil {
			return err
		}
		fmt.Printf("Total anime:\t%d\n", n)
		n, err = query.GetEpisodeCount(db)
		if err != nil {
			return err
		}
		fmt.Printf("Total episodes:\t%d\n", n)
		as, err := query.GetFinishedAnime(db)
		if err != nil {
			return err
		}
		fmt.Printf("Finished anime:\t%d\n", len(as))

		n, err = query.GetWatchedEpisodeCount(db)
		if err != nil {
			return err
		}
		fmt.Printf("Watched episodes:\t%d\n", n)
		m, err := query.GetWatchedMinutes(db)
		if err != nil {
			return err
		}
		fmt.Printf("Watched minutes:\t%d\n", m)
		fmt.Printf("Watched hours:\t%.3f\n", float64(m)/60)
		fmt.Printf("Watched days:\t%.3f\n", float64(m)/60/24)

		n, err = query.GetWatchingCount(db)
		if err != nil {
			return err
		}
		fmt.Printf("Watching anime:\t%d\n", n)
		return nil
	},
}
