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
	"log"

	"go.felesatra.moe/animanager/cmd/animanager/vars"
	"go.felesatra.moe/animanager/internal/fileid"
)

var findFilesCmd = command{
	usageLine: "findfiles",
	shortDesc: "find episode files",
	longDesc: `Find episode files.
`,
	run: func(cmd *command, args []string) error {
		f := cmd.flagSet()
		cfgv := vars.Config(f)
		if err := f.Parse(args); err != nil {
			return err
		}
		if f.NArg() != 0 {
			return errors.New("no arguments allowed")
		}
		cfg, err := cfgv.Load()
		if err != nil {
			return err
		}

		log.Printf("Finding video files...")
		files, err := fileid.FindVideoFiles(cfg.WatchDirs)
		if err != nil {
			return err
		}
		log.Printf("Finished finding video files")

		db, err := cfgv.OpenDB()
		if err != nil {
			return err
		}
		defer db.Close()
		if err := fileid.RefreshFiles(db, files); err != nil {
			return err
		}
		return nil
	},
}
