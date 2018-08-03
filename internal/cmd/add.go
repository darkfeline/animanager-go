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
	"log"
	"os"
	"strconv"

	"github.com/google/subcommands"
	"go.felesatra.moe/animanager/internal/config"
	"go.felesatra.moe/animanager/internal/database"
)

type Add struct {
}

func (*Add) Name() string     { return "add" }
func (*Add) Synopsis() string { return "Add an anime." }
func (*Add) Usage() string {
	return `add AID:
  Add an anime.
`
}

func (*Add) SetFlags(f *flag.FlagSet) {
}

func (a *Add) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() != 1 {
		fmt.Fprint(os.Stderr, a.Usage())
		return subcommands.ExitUsageError
	}
	aid, err := strconv.Atoi(f.Arg(0))
	log.Print(aid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid AID: %s\n", err)
		return subcommands.ExitUsageError
	}
	c := config.New()

	db, err := database.Open(c.DBPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %s\n", err)
		return subcommands.ExitFailure
	}
	defer db.Close()
	return subcommands.ExitSuccess
}
