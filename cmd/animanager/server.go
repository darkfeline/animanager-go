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
	"time"

	"go.felesatra.moe/animanager/cmd/animanager/vars"
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

EXPERIMENTAL; DO NOT USE
`,
	run: func(h *handle, args []string) error {
		f := h.flagSet()
		cfgv := vars.Config(f)
		if err := f.Parse(args); err != nil {
			return err
		}
		cfg, err := cfgv.Load()
		if err != nil {
			return err
		}

		ctx := context.Background()
		ctx, cancel := signal.NotifyContext(ctx, unix.SIGTERM, unix.SIGINT)
		defer cancel()
		c, err := cfgv.DialUDP(ctx)
		if err != nil {
			return err
		}
		s := server.NewServer(ctx, c)
		defer func(ctx context.Context) {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			s.Shutdown(ctx)
		}(context.WithoutCancel(ctx))

		rs := grpc.NewServer()
		api.RegisterApiServer(rs, s)

		l, err := net.Listen("tcp", cfg.ServerAddr)
		if err != nil {
			return err
		}
		go func() {
			<-ctx.Done()
			rs.Stop()
		}()
		return rs.Serve(l)
	},
}
