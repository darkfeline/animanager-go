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
	"compress/gzip"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"go.felesatra.moe/anidb"
	"go.felesatra.moe/animanager/internal/config"
)

type UpdateTitles struct {
	file string
}

func (*UpdateTitles) Name() string     { return "update-titles" }
func (*UpdateTitles) Synopsis() string { return "Update AniDB titles database." }
func (*UpdateTitles) Usage() string {
	return `Usage: update-titles
       update-titles -file FILE
Update AniDB titles database.
`
}

func (c *UpdateTitles) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.file, "file", "", "Titles file to use.")
}

func (c *UpdateTitles) Run(_ context.Context, f *flag.FlagSet, cfg config.Config) error {
	if c.file == "" {
		return updateCacheFromAPI()
	}
	f2, err := os.Open(c.file)
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
