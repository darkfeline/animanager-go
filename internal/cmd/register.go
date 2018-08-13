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
	"os"
	"regexp"
	"strconv"

	"github.com/google/subcommands"

	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/query"
)

type Register struct {
	pattern string
}

func (*Register) Name() string     { return "register" }
func (*Register) Synopsis() string { return "Register an anime." }
func (*Register) Usage() string {
	return `Usage: register [-pattern pattern] aid
Register an anime.
`
}

func (r *Register) SetFlags(f *flag.FlagSet) {
	f.StringVar(&r.pattern, "pattern", "", "File pattern")
}

func (r *Register) Execute(ctx context.Context, f *flag.FlagSet, x ...interface{}) subcommands.ExitStatus {
	if f.NArg() != 1 {
		fmt.Fprint(os.Stderr, r.Usage())
		return subcommands.ExitUsageError
	}
	aid, err := strconv.Atoi(f.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid AID: %s\n", err)
		return subcommands.ExitUsageError
	}

	c := getConfig(x)
	db, err := database.Open(ctx, c.DBPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %s\n", err)
		return subcommands.ExitFailure
	}
	defer db.Close()
	if r.pattern == "" {
		a, err := query.GetAnime(db, aid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return subcommands.ExitFailure
		}
		r.pattern = animeDefaultRegexp(a)
	}
	if err := query.InsertWatching(db, aid, r.pattern); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

var nonAlphaNum = regexp.MustCompile("[^a-zA-Z0-9]+")

// animeDefaultRegexp returns the default regexp watching pattern for the
// anime.
func animeDefaultRegexp(a *query.Anime) (re string) {
	var b bytes.Buffer
	fragments := nonAlphaNum.Split(a.Title, -1)
	for _, s := range fragments {
		io.WriteString(&b, regexp.QuoteMeta(s))
		io.WriteString(&b, `.*?`)
	}
	io.WriteString(&b, `\b([0-9]+)`)
	return b.String()
}
