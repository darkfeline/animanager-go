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
	"log/slog"
	"os"
	"path/filepath"

	"go.felesatra.moe/anidb/udpapi"
	"go.felesatra.moe/animanager/internal/query"
	"go.felesatra.moe/hash/ed2k"
)

var fmask udpapi.FileFmask
var amask udpapi.FileAmask

func init() {
	fmask.Set("aid", "eid")
}

type UDPClient interface {
	FileByHash(context.Context, int64, string, udpapi.FileFmask, udpapi.FileAmask) ([]string, error)
}

type Matcher struct {
	l  *slog.Logger
	db *sql.DB
	c  UDPClient
}

func NewMatcher(l *slog.Logger, db *sql.DB, c UDPClient) Matcher {
	return Matcher{
		l:  l,
		db: db,
		c:  c,
	}
}

// MatchEpisode adds the given file as an episode file.
// Episode matching is done via AniDB UDP API.
func (m Matcher) MatchEpisode(ctx context.Context, file string) error {
	// Safe because method has value receiver
	m.l = m.l.With("file", file)
	fh, err := m.matchFileToEpisode(ctx, file)
	if err != nil {
		return fmt.Errorf("match episode: %s", err)
	}
	if fh.EID == 0 {
		slog.Debug("file hash missing EID", "FileHash", fh)
		return nil
	}
	efs := []query.EpisodeFile{{EID: fh.EID, Path: file}}
	if err := query.InsertEpisodeFiles(m.db, m.l, efs); err != nil {
		return fmt.Errorf("match episode: %w", err)
	}
	return nil
}

// matchFileToEpisodes finds episode matches for the given file.
// Episode matching is done via AniDB UDP API.
func (m Matcher) matchFileToEpisode(ctx context.Context, file string) (*query.FileHash, error) {
	fk, err := getFileKey(file)
	if err != nil {
		return nil, fmt.Errorf("match file to episode: %s", err)
	}
	// Safe because method has value receiver
	m.l = m.l.With("size", fk.size, "hash", fk.hash)

	// Try getting from cache
	fh, err := query.GetFileHash(m.db, fk.size, fk.hash)
	if err == nil {
		m.l.Debug("got file hash from db", "FileHash", fh)
		return fh, nil
	}
	m.l.Debug("error getting file hash from db", "error", err)

	// Get from AniDB
	row, err := m.c.FileByHash(ctx, fk.size, string(fk.hash), fmask, amask)
	if err != nil {
		return nil, fmt.Errorf("match file to episode: %s", err)
	}
	m.l.Debug("got file hash response", "row", row)
	fh, err = parseFileHashRow(row)
	if err != nil {
		return nil, fmt.Errorf("match file to episode: %s", err)
	}
	fh.Size = fk.size
	fh.Hash = fk.hash
	fh.Filename = filepath.Base(file)

	// Add to cache
	if err := query.InsertFileHash(m.db, fh); err != nil {
		return nil, fmt.Errorf("match file to episode: %s", err)
	}
	return fh, nil
}

func parseFileHashRow(row []string) (*query.FileHash, error) {
	if n := len(row); n != 3 {
		return nil, fmt.Errorf("parse file has row: unexpected number of values in response: %d", n)
	}
	aid, err := query.ParseID[query.AID](row[1])
	if err != nil {
		return nil, fmt.Errorf("parse file has row: parse aid: %s", err)
	}
	eid, err := query.ParseID[query.EID](row[2])
	if err != nil {
		return nil, fmt.Errorf("parse file has row: parse eid: %s", err)
	}
	return &query.FileHash{
		EID: eid,
		AID: aid,
	}, nil
}

type fileKey struct {
	size int64
	hash query.Hash
}

func getFileKey(file string) (fileKey, error) {
	f, err := os.Open(file)
	if err != nil {
		return fileKey{}, fmt.Errorf("get file key: %s", err)
	}
	fi, err := f.Stat()
	if err != nil {
		return fileKey{}, fmt.Errorf("get file key: %s", err)
	}
	h := ed2k.New()
	sum := h.Sum(nil)
	return fileKey{
		size: fi.Size(),
		hash: formatHash(sum),
	}, nil
}

func formatHash(sum []byte) query.Hash {
	return query.Hash(fmt.Sprintf("%x", sum))
}
