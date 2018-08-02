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

// Package cmd implements subcommands.
package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
	"go.felesatra.moe/animanager/internal/config"
	"go.felesatra.moe/animanager/internal/database"
)

type Show struct {
}

func (*Show) Name() string     { return "show" }
func (*Show) Synopsis() string { return "Show information about a series." }
func (*Show) Usage() string {
	return `show AID:
  Show information about a series.
`
}

func (s *Show) SetFlags(f *flag.FlagSet) {
}

func (s *Show) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	c := config.New()
	_, err := database.Open(c.DBPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %s\n", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
