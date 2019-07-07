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

	"github.com/google/subcommands"

	"go.felesatra.moe/animanager/internal/anidb/titles"
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

func (c *UpdateTitles) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if c.file == "" {
		if err := titles.UpdateCacheFromAPI(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return subcommands.ExitFailure
		}
	} else {
		f, err := os.Open(c.file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return subcommands.ExitFailure
		}
		defer f.Close()
		r, err := gzip.NewReader(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return subcommands.ExitFailure
		}
		defer r.Close()
		d, err := ioutil.ReadAll(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return subcommands.ExitFailure
		}
		if err := titles.UpdateCache(d); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return subcommands.ExitFailure
		}
	}
	return subcommands.ExitSuccess
}
