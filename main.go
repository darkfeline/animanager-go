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

// Command animanager manages watched anime and anime to be watched.
package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"

	"go.felesatra.moe/animanager/internal/anidb/titles"
	"go.felesatra.moe/animanager/internal/cmd"
	"go.felesatra.moe/animanager/internal/config"
	"go.felesatra.moe/animanager/internal/migrate"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&cmd.Add{}, "")
	subcommands.Register(&cmd.Show{}, "")
	subcommands.Register(&cmd.TitleSearch{}, "")
	flag.Parse()
	setupLog(debug)
	ctx := context.Background()
	c := config.New()
	os.Exit(int(subcommands.Execute(ctx, c)))
}

func setupLog(debug bool) {
	if !debug {
		return
	}
	migrate.Logger.SetOutput(os.Stderr)
	titles.Logger.SetOutput(os.Stderr)
}
