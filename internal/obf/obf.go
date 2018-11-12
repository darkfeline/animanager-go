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

// Package obf implements object formatting.
package obf

import (
	"fmt"
	"io"

	"go.felesatra.moe/animanager/internal/query"
)

func PrintAnime(w io.Writer, a *query.Anime) {
	fmt.Fprintf(w, "AID: %d\n", a.AID)
	fmt.Fprintf(w, "Title: %s\n", a.Title)
	fmt.Fprintf(w, "Type: %s\n", a.Type)
	fmt.Fprintf(w, "Episodes: %d\n", a.EpisodeCount)
	fmt.Fprintf(w, "Start date: %s\n", a.StartDate())
	fmt.Fprintf(w, "End date: %s\n", a.EndDate())
}

func PrintAnimeShort(w io.Writer, a *query.Anime) {
	fmt.Fprintf(w, "%d\t%s\t%d eps\n", a.AID, a.Title, a.EpisodeCount)
}

func PrintEpisode(w io.Writer, e query.Episode) {
	fmt.Fprintf(w, "%d\t", e.ID)
	fmt.Fprintf(w, "%s%d\t", e.Type.Prefix(), e.Number)
	if e.UserWatched {
		fmt.Fprintf(w, "W ")
	} else {
		fmt.Fprintf(w, ". ")
	}
	fmt.Fprintf(w, "%s\t", e.Title)
	fmt.Fprintf(w, "(%d min)", e.Length)
	fmt.Fprintf(w, "\n")
}
