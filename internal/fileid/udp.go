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
	"go.felesatra.moe/animanager/internal/udp"
	"go.felesatra.moe/hash/ed2k"
)

var fmask udpapi.FileFmask
var amask udpapi.FileAmask

// MatchEpisode adds the given file as an episode file.
// Episode matching is done via AniDB UDP API.
func MatchEpisode(ctx context.Context, db *sql.DB, c *udp.Client, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("match episode: %s", err)
	}
	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("match episode: %s", err)
	}

	h := ed2k.New()
	hash := h.Sum(nil)
	// XXXXXXXXX string hash
	rows, err := c.FileByHash(ctx, fi.Size(), string(hash), fmask, amask)
	if err != nil {
		return err
	}
	_ = rows
	panic(nil)
}
