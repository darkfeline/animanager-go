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
	"os"
	"regexp"
	"strings"

	"go.felesatra.moe/anidb"
	"go.felesatra.moe/animanager/internal/afmt"
)

var searchCmd = command{
	usageLine: "search [terms]",
	shortDesc: "search for an anime by title",
	longDesc: `search for an anime by title.
`,
	run: func(h *handle, args []string) error {
		f := h.flagSet()
		if err := f.Parse(args); err != nil {
			return err
		}

		if f.NArg() == 0 {
			return errors.New("no search terms")
		}
		terms := f.Args()
		tc, err := anidb.DefaultTitlesCache()
		if err != nil {
			return err
		}
		defer tc.SaveIfUpdated()
		ts, err := tc.GetTitles()
		if err != nil {
			return err
		}
		ts = search(ts, terms)
		afmt.PrintAnimeT(os.Stdout, ts)
		return nil

	},
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
