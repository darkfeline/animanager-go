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

// Package anidb provides AniDB client functions for Animanager.
package anidb

import (
	"time"

	"go.felesatra.moe/anidb"
	"golang.org/x/time/rate"
)

// Client is the AniDB client for Animanager.
var Client = anidb.Client{
	Name:    "kfanimanager",
	Version: 1,
	Limiter: rate.NewLimiter(rate.Every(2*time.Second), 1),
}

// RequestAnime calls anidb.RequestAnime with the Animanager Client.
func RequestAnime(aid int) (*anidb.Anime, error) {
	return Client.RequestAnime(aid)
}
