// Copyright (C) 2023  Allen Li
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

package fileid

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"go.felesatra.moe/anidb/udpapi"
	"go.felesatra.moe/animanager/internal/query"
	"go.felesatra.moe/animanager/internal/udp"
	"go.felesatra.moe/hash/ed2k"
)

var fmask udpapi.FileFmask
var amask udpapi.FileAmask

func init() {
	fmask.Set("aid", "eid")
}

// MatchEpisode adds the given file as an episode file.
// Episode matching is done via AniDB UDP API.
func MatchEpisode(ctx context.Context, db *sql.DB, c *udp.Client, file string) error {
	m, err := matchFileToEpisodes(ctx, c, file)
	if err != nil {
		return fmt.Errorf("match episode: %s", err)
	}
	t, err := db.Begin()
	if err != nil {
		return fmt.Errorf("match episode: %w", err)
	}
	defer t.Rollback()
	efs := []query.EpisodeFile{{EID: m.eid, Path: file}}
	if err := query.InsertEpisodeFiles(db, efs); err != nil {
		return fmt.Errorf("match episode: %w", err)
	}
	return nil
}

// matchFileToEpisodes finds episode matches for the given file.
// Episode matching is done via AniDB UDP API.
func matchFileToEpisodes(ctx context.Context, c *udp.Client, file string) (epMatch, error) {
	f, err := os.Open(file)
	if err != nil {
		return epMatch{}, fmt.Errorf("match file to episode: %s", err)
	}
	fi, err := f.Stat()
	if err != nil {
		return epMatch{}, fmt.Errorf("match file to episode: %s", err)
	}
	h := ed2k.New()
	sum := h.Sum(nil)
	row, err := c.FileByHash(ctx, fi.Size(), fmt.Sprintf("%x", sum), fmask, amask)
	if err != nil {
		return epMatch{}, fmt.Errorf("match file to episode: %s", err)
	}
	if n := len(row); n != 3 {
		return epMatch{}, fmt.Errorf("match file to episode: unexpected number of values in response: %d", n)
	}
	aid, err := query.ParseID[query.AID](row[1])
	if err != nil {
		return epMatch{}, fmt.Errorf("match file to episode: parse aid: %s", err)
	}
	eid, err := query.ParseID[query.EID](row[2])
	if err != nil {
		return epMatch{}, fmt.Errorf("match file to episode: parse eid: %s", err)
	}
	return epMatch{
		aid: aid,
		eid: eid,
	}, nil
}

type epMatch struct {
	aid query.AID
	eid query.EID
}
