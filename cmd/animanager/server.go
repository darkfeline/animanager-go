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
	"log"
	"net"
	"time"

	"go.felesatra.moe/anidb/udpapi"
	"go.felesatra.moe/animanager/cmd/animanager/vars"
	"go.felesatra.moe/animanager/internal/clog"
	"go.felesatra.moe/animanager/internal/config"
	"go.felesatra.moe/animanager/internal/server"
	"go.felesatra.moe/animanager/internal/server/api"
	"go.felesatra.moe/animanager/internal/udp"
	"google.golang.org/grpc"
)

var serverCmd = command{
	usageLine: "server",
	shortDesc: "Run AniDB UDP API server",
	longDesc: `Run AniDB UDP API server.

Used internally to maintain a UDP session for reuse across commands.

EXPERIMENTAL; DO NOT USE
`,
	run: func(cmd *command, args []string) error {
		f := cmd.flagSet()
		ctxv := vars.Context(f)
		cfgv := vars.Config(f)
		if err := f.Parse(args); err != nil {
			return err
		}
		cfg, err := cfgv.Load()
		if err != nil {
			return err
		}

		ctx, cancel := ctxv.Context()
		defer cancel()
		s, err := server.NewServer(ctx, &udp.Config{
			ServerAddr: cfg.UDPServerAddr,
			UserInfo:   userInfo(cfg),
			Logger:     log.Default(),
		})
		if err != nil {
			return err
		}
		defer func(ctx context.Context) {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			s.Shutdown(ctx)
		}(context.WithoutCancel(ctx))

		rs := grpc.NewServer(grpc.UnaryInterceptor(withLogger{log.Default()}.Unary))
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

type withLogger struct {
	logger clog.Logger
}

func (w withLogger) Unary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	ctx = clog.WithLogger(ctx, w.logger)
	return handler(ctx, req)
}

func userInfo(cfg *config.Config) udpapi.UserInfo {
	return udpapi.UserInfo{
		UserName:     cfg.AniDB.Username,
		UserPassword: cfg.AniDB.Password,
		APIKey:       cfg.AniDB.APIKey,
	}
}
