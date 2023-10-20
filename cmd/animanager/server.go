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
	"net"
	"os/signal"

	"go.felesatra.moe/animanager/internal/server"
	"go.felesatra.moe/animanager/internal/server/api"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc"
)

var serverCmd = command{
	usageLine: "server",
	shortDesc: "Run AniDB UDP API server",
	longDesc: `Run AniDB UDP API server.
Used internally to maintain a UDP session for reuse across commands.
`,
	run: func(cmd *command, args []string) error {
		f := cmd.flagSet()
		if err := f.Parse(args); err != nil {
			return err
		}

		ctx := context.Background()
		s, err := server.NewServer(ctx, &server.Config{})
		if err != nil {
			return err
		}
		defer s.Shutdown(context.Background())

		ctx, stop := signal.NotifyContext(ctx, unix.SIGTERM, unix.SIGINT)
		defer stop()

		rs := grpc.NewServer()
		api.RegisterApiServer(rs, s)
		l, err := net.Listen("tcp", ":1234")
		if err != nil {
			return err
		}
		go func() {
			select {
			case <-ctx.Done():
				rs.Stop()
			}
		}()
		return rs.Serve(l)
	},
}
