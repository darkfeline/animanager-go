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
	"fmt"
	"net/rpc"

	"go.felesatra.moe/animanager/internal/server/api"
)

var callCmd = command{
	usageLine: "call",
	shortDesc: "Call method on AniDB UDP API server",
	longDesc: `Run AniDB UDP API server.
Used for testing.
`,
	run: func(cmd *command, args []string) error {
		f := cmd.flagSet()
		if err := f.Parse(args); err != nil {
			return err
		}

		c, err := rpc.Dial("tcp", ":1234")
		if err != nil {
			return err
		}
		resp := api.PingResponse{
			Message: "vanitas",
		}
		if err := c.Call("API.Ping", api.PingRequest{}, &resp); err != nil {
			return err
		}
		fmt.Printf("%s\n", resp.Message)
		return nil
	},
}
