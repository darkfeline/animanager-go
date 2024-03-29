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

// Package udp implements stuff for the AniDB UDP API.
package udp

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"go.felesatra.moe/anidb/udpapi"
	"go.felesatra.moe/animanager/internal/clientid"
)

type Client struct {
	*udpapi.Client
	userinfo udpapi.UserInfo
	logger   *slog.Logger
}

type Config struct {
	ServerAddr string
	UserInfo   udpapi.UserInfo
	Logger     *slog.Logger
}

func Dial(ctx context.Context, cfg *Config) (*Client, error) {
	c, err := udpapi.Dial(cfg.ServerAddr, cfg.Logger.WithGroup("udpapi"))
	if err != nil {
		return nil, err
	}
	c.ClientName = clientid.UDPName
	c.ClientVersion = clientid.UDPVersion
	c2 := &Client{
		Client:   c,
		userinfo: cfg.UserInfo,
		logger:   cfg.Logger,
	}
	if err := c2.login(ctx); err != nil {
		return nil, err
	}
	return c2, nil
}

// Shutdown the client.
//
// The underlying AniDB client connection is logged out and closed.
//
// The context should have enough time to allow the client to log out,
// especially when using encryption.
// Otherwise, you must wait for the encryption session to timeout
// before starting another server.
func (c *Client) Shutdown(ctx context.Context) error {
	if err := c.logout(ctx); err != nil {
		return fmt.Errorf("server shutdown: %s", err)
	}
	c.Client.Close()
	return nil
}

func (c *Client) login(ctx context.Context) error {
	log.Printf("Logging in to AniDB...")
	if c.userinfo.APIKey != "" {
		if err := c.Client.Encrypt(ctx, c.userinfo); err != nil {
			return fmt.Errorf("server login: %s", err)
		}
	}
	if _, err := c.Client.Auth(ctx, c.userinfo); err != nil {
		return fmt.Errorf("server login: %s", err)
	}
	log.Printf("Logged in to AniDB")
	return nil
}

func (c *Client) logout(ctx context.Context) error {
	log.Printf("Logging out of AniDB...")
	if err := c.Client.Logout(ctx); err != nil {
		return fmt.Errorf("server logout: %s", err)
	}
	log.Printf("Logged out of AniDB")
	return nil
}
