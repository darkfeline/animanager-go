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

package migrate

import (
	"context"
	"database/sql"
	"fmt"
)

func migrate13(ctx context.Context, d *sql.DB) error {
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
CREATE TABLE filehash_new (
    size INTEGER NOT NULL,
    hash TEXT NOT NULL,
    eid INTEGER NOT NULL,
    aid INTEGER NOT NULL,
    filename TEXT NOT NULL,
    UNIQUE(size, hash)
)`)
	if err != nil {
		return fmt.Errorf("CREATE TABLE filehash_new: %s", err)
	}

	_, err = t.Exec(`
INSERT INTO filehash_new
(size, hash, eid, aid, filename)
SELECT size, hash, IFNULL(eid, 0), IFNULL(aid, 0), IFNULL(filename, '')
FROM filehash`)
	if err != nil {
		return fmt.Errorf("INSERT INTO filehash_new: %s", err)
	}
	_, err = t.Exec(`DROP TABLE filehash`)
	if err != nil {
		return fmt.Errorf("DROP TABLE filehash: %s", err)
	}
	_, err = t.Exec(`ALTER TABLE filehash_new RENAME TO filehash`)
	if err != nil {
		return fmt.Errorf("ALTER TABLE filehash_new: %s", err)
	}

	return t.Commit()
}
