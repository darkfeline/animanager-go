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

// Package server implements the internal server.
// Used to maintain AniDB UDP sessions.
package server

import (
	"context"
	"fmt"

	"go.felesatra.moe/animanager/internal/server/api"
	"go.felesatra.moe/animanager/internal/udp"
)

type Server struct {
	api.UnimplementedApiServer
	client *udp.Client
}

// NewServer starts a new server.
// You must call Shutdown, especially when using encryption.
// The context is used only for login.
func NewServer(ctx context.Context, cfg *udp.Config) (*Server, error) {
	c, err := udp.Dial(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("new server: %s", err)
	}
	s := &Server{
		client: c,
	}
	return s, nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.client.Shutdown(ctx)
}
