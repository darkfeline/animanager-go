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

// Package date implements a date type stored as Unix timestamps.
package date

import (
	"database/sql"
	"time"
)

// Date represents a date as a Unix timestamp at 00:00 UTC of the
// date.
type Date int64

// Zero is 0000-01-01.
const Zero Date = -62167219200

// Parse parses and returns a Date from a string in YYYY-MM-DD
// format.
func Parse(s string) (Date, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return 0, err
	}
	return Date(t.Unix()), nil
}

// Time returns the Time representation of the date.
func (d Date) Time() time.Time {
	return time.Unix(int64(d), 0).UTC()
}

// NullInt64 returns the SQL representation of the date.
func (d Date) NullInt64() sql.NullInt64 {
	return sql.NullInt64{
		Int64: int64(d),
		Valid: true,
	}
}

// String returns the date formatted as YYYY-MM-DD.
func (d Date) String() string {
	t := d.Time()
	return t.Format("2006-01-02")
}

// FromTime returns the date of the given time.
func FromTime(t time.Time) Date {
	return Date(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).Unix())
}

// New returns a new date.
func New(year int, month time.Month, day int) Date {
	return FromTime(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
}

// Today returns the date for today.
func Today() Date {
	return FromTime(time.Now())
}
