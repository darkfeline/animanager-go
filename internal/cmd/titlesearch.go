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
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/google/subcommands"
	"go.felesatra.moe/anidb"

	"go.felesatra.moe/animanager/internal/anidb/titles"
)

type TitleSearch struct {
	skipCache bool
}

func (*TitleSearch) Name() string     { return "titlesearch" }
func (*TitleSearch) Synopsis() string { return "Search for an anime title." }
func (*TitleSearch) Usage() string {
	return `Usage: titlesearch terms...
Search for an anime title.
`
}

func (t *TitleSearch) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&t.skipCache, "skipcache", false, "Ignore local titles cache.")
}

func (t *TitleSearch) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() == 0 {
		fmt.Fprint(os.Stderr, t.Usage())
		return subcommands.ExitUsageError
	}
	terms := f.Args()
	var ts []anidb.AnimeT
	var err error
	if t.skipCache {
		ts, err = titles.GetSkipCache()
	} else {
		ts, err = titles.Get()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return subcommands.ExitFailure
	}
	ts = titles.Search(ts, terms)
	printAnimeT(os.Stdout, ts)
	return subcommands.ExitSuccess
}

func printAnimeT(w io.Writer, ts []anidb.AnimeT) error {
	bw := bufio.NewWriter(w)
	for _, at := range ts {
		fmt.Fprintf(bw, "%d\t", at.AID)
		first := true
		for _, t := range at.Titles {
			if t.Lang != "x-jat" && t.Lang != "en" {
				continue
			}
			if !first {
				fmt.Fprint(bw, ", ")
			}
			fmt.Fprint(bw, t.Name)
			first = false
		}
		fmt.Fprint(bw, "\n")
	}
	return bw.Flush()
}
