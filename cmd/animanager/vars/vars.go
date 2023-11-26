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

// Package vars defines common flag variables.
package vars

import (
	"flag"
	"fmt"

	"go.felesatra.moe/animanager/internal/config"
)

type ConfigVar struct {
	cfgPath string
}

func Config(fs *flag.FlagSet) *ConfigVar {
	v := &ConfigVar{}
	fs.StringVar(&v.cfgPath, "config", config.DefaultPath, "Path to config file")
	return v
}

func (v ConfigVar) Load() (*config.Config, error) {
	cfg, err := config.Load(v.cfgPath)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %s", err)
	}
	return cfg, nil
}
