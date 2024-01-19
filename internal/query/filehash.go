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

type FileHash struct {
	_table   struct{}  `sql:"filehash"`
	Size     int64     `sql:"size"`
	Hash     sqlc.Hash `sql:"hash"`
	EID      sqlc.EID  `sql:"eid"`
	AID      sqlc.AID  `sql:"aid"`
	Filename string    `sql:"filename"`
}

func InsertFileHash(db sqlc.DBTX, fh *FileHash) error {
	ctx := context.Background()
	p := sqlc.InsertFileHashParams{
		Size:     fh.Size,
		Hash:     fh.Hash,
		Eid:      fh.EID,
		Aid:      fh.AID,
		Filename: fh.Filename,
	}
	return sqlc.New(db).InsertFileHash(ctx, p)
}

func GetFileHash(db sqlc.DBTX, size int64, hash sqlc.Hash) (*FileHash, error) {
	ctx := context.Background()
	p := sqlc.GetFileHashParams{
		Size: size,
		Hash: hash,
	}
	fh, err := sqlc.New(db).GetFileHash(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("GetFileHash %d %q: %s", size, hash, err)
	}
	fh2 := convertFileHash(fh)
	return &fh2, nil
}

func GetFileHashBySize(db sqlc.DBTX, size int64) ([]FileHash, error) {
	ctx := context.Background()
	fh, err := sqlc.New(db).GetFileHashBySize(ctx, size)
	if err != nil {
		return nil, fmt.Errorf("GetFileHashBySize %d: %s", size, err)
	}
	return smap(fh, convertFileHash), nil
}

func convertFileHash(v sqlc.Filehash) FileHash {
	return FileHash{
		Size:     v.Size,
		Hash:     v.Hash,
		EID:      sqlc.EID(v.Eid),
		AID:      sqlc.AID(v.Aid),
		Filename: v.Filename,
	}
}
