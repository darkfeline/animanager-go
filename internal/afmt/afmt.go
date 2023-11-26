// Copyright (C) 2019  Allen Li
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package afmt provides Animanager object formatting functions.
package afmt

import (
	"bufio"
	"fmt"
	"io"

	"go.felesatra.moe/anidb"
	"go.felesatra.moe/animanager/internal/query"
)

func PrintAnime(w io.Writer, a *query.Anime) error {
	bw := bufio.NewWriter(w)
	fmt.Fprintf(bw, "AID: %d\n", a.AID)
	fmt.Fprintf(bw, "Title: %s\n", a.Title)
	fmt.Fprintf(bw, "Type: %s\n", a.Type)
	fmt.Fprintf(bw, "Episodes: %d\n", a.EpisodeCount)
	fmt.Fprintf(bw, "Start date: %s\n", a.StartDate())
	fmt.Fprintf(bw, "End date: %s\n", a.EndDate())
	return bw.Flush()
}

func PrintAnimeShort(w io.Writer, a *query.Anime) error {
	_, err := fmt.Fprintf(w, "%d\t%s\t%d eps\n", a.AID, a.Title, a.EpisodeCount)
	return err
}

func PrintAnimeT(w io.Writer, ts []anidb.AnimeT) error {
	bw := bufio.NewWriter(w)
	for _, at := range ts {
		fmt.Fprintf(bw, "%d\t", at.AID)
		first := true
		for _, t := range at.Titles {
			if t.Lang != "x-jat" && t.Lang != "en" {
				continue
			}
			if !first {
				fmt.Fprint(bw, " | ")
			}
			fmt.Fprint(bw, t.Name)
			first = false
		}
		fmt.Fprint(bw, "\n")
	}
	return bw.Flush()
}

func PrintEpisode(w io.Writer, e query.Episode) error {
	bw := bufio.NewWriter(w)
	fmt.Fprintf(bw, "%d\t", e.EID)
	fmt.Fprintf(bw, "%s%d\t", e.Type.Prefix(), e.Number)
	if e.UserWatched {
		fmt.Fprintf(bw, "W ")
	} else {
		fmt.Fprintf(bw, ". ")
	}
	fmt.Fprintf(bw, "%s\t", e.Title)
	fmt.Fprintf(bw, "(%d min)", e.Length)
	fmt.Fprintf(bw, "\n")
	return bw.Flush()
}

// PrintWatchableOption provides options for PrintWatchable.
type PrintWatchableOption struct {
	IncludeWatched      bool
	IncludeMissingFiles bool
	// NumWatchable is the number of watchable episodes to print.
	// If zero, use the default value.  If negative, print all
	// watchable episodes.
	NumWatchable int
}

const defaultNumWatchable = 1

// PrintWatchable prints the watchable episodes of an anime.
func PrintWatchable(w io.Writer, a *query.Anime, efs []query.EpisodeFiles, o PrintWatchableOption) error {
	if o.NumWatchable == 0 {
		o.NumWatchable = defaultNumWatchable
	}
	bw := bufio.NewWriter(w)
	var printed int
	for i, ef := range efs {
		e := ef.Episode
		// Skip uninteresting episode types.
		if e.Type == query.EpCredit || e.Type == query.EpTrailer {
			continue
		}
		// Skip if done.
		if e.UserWatched && !o.IncludeWatched {
			continue
		}
		// Skip if no files.
		if len(ef.Files) == 0 && !o.IncludeMissingFiles {
			continue
		}
		// If we have already printed enough episodes,
		// stop looping and just print that there are
		// more.
		if o.NumWatchable > 0 && printed >= o.NumWatchable {
			fmt.Fprint(bw, "MORE\t...\n")
			break
		}
		// Print anime and previous episode if we are
		// printing the first episode for an anime.
		if printed == 0 {
			PrintAnimeShort(bw, a)
			if i > 0 {
				PrintEpisode(bw, efs[i-1].Episode)
			}
		}
		PrintEpisode(bw, e)
		printed++
		for _, f := range ef.Files {
			fmt.Fprintf(bw, "\t\t  %s\n", f.Path)
		}
		if len(ef.Files) == 0 {
			fmt.Fprintf(bw, "\t\t  <NO FILES>\n")
		}
	}
	if printed > 0 {
		fmt.Fprintln(bw)
	}
	return bw.Flush()
}
