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

package cmd

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"go.felesatra.moe/animanager/internal/config"
	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/query"
)

type Register struct {
	pattern string
	offset  int
}

func (*Register) Name() string     { return "register" }
func (*Register) Synopsis() string { return "Register an anime." }
func (*Register) Usage() string {
	return `Usage: register [-pattern pattern] [-offset int] aid
Register an anime.
`
}

func (c *Register) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.pattern, "pattern", "", "File pattern")
	f.IntVar(&c.offset, "offset", 0, "Episode offset")
}

func (c *Register) Run(ctx context.Context, f *flag.FlagSet, cfg config.Config) error {
	if f.NArg() != 1 {
		return usageError{"must pass exactly one argument"}
	}
	aid, err := strconv.Atoi(f.Arg(0))
	if err != nil {
		return fmt.Errorf("invalid AID %v: %v", aid, err)
	}

	db, err := database.Open(ctx, cfg.DBPath)
	if err != nil {
		return err
	}
	defer db.Close()
	if c.pattern == "" {
		a, err := query.GetAnime(db, aid)
		if err != nil {
			return err
		}
		c.pattern = animeDefaultRegexp(a)
	}
	w := query.Watching{
		AID:    aid,
		Regexp: c.pattern,
		Offset: c.offset,
	}
	if err := query.InsertWatching(db, w); err != nil {
		return err
	}
	return nil
}

var nonAlphaNum = regexp.MustCompile("[^a-zA-Z0-9]+")

// animeDefaultRegexp returns the default regexp watching pattern for the
// anime.
func animeDefaultRegexp(a *query.Anime) (re string) {
	var b bytes.Buffer
	io.WriteString(&b, "(?i)")
	fragments := nonAlphaNum.Split(a.Title, -1)
	for _, s := range fragments {
		io.WriteString(&b, regexp.QuoteMeta(s))
		io.WriteString(&b, `.*?`)
	}
	io.WriteString(&b, `\b([0-9]+)`)
	return b.String()
}
