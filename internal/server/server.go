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
	"log"

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
	ServerAddr string
	UserInfo   udpapi.UserInfo
	Logger     Logger
}

// A Logger can be used for logging.
// A Logger must be safe to use concurrently.
type Logger interface {
	Printf(string, ...any)
}

// NewServer starts a new server.
// You must call Shutdown, especially when using encryption.
// The context is used only for login.
func NewServer(ctx context.Context, cfg *Config) (*Server, error) {
	c, err := udpapi.NewClient(cfg.ServerAddr)
	if err != nil {
		return nil, fmt.Errorf("new server: %s", err)
	}
	c.ClientName = clientid.UDPName
	c.ClientVersion = clientid.UDPVersion
	c.SetLogger(cfg.Logger)
	s := &Server{
		client:   c,
		userinfo: cfg.UserInfo,
		logger:   cfg.Logger,
	}
	if err := s.login(ctx); err != nil {
		return nil, err
	}
	return s, nil
}

// Shutdown the server.
// The underlying AniDB client connection is logged out and closed.
// The context should have enough time to allow the client to log out,
// especially when using encryption.
// Otherwise, you must wait for the encryption session to timeout
// before starting another server.
// No new requests will be accepted (as the connection is closed).
// Outstanding requests will be unblocked.
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.logout(ctx); err != nil {
		return fmt.Errorf("server shutdown: %s", err)
	}
	s.client.Close()
	return nil
}

func (s *Server) login(ctx context.Context) error {
	s.logger.Printf("Logging in to AniDB...")
	if s.userinfo.APIKey != "" {
		if err := s.client.Encrypt(ctx, s.userinfo); err != nil {
			return fmt.Errorf("server login: %s", err)
		}
	}
	if _, err := s.client.Auth(ctx, s.userinfo); err != nil {
		return fmt.Errorf("server login: %s", err)
	}
	log.Printf("Logged in to AniDB")
	return nil
}

func (s *Server) logout(ctx context.Context) error {
	s.logger.Printf("Logging out of AniDB...")
	if err := s.client.Logout(ctx); err != nil {
		return fmt.Errorf("server logout: %s", err)
	}
	s.logger.Printf("Logged out of AniDB")
	return nil
}
