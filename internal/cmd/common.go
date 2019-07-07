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
	"io/ioutil"
	"log"

	"github.com/google/subcommands"
	"go.felesatra.moe/animanager/internal/config"
)

// Logger is used by this package for logging.
var Logger = log.New(ioutil.Discard, "cmd: ", log.LstdFlags)

// getConfig gets the Config passed into a subcommand.
func getConfig(x []interface{}) config.Config {
	return x[0].(config.Config)
}

type command interface {
	Name() string
	Synopsis() string
	Usage() string
	SetFlags(*flag.FlagSet)
	Run(context.Context, *flag.FlagSet, config.Config) error
}

type wrapper struct {
	command
}

func Wrap(c command) subcommands.Command {
	return wrapper{c}
}

func Wrap2(c subcommands.Command) subcommands.Command {
	return c
}

func (w wrapper) Execute(ctx context.Context, f *flag.FlagSet, x ...interface{}) subcommands.ExitStatus {
	cfg := getConfig(x)
	if err := w.command.Run(ctx, f, cfg); err != nil {
		log.Printf("Error: %s", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
