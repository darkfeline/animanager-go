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
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"

	"go.felesatra.moe/animanager/cmd/animanager/vars"
	"go.felesatra.moe/animanager/internal/afmt"
	"go.felesatra.moe/animanager/internal/query"
)

var showCmd = command{
	usageLine: "show aid",
	shortDesc: "show information about a show",
	longDesc: `Show information about a show.
`,
	run: func(h *handle, args []string) error {
		f := h.flagSet()
		cfgv := vars.Config(f)
		if err := f.Parse(args); err != nil {
			return err
		}

		if f.NArg() != 1 {
			return errors.New("must pass exactly one argument")
		}
		aid, err := query.ParseID[query.AID](f.Arg(0))
		if err != nil {
			return fmt.Errorf("invalid AID %v: %v", aid, err)
		}

		db, err := cfgv.OpenDB()
		if err != nil {
			return err
		}
		defer db.Close()
		a, err := query.GetAnime(db, aid)
		if err != nil {
			return err
		}
		bw := bufio.NewWriter(os.Stdout)
		afmt.PrintAnime(bw, a)
		w, err := query.GetWatching(db, aid)
		switch {
		case err == nil:
			fmt.Fprintf(bw, "Registered: %#v (offset %d)\n", w.Regexp, w.Offset)
		case errors.Is(err, sql.ErrNoRows):
			io.WriteString(bw, "Not registered\n")
		default:
			return err
		}
		es, err := query.GetEpisodes(db, aid)
		if err != nil {
			return err
		}
		for _, e := range es {
			afmt.PrintEpisode(bw, e)
		}
		bw.Flush()
		return nil
	},
}
