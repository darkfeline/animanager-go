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

package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/google/subcommands"
	_ "github.com/mattn/go-sqlite3"
	"go.felesatra.moe/animanager/internal/migrate"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&cliCmd{}, "")
	flag.Parse()
	setupLog(debug)
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}

var ilog = log.New(os.Stderr, "", log.LstdFlags)
var dlog = log.New(ioutil.Discard, "", log.LstdFlags)

func setupLog(debug bool) {
	if !debug {
		return
	}
	dlog.SetOutput(os.Stderr)
	migrate.Logger.SetOutput(os.Stderr)
	migrate.Logger.SetPrefix("migrate: ")
}
