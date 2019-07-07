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
	"regexp"
	"strings"

	"github.com/google/subcommands"

	"go.felesatra.moe/anidb"
	"go.felesatra.moe/anidb/cache/titles"
	"go.felesatra.moe/animanager/internal/afmt"
)

type Search struct {
}

func (*Search) Name() string     { return "search" }
func (*Search) Synopsis() string { return "Search for an anime title." }
func (*Search) Usage() string {
	return `Usage: search terms...
Search for an anime title.
`
}

func (s *Search) SetFlags(f *flag.FlagSet) {
}

func (s *Search) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() == 0 {
		fmt.Fprint(os.Stderr, s.Usage())
		return subcommands.ExitUsageError
	}
	terms := f.Args()
	ts, err := titles.LoadDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return subcommands.ExitFailure
	}
	ts = search(ts, terms)
	afmt.PrintAnimeT(os.Stdout, ts)
	return subcommands.ExitSuccess
}

// search returns a slice of anime whose title matches the given
// terms.  A title is matched if it contains all terms in order,
// ignoring case and intervening characters.
func search(at []anidb.AnimeT, terms []string) []anidb.AnimeT {
	r := globTerms(terms)
	return filterTitles(r, at)
}

// globTerms returns a regexp that matches strings containing the
// terms in order, ignoring case and intervening characters.
func globTerms(terms []string) *regexp.Regexp {
	for i, t := range terms {
		terms[i] = regexp.QuoteMeta(t)
	}
	return regexp.MustCompile("(?i)" + strings.Join(terms, ".*"))
}

// filterTitles returns a slice of anime whose title matches the regexp.
func filterTitles(r *regexp.Regexp, ts []anidb.AnimeT) []anidb.AnimeT {
	var matched []anidb.AnimeT
	for _, at := range ts {
		if titleMatches(r, at.Titles) {
			matched = append(matched, at)
		}
	}
	return matched
}

// titleMatches returns true if any of the titles matches the regexp.
func titleMatches(r *regexp.Regexp, ts []anidb.Title) bool {
	for _, t := range ts {
		if r.FindStringIndex(t.Name) != nil {
			return true
		}
	}
	return false
}
