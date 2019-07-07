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

// Package titles implements convenient AniDB titles getting and searching.
package titles

import (
	"go.felesatra.moe/anidb"
	"go.felesatra.moe/anidb/cache/titles"
	"golang.org/x/xerrors"
)

// TODO: move these where used

// UpdateCacheFromAPI updates the anime titles cache from the AniDB titles dump.
func UpdateCacheFromAPI() error {
	ts, err := anidb.RequestTitles()
	if err != nil {
		return xerrors.Errorf("update cache from api: %w", err)
	}
	if err := titles.SaveDefault(ts); err != nil {
		return xerrors.Errorf("update cache from api: %w", err)
	}
	return nil
}

// UpdateCache updates the anime titles cache with an AniDB XML titles dump.
func UpdateCache(d []byte) error {
	ts, err := anidb.DecodeTitles(d)
	if err != nil {
		return xerrors.Errorf("update cache: %w", err)
	}
	if err := titles.SaveDefault(ts); err != nil {
		return xerrors.Errorf("update cache: %w", err)
	}
	return nil
}
