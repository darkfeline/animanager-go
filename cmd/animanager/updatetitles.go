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
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"

	"go.felesatra.moe/anidb"
)

var updateTitlesCmd = command{
	usageLine: "update-titles [-file name]",
	shortDesc: "update AniDB titles database",
	longDesc: `Update AniDB titles database.
`,
	run: func(cmd *command, args []string) error {
		f := cmd.flagSet()
		file := f.String("file", "", "Titles file to use.")
		if err := f.Parse(args); err != nil {
			return err
		}

		if *file == "" {
			return updateCacheFromAPI()
		}
		f2, err := os.Open(*file)
		if err != nil {
			return err
		}
		defer f2.Close()
		r, err := gzip.NewReader(f2)
		if err != nil {
			return err
		}
		defer r.Close()
		d, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}
		if err := updateCache(d); err != nil {
			return err
		}
		return nil
	},
}

func updateCacheFromAPI() error {
	c, err := anidb.DefaultTitlesCache()
	if err != nil {
		return fmt.Errorf("update cache from api: %w", err)
	}
	if _, err := c.GetFreshTitles(); err != nil {
		return fmt.Errorf("update cache from api: %w", err)
	}
	if err := c.Save(); err != nil {
		return fmt.Errorf("update cache from api: %w", err)
	}
	return nil
}

func updateCache(d []byte) error {
	ts, err := anidb.DecodeTitles(d)
	if err != nil {
		return fmt.Errorf("update cache: %w", err)
	}
	c, err := anidb.DefaultTitlesCache()
	if err != nil {
		return fmt.Errorf("update cache from api: %w", err)
	}
	c.Titles = ts
	if err := c.Save(); err != nil {
		return fmt.Errorf("update cache from api: %w", err)
	}
	return nil
}
