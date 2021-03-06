// Copyright (C) 2018  Allen Li
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
)

func migrate3(ctx context.Context, d *sql.DB) error {
	t, err := d.Begin()
	if err != nil {
		return err
	}
	defer t.Rollback()
	_, err = t.Exec(`
CREATE TABLE anime (
    aid INTEGER,
    title TEXT NOT NULL UNIQUE,
    type TEXT NOT NULL,
    episodecount INTEGER NOT NULL,
    startdate INTEGER,
    enddate INTEGER,
    PRIMARY KEY (aid)
)`)
	if err != nil {
		return err
	}
	_, err = t.Exec("CREATE INDEX anime_titles ON anime (title)")
	if err != nil {
		return err
	}
	_, err = t.Exec(`
CREATE TABLE episode_type (
    id INTEGER,
    name TEXT NOT NULL UNIQUE,
    prefix TEXT NOT NULL UNIQUE,
    PRIMARY KEY(id)
)`)
	if err != nil {
		return err
	}
	_, err = t.Exec(`
INSERT INTO episode_type (id, name, prefix) VALUES
(1, 'regular', ''),
(2, 'special', 'S'),
(3, 'credit', 'C'),
(4, 'trailer', 'T'),
(5, 'parody', 'P'),
(6, 'other', 'O')`)
	if err != nil {
		return err
	}
	_, err = t.Exec(`
CREATE TABLE episode (
    id INTEGER,
    aid INTEGER NOT NULL,
    type INTEGER NOT NULL,
    number INTEGER NOT NULL,
    title TEXT NOT NULL,
    length INTEGER NOT NULL,
    user_watched INTEGER NOT NULL CHECK (user_watched IN (0, 1))
	DEFAULT 0,
    PRIMARY KEY (id),
    UNIQUE (aid, type, number),
    FOREIGN KEY (aid) REFERENCES anime (aid)
	ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (type) REFERENCES episode_type (id)
	ON DELETE RESTRICT ON UPDATE CASCADE
)`)
	if err != nil {
		return err
	}
	_, err = t.Exec(`
CREATE TABLE watching (
    aid INTEGER,
    regexp TEXT NOT NULL,
    PRIMARY KEY (aid),
    FOREIGN KEY (aid) REFERENCES anime (aid)
	ON DELETE CASCADE ON UPDATE CASCADE
)`)
	if err != nil {
		return err
	}
	_, err = t.Exec(`
CREATE TABLE file_priority (
     id INTEGER PRIMARY KEY,
     regexp TEXT NOT NULL,
     priority INTEGER NOT NULL
)`)
	if err != nil {
		return err
	}
	return t.Commit()
}
