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
	"fmt"

	"go.felesatra.moe/anidb/udpapi"
	"go.felesatra.moe/animanager/internal/server/api"
)

type Server struct {
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

func (*Server) Ping(req api.PingRequest, resp *api.PingResponse) error {
	resp.Message = req.Message
	return nil
}

func (*Server) Login(req api.LoginRequest, resp *api.LoginResponse) error {
	return nil
}

func (*Server) Logout(req api.LogoutRequest, resp *api.LogoutResponse) error {
	return nil
}

func (*Server) File(req api.FileRequest, resp *api.FileResponse) error {
	return nil
}
