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

func migrate8(ctx context.Context, d *sql.DB) error {
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
CREATE TABLE episode_file_new (
    id INTEGER,
    episode_id INTEGER NOT NULL,
    eid INTEGER NOT NULL,
    path TEXT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (eid) REFERENCES episode (eid)
        ON DELETE CASCADE ON UPDATE CASCADE
)`)
	if err != nil {
		return fmt.Errorf("CREATE TABLE episode_file_new: %s", err)
	}
	_, err = t.Exec(`
INSERT INTO episode_file_new
(id, episode_id, eid, path)
SELECT episode_file.id, episode_file.episode_id, episode.eid, episode_file.path
FROM episode_file
JOIN episode ON (episode_file.episode_id = episode.id)`)
	if err != nil {
		return fmt.Errorf("INSERT INTO episode_file_new: %s", err)
	}
	_, err = t.Exec(`DROP TABLE episode_file`)
	if err != nil {
		return fmt.Errorf("DROP TABLE episode_file: %s", err)
	}
	_, err = t.Exec(`ALTER TABLE episode_file_new RENAME TO episode_file`)
	if err != nil {
		return fmt.Errorf("ALTER TABLE episode_file_new: %s", err)
	}
	return t.Commit()
}
