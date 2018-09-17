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
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/subcommands"

	"go.felesatra.moe/animanager/internal/anidb"
	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/query"
)

type Add struct {
}

func (*Add) Name() string     { return "add" }
func (*Add) Synopsis() string { return "Add an anime." }
func (*Add) Usage() string {
	return `Usage: add aids...
Add an anime.
`
}

func (*Add) SetFlags(f *flag.FlagSet) {
}

func (a *Add) Execute(ctx context.Context, f *flag.FlagSet, x ...interface{}) subcommands.ExitStatus {
	if f.NArg() < 1 {
		fmt.Fprint(os.Stderr, a.Usage())
		return subcommands.ExitUsageError
	}
	aids := make([]int, f.NArg())
	for i, s := range f.Args() {
		aid, err := strconv.Atoi(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid AID: %s\n", err)
			return subcommands.ExitUsageError
		}
		aids[i] = aid
	}

	c := getConfig(x)
	db, err := database.Open(ctx, c.DBPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %s\n", err)
		return subcommands.ExitFailure
	}
	defer db.Close()
	for i, aid := range aids {
		if err := addAnime(db, aid); err != nil {
			fmt.Fprintf(os.Stderr, "Error adding anime: %s\n", err)
			return subcommands.ExitFailure
		}
		if i < len(aids)-1 {
			time.Sleep(2 * time.Second)
		}
	}
	return subcommands.ExitSuccess
}

func addAnime(db *sql.DB, aid int) error {
	a, err := anidb.RequestAnime(aid)
	if err != nil {
		return err
	}
	if err := query.InsertAnime(db, a); err != nil {
		return err
	}
	return nil
}
