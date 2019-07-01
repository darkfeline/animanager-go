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
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/google/subcommands"

	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/obx"
	"go.felesatra.moe/animanager/internal/query"
)

type Unregister struct {
	complete bool
}

func (*Unregister) Name() string     { return "unregister" }
func (*Unregister) Synopsis() string { return "Unregister an anime." }
func (*Unregister) Usage() string {
	return `Usage: unregister aid...
       unregister -complete [aids...]
Unregister an anime.
`
}

func (u *Unregister) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&u.complete, "complete", false, "Unregister complete anime")
}

func (u *Unregister) Execute(ctx context.Context, f *flag.FlagSet, x ...interface{}) subcommands.ExitStatus {
	if f.NArg() < 1 && !u.complete {
		fmt.Fprint(os.Stderr, u.Usage())
		return subcommands.ExitUsageError
	}
	aids := make([]int, f.NArg())
	for _, s := range f.Args() {
		aid, err := strconv.Atoi(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid AID: %s\n", err)
			return subcommands.ExitUsageError
		}
		aids = append(aids, aid)
	}

	c := getConfig(x)
	db, err := database.Open(ctx, c.DBPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %s\n", err)
		return subcommands.ExitFailure
	}
	defer db.Close()
	if u.complete {
		as, err := obx.GetCompleteAnime(db)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return subcommands.ExitFailure
		}
		aids = append(aids, as...)
	}
	for _, aid := range aids {
		fmt.Println(aid)
		if err := query.DeleteWatching(db, aid); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return subcommands.ExitFailure
		}
	}
	return subcommands.ExitSuccess
}
