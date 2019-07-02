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
	"bufio"
	"context"
	"flag"
	"os"
	"sort"

	"github.com/google/subcommands"

	"go.felesatra.moe/animanager/internal/afmt"
	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/query"
)

type Unfinished struct {
}

func (*Unfinished) Name() string     { return "unfinished" }
func (*Unfinished) Synopsis() string { return "Print unfinished anime." }
func (*Unfinished) Usage() string {
	return `Usage: unfinished
Print unfinished anime.
`
}

func (s *Unfinished) SetFlags(f *flag.FlagSet) {
}

func (s *Unfinished) Execute(ctx context.Context, f *flag.FlagSet, x ...interface{}) subcommands.ExitStatus {
	return executeInner(s, ctx, f, x)
}

func (s *Unfinished) innerExecute(ctx context.Context, f *flag.FlagSet, x ...interface{}) error {
	c := getConfig(x)
	db, err := database.Open(ctx, c.DBPath)
	if err != nil {
		return err
	}
	defer db.Close()
	as, err := query.GetUnwatchedAnime(db)
	bw := bufio.NewWriter(os.Stdout)
	sort.Slice(as, func(i, j int) bool { return as[i].AID < as[j].AID })
	for _, a := range as {
		afmt.PrintAnimeShort(bw, &a)
	}
	bw.Flush()
	return nil
}
