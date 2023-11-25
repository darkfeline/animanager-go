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

func migrate7(ctx context.Context, d *sql.DB) error {
	c, err := d.Conn(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	_, err = c.ExecContext(ctx, "PRAGMA foreign_keys = 0")
	if err != nil {
		return err
	}
	defer c.ExecContext(ctx, "PRAGMA foreign_keys = 1")

	t, err := c.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer t.Rollback()
	_, err = t.Exec(`
CREATE TABLE episode_new (
    id INTEGER,
    eid INTEGER NOT NULL,
    aid INTEGER NOT NULL,
    type INTEGER NOT NULL,
    number INTEGER NOT NULL,
    title TEXT NOT NULL,
    length INTEGER NOT NULL,
    user_watched INTEGER NOT NULL CHECK (user_watched IN (0, 1))
        DEFAULT 0,
    PRIMARY KEY (id),
    UNIQUE (eid),
    UNIQUE (aid, type, number),
    FOREIGN KEY (aid) REFERENCES anime (aid)
        ON DELETE CASCADE ON UPDATE CASCADE
)`)
	if err != nil {
		return fmt.Errorf("CREATE TABLE episode_new: %s", err)
	}
	_, err = t.Exec(`
INSERT INTO episode_new
(id, eid, aid, type, number, title, length, user_watched)
SELECT id, eid, aid, type, number, title, length, user_watched
FROM episode`)
	if err != nil {
		return fmt.Errorf("INSERT INTO episode_new: %s", err)
	}
	_, err = t.Exec(`DROP TABLE episode`)
	if err != nil {
		return fmt.Errorf("DROP TABLE episode: %s", err)
	}
	_, err = t.Exec(`ALTER TABLE episode_new RENAME TO episode`)
	if err != nil {
		return fmt.Errorf("ALTER TABLE episode_new: %s", err)
	}
	return t.Commit()
}
