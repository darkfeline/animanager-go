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
	"database/sql"
	"fmt"
)

func getUserVersion(d *sql.DB) (int, error) {
	r, err := d.Query("PRAGMA user_version")
	if err != nil {
		return 0, fmt.Errorf("get user version: %s", err)
	}
	defer r.Close()
	ok := r.Next()
	if !ok {
		return 0, fmt.Errorf("get user version: %s", r.Err())
	}
	var v int
	if err := r.Scan(&v); err != nil {
		return 0, fmt.Errorf("get user version: %s", err)
	}
	r.Close()
	if err := r.Err(); err != nil {
		return 0, fmt.Errorf("get user version: %s", err)
	}
	return v, nil
}

func setUserVersion(d *sql.DB, v int) error {
	_, err := d.Exec(fmt.Sprintf("PRAGMA user_version=%d", v))
	if err != nil {
		return fmt.Errorf("set user version %d: %s", v, err)
	}
	return nil
}
