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
	"os"
	"path/filepath"
)

type Config struct {
	configDir string
}

func (c *Config) DBPath() string {
	return filepath.Join(c.configDir, "database.db")
}

func New() *Config {
	home := os.Getenv("HOME")
	return &Config{
		configDir: filepath.Join(home, ".animanager"),
	}
}
