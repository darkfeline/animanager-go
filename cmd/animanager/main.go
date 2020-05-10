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
	"fmt"
	"os"

	"github.com/google/subcommands"

	"go.felesatra.moe/animanager/internal/config"

	"go.felesatra.moe/animanager/cmd/animanager/internal/cmd"
)

var configPath = flag.String("config", config.DefaultPath, "Config file")

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	cmd.AddCommands(subcommands.DefaultCommander)
	flag.Parse()
	ctx := context.Background()
	c, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %s\n", err)
	}
	os.Exit(int(subcommands.Execute(ctx, c)))
}
