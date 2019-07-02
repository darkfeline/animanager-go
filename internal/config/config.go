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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"go.felesatra.moe/go2/errors"
)

// Config is the configuration for Animanager.
type Config struct {
	DBPath    string   `toml:"database"`
	WatchDirs []string `toml:"watch_dirs"`
	Player    []string `toml:"player"`
}

var defaultDir = filepath.Join(os.Getenv("HOME"), ".animanager")

// New loads the configuration file.  If an error occurs, an error is
// returned along with the default configuration.
func New(p string) (Config, error) {
	c := Default()
	f, err := os.Open(p)
	if err != nil {
		return c, errors.Wrap(err, "load config")
	}
	defer f.Close()
	d, err := ioutil.ReadAll(f)
	if err != nil {
		return c, errors.Wrapf(err, "load config %s", p)
	}
	if err := toml.Unmarshal(d, &c); err != nil {
		return c, errors.Wrapf(err, "load config %s", p)
	}
	c.DBPath = os.ExpandEnv(c.DBPath)
	for i, d := range c.WatchDirs {
		c.WatchDirs[i] = os.ExpandEnv(d)
	}
	return c, nil
}

// Default returns the default configuration.
func Default() Config {
	return Config{
		DBPath: filepath.Join(defaultDir, "database.db"),
		Player: []string{"mpv", "--quiet"},
	}
}
