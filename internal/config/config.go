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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config is the configuration for Animanager.
type Config struct {
	DBPath        string      `toml:"database"`
	WatchDirs     []string    `toml:"watch_dirs"`
	Player        []string    `toml:"player"`
	AniDB         AniDBConfig `toml:"anidb"`
	ServerAddress string      `toml:"server_address"`
}

// AniDBConfig is the configuration for AniDB (mainly UDP API).
type AniDBConfig struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	APIKey   string `toml:"api_key"`
}

// DefaultPath is the default config file path.
var DefaultPath string

var defaultConfig = Config{
	Player: []string{"mpv", "--quiet"},
}

func init() {
	d := os.Getenv("XDG_CONFIG_HOME")
	if d == "" {
		d = filepath.Join(os.Getenv("HOME"), ".config")
	}
	DefaultPath = filepath.Join(d, "animanager", "config.toml")

	d = os.Getenv("XDG_STATE_HOME")
	if d == "" {
		d = filepath.Join(os.Getenv("HOME"), ".local", "share")
	}
	defaultConfig.DBPath = filepath.Join(d, "animanager", "database.db")
}

// Load loads the configuration file.  If an error occurs, an error is
// returned along with the default configuration.
func Load(p string) (*Config, error) {
	// Copy default config.
	c := defaultConfig
	d, err := ioutil.ReadFile(p)
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
