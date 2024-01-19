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
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"go.felesatra.moe/anidb/udpapi"
	"go.felesatra.moe/anidb/udpapi/codes"
	"go.felesatra.moe/animanager/internal/query"
	"go.felesatra.moe/animanager/internal/sqlc"
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
	fh, err := m.findFileHash(ctx, file)
	if err != nil {
		return fmt.Errorf("match episode: %s", err)
	}
	m.l = m.l.With("FileHash", fh)
	if fh.EID == 0 {
		slog.Debug("file hash missing EID")
		return nil
	}
	efs := []query.EpisodeFile{{Eid: fh.EID, Path: file}}
	slog.Debug("matched fie to EID", "eid", fh.EID)
	if err := query.InsertEpisodeFiles(ctx, sqlc.New(m.db), m.l, efs); err != nil {
		return fmt.Errorf("match episode: %w", err)
	}
	return nil
}

// findFileHash finds episode matches for the given file.
// Episode matching is done via AniDB UDP API.
func (m Matcher) findFileHash(ctx context.Context, file string) (*query.FileHash, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("match file to episode: %s", err)
	}
	fk := fileKey{
		Path: file,
	}
	if err := fk.populateSize(); err != nil {
		return nil, fmt.Errorf("match file to episode: %s", err)
	}

	// Try getting from cache by size first, since hashing is slow.
	fh, err := m.findFileHashBySize(ctx, fk)
	if err == nil {
		m.l.Debug("got file hash from cache by size", "FileHash", fh)
		return fh, nil
	}
	m.l.Debug("error getting file hash from cache by size", "error", err)

	if err := fk.populateHash(); err != nil {
		return nil, fmt.Errorf("match file to episode: %s", err)
	}

	// Safe because method has value receiver
	m.l = m.l.With("fileKey", fk)
	m.l.Debug("got file key")

	// Try getting from cache
	fh, err = query.GetFileHash(m.db, fk.Size, fk.Hash)
	if err == nil {
		m.l.Debug("got file hash from cache", "FileHash", fh)
		return fh, nil
	}
	m.l.Debug("error getting file hash from cache", "error", err)

	// Get from AniDB
	fh, err = m.requestFileHash(ctx, fk)
	if err != nil {
		return nil, fmt.Errorf("match file to episode: %s", err)
	}
	fh.Filename = filepath.Base(file)
	m.l = m.l.With("FileHash", fh)
	m.l.Debug("lookup file hash completed")

	// Add to cache
	if err := query.InsertFileHash(m.db, fh); err != nil {
		return nil, fmt.Errorf("match file to episode: %s", err)
	}
	m.l.Debug("added file hash to cache")
	return fh, nil
}

func (m Matcher) findFileHashBySize(ctx context.Context, fk fileKey) (*query.FileHash, error) {
	fhs, err := query.GetFileHashBySize(m.db, fk.Size)
	if err != nil {
		return nil, fmt.Errorf("findFileHashBySize: %s", err)
	}
	if len(fhs) == 0 {
		return nil, errors.New("findFileHashBySize: no matches by size")
	}
	if len(fhs) == 1 {
		fh := fhs[0]
		m.l.Debug("got single file hash from cache by size", "FileHash", fh)
		return &fh, nil
	}
	// Try to match by filename.
	name := filepath.Base(fk.Path)
	for _, fh := range fhs {
		if fh.Filename == name {
			return &fh, nil
		}
	}
	return nil, fmt.Errorf("findFileHashBySize: no matches by name for %d resultss", len(fhs))
}

// requestFileHash requests a [FileHash] with the [UDPClient].
// If the request does not find a file, then the returned FileHash will
// have zero EID and AID values.
func (m Matcher) requestFileHash(ctx context.Context, fk fileKey) (*query.FileHash, error) {
	fh := query.FileHash{
		Size: fk.Size,
		Hash: fk.Hash,
	}
	row, err := m.c.FileByHash(ctx, fk.Size, string(fk.Hash), fmask, amask)
	if err != nil {
		if !errors.Is(err, codes.NO_SUCH_FILE) {
			return nil, fmt.Errorf("lookup file hash: %s", err)
		}
		m.l.Debug("got no such file response")
	} else {
		m.l.Debug("got file hash response", "row", row)
		if err := parseFileHashRow(&fh, row); err != nil {
			return nil, fmt.Errorf("lookup file hash: %s", err)
		}
	}
	return &fh, nil
}

func parseFileHashRow(fh *query.FileHash, row []string) error {
	if n := len(row); n != 3 {
		return fmt.Errorf("parse file has row: unexpected number of values in response: %d", n)
	}
	aid, err := sqlc.ParseID[sqlc.AID](row[1])
	if err != nil {
		return fmt.Errorf("parse file has row: parse aid: %s", err)
	}
	eid, err := sqlc.ParseID[sqlc.EID](row[2])
	if err != nil {
		return fmt.Errorf("parse file has row: parse eid: %s", err)
	}
	fh.EID = eid
	fh.AID = aid
	return nil
}

type fileKey struct {
	Path string
	Size int64
	Hash sqlc.Hash
}

func (k *fileKey) populateSize() error {
	f, err := os.Open(k.Path)
	if err != nil {
		return fmt.Errorf("populate size: %s", err)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("populate size: %s", err)
	}
	k.Size = fi.Size()
	return nil

}

func (k *fileKey) populateHash() error {
	f, err := os.Open(k.Path)
	if err != nil {
		return fmt.Errorf("populate hash: %s", err)
	}
	defer f.Close()
	h := ed2k.New()
	if _, err := io.CopyBuffer(h, f, copyBuf()); err != nil {
		return fmt.Errorf("populate hash: %s", err)
	}
	sum := h.Sum(nil)
	k.Hash = formatHash(sum)
	return nil
}

// getCopyBuf returns a large copy buffer for [getFileKey].
// Since video files are kinda big, a small buffer makes hashing slow.
// Allocate lazily so most commands don't do this.
// Use 3D2k chunk size for this as that's the hash that is used.
var copyBuf = sync.OnceValue(func() []byte { return make([]byte, 9728000) })

func formatHash(sum []byte) sqlc.Hash {
	return sqlc.Hash(fmt.Sprintf("%x", sum))
}
