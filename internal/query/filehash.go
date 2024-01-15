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

package query

import (
	"context"
	"fmt"

	"go.felesatra.moe/animanager/internal/sqlc"
)

// A Hash is an eD2k formatted as a hex string.
type Hash string

func (h *Hash) Scan(src any) error {
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("wrong type %T for %T", src, h)
	}
	*h = Hash(s)
	return nil
}

type FileHash struct {
	_table   struct{} `sql:"filehash"`
	Size     int64    `sql:"size"`
	Hash     Hash     `sql:"hash"`
	EID      sqlc.EID `sql:"eid"`
	AID      sqlc.AID `sql:"aid"`
	Filename string   `sql:"filename"`
}

func InsertFileHash(db sqlc.DBTX, fh *FileHash) error {
	ctx := context.Background()
	p := sqlc.InsertFileHashParams{
		Size: fh.Size,
		Hash: string(fh.Hash),
	}
	if fh.EID != 0 {
		p.Eid.Int64 = int64(fh.EID)
		p.Eid.Valid = true
	}
	if fh.AID != 0 {
		p.Aid.Int64 = int64(fh.AID)
		p.Aid.Valid = true
	}
	if fh.Filename != "" {
		p.Filename.String = fh.Filename
		p.Filename.Valid = true
	}
	return sqlc.New(db).InsertFileHash(ctx, p)
}

func GetFileHash(db sqlc.DBTX, size int64, hash Hash) (*FileHash, error) {
	ctx := context.Background()
	p := sqlc.GetFileHashParams{
		Size: size,
		Hash: string(hash),
	}
	fh, err := sqlc.New(db).GetFileHash(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("GetFileHash %d %q: %s", size, hash, err)
	}
	fh2 := convertFileHash(fh)
	return &fh2, nil
}

func convertFileHash(v sqlc.Filehash) FileHash {
	return FileHash{
		Size:     v.Size,
		Hash:     Hash(v.Hash),
		EID:      sqlc.EID(v.Eid.Int64),
		AID:      sqlc.AID(v.Aid.Int64),
		Filename: string(v.Filename.String),
	}
}
