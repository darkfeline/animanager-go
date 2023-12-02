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

	"go.felesatra.moe/animanager/cmd/animanager/vars"
	"go.felesatra.moe/animanager/internal/query"
)

var clearFilesCmd = command{
	usageLine: "clearfiles",
	shortDesc: "clears episode files",
	longDesc: `Clears all episode files.
`,
	run: func(h *handle, args []string) error {
		f := h.flagSet()
		cfgv := vars.Config(f)
		if err := f.Parse(args); err != nil {
			return err
		}
		if f.NArg() != 0 {
			return errors.New("no arguments allowed")
		}

		db, err := cfgv.OpenDB()
		if err != nil {
			return err
		}
		defer db.Close()
		if err := query.DeleteAllEpisodeFiles(db); err != nil {
			return err
		}
		return nil
	},
}
