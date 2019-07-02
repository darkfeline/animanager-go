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
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
	"go.felesatra.moe/anidb"

	"go.felesatra.moe/animanager/internal/afmt"
	"go.felesatra.moe/animanager/internal/anidb/titles"
)

type Search struct {
	skipCache bool
}

func (*Search) Name() string     { return "search" }
func (*Search) Synopsis() string { return "Search for an anime title." }
func (*Search) Usage() string {
	return `Usage: search terms...
Search for an anime title.
`
}

func (s *Search) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&s.skipCache, "skipcache", false, "Ignore local titles cache.")
}

func (s *Search) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() == 0 {
		fmt.Fprint(os.Stderr, s.Usage())
		return subcommands.ExitUsageError
	}
	terms := f.Args()
	var ts []anidb.AnimeT
	var err error
	if s.skipCache {
		ts, err = titles.GetSkipCache()
	} else {
		ts, err = titles.Get()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return subcommands.ExitFailure
	}
	ts = titles.Search(ts, terms)
	afmt.PrintAnimeT(os.Stdout, ts)
	return subcommands.ExitSuccess
}
