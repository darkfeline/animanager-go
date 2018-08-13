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

package query

import (
	"database/sql"
)

// DeleteWatching deletes the watching entry for an anime from the
// database.  This function returns ErrMissing if no such entry
// exists.
func DeleteWatching(db *sql.DB, aid int) error {
	t, err := db.Begin()
	defer t.Rollback()
	r, err := t.Exec(`DELETE FROM watching WHERE aid=?`, aid)
	if err != nil {
		return err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrMissing
	}
	return t.Commit()
}
