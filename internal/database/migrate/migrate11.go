// Copyright (C) 2021  Allen Li
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

package migrate

import (
	"context"
	"database/sql"
	"fmt"
)

func migrate11(ctx context.Context, d *sql.DB) error {
	c, err := d.Conn(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	t, err := c.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer t.Rollback()

	_, err = t.Exec(`
CREATE TABLE filehash (
    size INTEGER NOT NULL,
    hash TEXT NOT NULL,
    eid INTEGER,
    aid INTEGER,
    UNIQUE(size, hash)
)`)
	if err != nil {
		return fmt.Errorf("CREATE TABLE filehash: %s", err)
	}

	return t.Commit()
}
