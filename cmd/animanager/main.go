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
	"log/slog"
	"os"
	"strings"
	"sync"
)

func main() {
	log.SetPrefix("animanager: ")
	log.SetFlags(log.LstdFlags | log.Lmsgprefix)
	if len(os.Args) < 2 {
		printUsage(os.Stdout)
		os.Exit(0)
	}
	cmd, args := os.Args[1], os.Args[2:]
	for _, c := range commands {
		if cmd != c.name() {
			continue
		}
		if err := c.run(newHandle(&c), args); err != nil {
			// ErrHelp is returned when -h is passed and
			// the command FlagSet doesn't define it.
			// FlagSet.Parse will still print help, so we
			// only need to ignore the error.
			if errors.Is(err, flag.ErrHelp) {
				os.Exit(0)
			}
			log.Fatal(err)
		}
		os.Exit(0)
	}
	log.Fatalf("Unknown command %q", cmd)
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
	callCmd,
	clearFilesCmd,
	findFilesCmd,
	findFilesUDPCmd,
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
	run       func(*handle, []string) error
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

// A handle wraps a [command] and holds state while running the
// command.
// This lets us not store state on the command struct directly.
type handle struct {
	cmd     *command
	flagSet func() *flagSet

	// Flags
	verbose bool
}

func newHandle(cmd *command) *handle {
	h := &handle{
		cmd: cmd,
	}
	h.flagSet = sync.OnceValue(func() *flagSet {
		f := cmd.flagSet()
		h.setupFlags(f)
		f2 := &flagSet{
			FlagSet:   f,
			postParse: h.postFlagParse,
		}
		return f2
	})
	return h
}

func (h *handle) setupFlags(f *flag.FlagSet) {
	f.BoolVar(&h.verbose, "verbose", false, "Enable verbose logging")
}

func (h *handle) postFlagParse() {
	h.configureLogging()
}

func (h *handle) configureLogging() {
	if !h.verbose {
		return
	}
	th := slog.NewTextHandler(log.Default().Writer(), &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})
	slog.SetDefault(slog.New(th))
}

type flagSet struct {
	*flag.FlagSet
	postParse func()
}

func (f *flagSet) Parse(args []string) error {
	if err := f.FlagSet.Parse(args); err != nil {
		return err
	}
	f.postParse()
	return nil
}
