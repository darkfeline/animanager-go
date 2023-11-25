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
	"os/signal"
	"time"

	"go.felesatra.moe/animanager/internal/clog"
	"go.felesatra.moe/animanager/internal/udp"
	"golang.org/x/sys/unix"
)

var findFilesUDPCmd = command{
	usageLine: "findfilesudp",
	shortDesc: "find episode files (UDP)",
	longDesc: `Find episode files using the UDP API.

EXPERIMENTAL; DO NOT USE
`,
	run: func(cmd *command, args []string) error {
		f := cmd.flagSet()
		if err := f.Parse(args); err != nil {
			return err
		}
		if f.NArg() != 0 {
			return errors.New("no arguments allowed")
		}
		cfg, err := cmd.loadConfig()
		if err != nil {
			return err
		}

		log.Printf("Finding video files...")
		files, err := findVideoFilesMany(cfg.WatchDirs)
		if err != nil {
			return err
		}
		log.Printf("Finished finding video files")

		db, err := openDB(cfg)
		if err != nil {
			return err
		}
		defer db.Close()

		ctx := context.Background()
		ctx = clog.WithLogger(ctx, log.Default())
		ctx, stop := signal.NotifyContext(ctx, unix.SIGTERM, unix.SIGINT)
		defer stop()

		c, err := udp.Dial(ctx, &udp.Config{
			ServerAddr: "api.anidb.net:9000",
			UserInfo:   userInfo(cfg),
			Logger:     log.Default(),
		})
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
	return nil
}
