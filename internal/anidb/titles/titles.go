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

// Package titles implements convenient AniDB titles getting and searching.
package titles

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"go.felesatra.moe/anidb"
	"go.felesatra.moe/anidb/cache/titles"
)

// ErrLogger is used by this package to log errors.
var ErrLogger = log.New(os.Stderr, "titles: ", log.LstdFlags)

// Logger is used by this package for informational logging.
var Logger = log.New(ioutil.Discard, "titles: ", log.LstdFlags)

// Get returns a slice of anime titles.  This function uses a cache if
// it exists, getting the data from AniDB and caching it otherwise.
func Get() ([]anidb.AnimeT, error) {
	ts, err := titles.LoadDefault()
	if err == nil {
		return ts, nil
	}
	Logger.Printf("Failed to load cached titles: %s", err)
	return GetSkipCache()
}

// GetSkipCache returns a slice of anime titles.  This function ignores
// the cache, getting the data from AniDB and caching it.
func GetSkipCache() ([]anidb.AnimeT, error) {
	ts, err := anidb.RequestTitles()
	if err != nil {
		return nil, err
	}
	err = titles.SaveDefault(ts)
	if err != nil {
		ErrLogger.Printf("Error caching titles: %s", err)
	}
	return ts, nil
}

// Search returns a slice of anime whose title matches the given
// terms.  A title is matched if it contains all terms in order,
// ignoring case and intervening characters.
func Search(at []anidb.AnimeT, terms []string) []anidb.AnimeT {
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
