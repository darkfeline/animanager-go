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
	"go.felesatra.moe/animanager/internal/clientid"
	"go.felesatra.moe/animanager/internal/server/api"
)

type Server struct {
	api.UnimplementedApiServer
	client   *udpapi.Client
	userinfo udpapi.UserInfo
	logger   Logger
}

type Config struct {
	UserInfo udpapi.UserInfo
	Logger   Logger
}

// NewServer starts a new server.
// You must call Shutdown, especially when using encryption.
// The context is used only for login.
func NewServer(ctx context.Context, cfg *Config) (*Server, error) {
	c, err := udpapi.NewClient()
	if err != nil {
		return nil, fmt.Errorf("new server: %s", err)
	}
	c.ClientName = clientid.Name
	c.ClientVersion = clientid.Version
	c.SetLogger(cfg.Logger)
	s := &Server{
		client:   c,
		userinfo: cfg.UserInfo,
		logger:   cfg.Logger,
	}
	return s, nil
}
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

// A Logger can be used for logging.
// A Logger must be safe to use concurrently.
type Logger interface {
	Printf(string, ...any)
}
