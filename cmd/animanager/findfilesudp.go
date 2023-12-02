// Copyright (C) 2023  Allen Li
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
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"os/signal"
	"time"

	"go.felesatra.moe/animanager/cmd/animanager/vars"
	"go.felesatra.moe/animanager/internal/fileid"
	"go.felesatra.moe/animanager/internal/udp"
	"golang.org/x/sys/unix"
)

var findFilesUDPCmd = command{
	usageLine: "findfilesudp",
	shortDesc: "find episode files (UDP)",
	longDesc: `Find episode files using the UDP API.

EXPERIMENTAL; DO NOT USE
`,
	run: func(h *handle, args []string) error {
		f := h.flagSet()
		cfgv := vars.Config(f)
		if err := f.Parse(args); err != nil {
			return err
		}
		if f.NArg() != 0 {
			return errors.New("no arguments allowed")
		}
		cfg, err := cfgv.Load()
		if err != nil {
			return err
		}

		log.Printf("Finding video files...")
		files, err := fileid.FindVideoFiles(cfg.WatchDirs)
		if err != nil {
			return err
		}
		log.Printf("Finished finding video files")

		db, err := cfgv.OpenDB()
		if err != nil {
			return err
		}
		defer db.Close()

		ctx := context.Background()
		ctx, cancel := signal.NotifyContext(ctx, unix.SIGTERM, unix.SIGINT)
		defer cancel()
		c, err := cfgv.DialUDP(ctx)
		if err != nil {
			return err
		}
		defer func(ctx context.Context) {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			c.Shutdown(ctx)
		}(context.WithoutCancel(ctx))

		if err := refreshFilesUDP(ctx, db, c, files); err != nil {
			return err
		}
		return nil
	},
}

// refreshFilesUDP updates episode files using the given video file
// paths.
func refreshFilesUDP(ctx context.Context, db *sql.DB, c *udp.Client, files []string) error {
	l := slog.Default().With("method", "udp")
	for _, f := range files {
		if err := fileid.MatchEpisode(ctx, db, c, f); err != nil {
			l.Warn("match file to episode", "file", f, "error", err)
		}
	}
	return nil
}
