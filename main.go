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
	"path/filepath"

	"github.com/google/subcommands"

	"go.felesatra.moe/animanager/internal/cmd"
	"go.felesatra.moe/animanager/internal/config"
	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/migrate"
)

var defaultConfig = filepath.Join(os.Getenv("HOME"), ".animanager", "config.toml")

func main() {
	var debug bool
	var configPath string
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.StringVar(&configPath, "config", defaultConfig, "Config file")
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(cmd.Wrap(&cmd.Add{}), "")
	subcommands.Register(cmd.Wrap(&cmd.FindFiles{}), "")
	subcommands.Register(cmd.Wrap(&cmd.Register{}), "")
	subcommands.Register(cmd.Wrap(&cmd.Show{}), "")
	subcommands.Register(cmd.Wrap(&cmd.ShowFiles{}), "")
	subcommands.Register(cmd.Wrap(&cmd.Search{}), "")
	subcommands.Register(cmd.Wrap(&cmd.SetDone{}), "")
	subcommands.Register(cmd.Wrap(&cmd.Stats{}), "")
	subcommands.Register(cmd.Wrap(&cmd.UpdateTitles{}), "")
	subcommands.Register(cmd.Wrap(&cmd.Unfinished{}), "")
	subcommands.Register(cmd.Wrap(&cmd.Unregister{}), "")
	subcommands.Register(cmd.Wrap(&cmd.Watch{}), "")
	subcommands.Register(cmd.Wrap(&cmd.Watchable{}), "")
	flag.Parse()
	if debug {
		setupDebug()
	}
	ctx := context.Background()
	c, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %s\n", err)
	}
	os.Exit(int(subcommands.Execute(ctx, c)))
}

func setupDebug() {
	cmd.Logger.SetOutput(os.Stderr)
	database.Logger.SetOutput(os.Stderr)
	migrate.Logger.SetOutput(os.Stderr)
}
