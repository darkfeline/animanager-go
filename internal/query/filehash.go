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
	"database/sql"
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
	EID      EID      `sql:"eid"`
	AID      AID      `sql:"aid"`
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

func GetFileHash(db *sql.DB, size int64, hash Hash) (*FileHash, error) {
	r := db.QueryRow(`
SELECT size, hash, eid, aid, filename FROM filehash WHERE size=? AND hash=?`,
		size, hash)
	var fh FileHash
	if err := r.Scan(&fh.Size, &fh.Hash, &fh.EID, &fh.AID, &fh.Filename); err != nil {
		return nil, err
	}
	return &fh, nil
}
