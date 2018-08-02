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

// Package date implements a date type.
package date

import "time"

type Date int64

func NewString(s string) (Date, error) {
	t, err := time.Parse("2006-01-02", s)
	return Date(t.Unix()), err
}

func (d Date) Time() time.Time {
	return time.Unix(int64(d), 0).UTC()
}

func (d Date) String() string {
	t := d.Time()
	return t.Format("2006-01-02")
}
