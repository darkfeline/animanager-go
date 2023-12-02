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
	"fmt"

	"go.felesatra.moe/animanager/cmd/animanager/vars"
	"go.felesatra.moe/animanager/internal/server/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var callCmd = command{
	usageLine: "call",
	shortDesc: "Call method on AniDB UDP API server",
	longDesc: `Run AniDB UDP API server.

Used for testing.

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

		conn, err := grpc.Dial(cfg.ServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
		defer conn.Close()
		c := api.NewApiClient(conn)
		ctx := context.Background()
		resp, err := c.Ping(ctx, &api.PingRequest{Message: "vanitas"})
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", resp.Message)
		return nil
	},
}
