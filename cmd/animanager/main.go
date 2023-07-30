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
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"go.felesatra.moe/animanager/internal/config"
)

// BUG(): This flag doesn't work
var configPath = flag.String("config", config.DefaultPath, "Config file")

func main() {
	log.SetPrefix("animanager: ")
	if len(os.Args) < 2 {
		printUsage(os.Stdout)
		os.Exit(0)
	}
	cmd, args := os.Args[1], os.Args[2:]
	for _, c := range commands {
		if cmd != c.name() {
			continue
		}
		cfg, err := config.Load(*configPath)
		if err != nil {
			log.Printf("error loading config: %s\n", err)
		}
		if err := c.run(&c, cfg, args); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				os.Exit(0)
			}
			log.Fatal(err)
		}
		os.Exit(0)
	}
	log.Printf("unknown command %s", cmd)
	printUsage(os.Stderr)
	os.Exit(2)
}

func printUsage(w io.Writer) error {
	bw := bufio.NewWriter(w)
	fmt.Fprintf(bw, `Usage: animanager <command> [arguments]

The commands are:

`)
	for _, c := range commands {
		fmt.Fprintf(bw, "\t%s\t%s\n", c.name(), c.shortDesc)
	}
	fmt.Fprintf(bw, `
Use "animanager help <command>" for more information about a command.
`)
	return bw.Flush()
}

var commands = []command{
	addCmd,
	findFilesCmd,
	registerCmd,
	showCmd,
	showFilesCmd,
	searchCmd,
	serverCmd,
	setDoneCmd,
	statsCmd,
	updateTitlesCmd,
	unfinishedCmd,
	unregisterCmd,
	watchCmd,
	watchableCmd,
}

type command struct {
	usageLine string
	shortDesc string
	longDesc  string
	run       func(*command, *config.Config, []string) error
}

func (c *command) name() string {
	return strings.SplitN(c.usageLine, " ", 2)[0]
}

func (c *command) flagSet() *flag.FlagSet {
	fs := flag.NewFlagSet(c.name(), flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: animanager %s\n\n", c.usageLine)
		fmt.Fprintf(fs.Output(), "%s\n", c.longDesc)
		fs.PrintDefaults()
	}
	return fs
}
