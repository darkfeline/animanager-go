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

	"go.felesatra.moe/anidb/udpapi"
	"go.felesatra.moe/animanager/internal/server/api"
)

type Server struct {
	api.UnimplementedApiServer
	Client *udpapi.Client
}

func NewServer() (*Server, error) {
	c, err := udpapi.NewClient(&udpapi.ClientConfig{})
	if err != nil {
		return nil, fmt.Errorf("new server: %s", err)
	}
	return &Server{
		Client: c,
	}, nil
}

// Shutdown the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}

func (*Server) Ping(ctx context.Context, req *api.PingRequest) (*api.PingResponse, error) {
	return &api.PingResponse{
		Message: req.GetMessage(),
	}, nil
}
