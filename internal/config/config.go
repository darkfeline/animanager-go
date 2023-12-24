// Copyright (C) 2018  Allen Li
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

// Package config implements configuration for Animanager.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/BurntSushi/toml"
)

// Config is the configuration for Animanager.
type Config struct {
	DBPath       string      `toml:"database"`
	WatchDirs    []string    `toml:"watch_dirs"`
	Player       []string    `toml:"player"`
	AniDB        AniDBConfig `toml:"anidb"`
	ServerAddr   string      `toml:"server_addr"`
	FilePatterns []string    `toml:"file_patterns"`

	regexps func() ([]*regexp.Regexp, error)
}

func (c *Config) FileRegexps() ([]*regexp.Regexp, error) {
	if c.regexps == nil {
		c.regexps = sync.OnceValues(func() ([]*regexp.Regexp, error) {
			rs := make([]*regexp.Regexp, len(c.FilePatterns))
			for i, p := range c.FilePatterns {
				r, err := regexp.Compile(p)
				if err != nil {
					return nil, err
				}
				rs[i] = r
			}
			return rs, nil
		})
	}
	return c.regexps()
}

// AniDBConfig is the configuration for AniDB (mainly UDP API).
type AniDBConfig struct {
	UDPServerAddr string `toml:"udp_server_addr"`
	Username      string `toml:"username"`
	Password      string `toml:"password"`
	APIKey        string `toml:"api_key"`
}

// DefaultPath is the default config file path.
var DefaultPath string

func init() {
	d := os.Getenv("XDG_CONFIG_HOME")
	if d == "" {
		d = filepath.Join(os.Getenv("HOME"), ".config")
	}
	DefaultPath = filepath.Join(d, "animanager", "config.toml")
}

var defaultConfig = Config{
	Player:     []string{"mpv", "--quiet"},
	ServerAddr: "127.0.0.1:1234",
	AniDB: AniDBConfig{
		UDPServerAddr: "api.anidb.net:9000",
	},
}

func init() {
	d := os.Getenv("XDG_STATE_HOME")
	if d == "" {
		d = filepath.Join(os.Getenv("HOME"), ".local", "state")
	}
	defaultConfig.DBPath = filepath.Join(d, "animanager", "database.db")
}

// Load loads the configuration file.  If an error occurs, an error is
// returned along with the default configuration.
func Load(p string) (*Config, error) {
	// Copy default config.
	c := defaultConfig
	d, err := os.ReadFile(p)
	if err != nil {
		return &c, fmt.Errorf("load config: %s", err)
	}
	if err := toml.Unmarshal(d, &c); err != nil {
		return &c, fmt.Errorf("load config %s: %s", p, err)
	}
	c.DBPath = os.ExpandEnv(c.DBPath)
	for i, d := range c.WatchDirs {
		c.WatchDirs[i] = os.ExpandEnv(d)
	}
	return &c, nil
}
