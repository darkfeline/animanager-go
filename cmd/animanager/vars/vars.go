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
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"sync"

	"go.felesatra.moe/animanager/internal/clog"
	"go.felesatra.moe/animanager/internal/config"
	"go.felesatra.moe/animanager/internal/database"
	"golang.org/x/sys/unix"
)

type ConfigVar struct {
	cfgPath string
	getCfg  func() (*config.Config, error)
}

func Config(fs *flag.FlagSet) *ConfigVar {
	v := &ConfigVar{}
	v.getCfg = sync.OnceValues(func() (*config.Config, error) { return config.Load(v.cfgPath) })
	fs.StringVar(&v.cfgPath, "config", config.DefaultPath, "Path to config file")
	return v
}

func (v ConfigVar) Load() (*config.Config, error) {
	cfg, err := v.getCfg()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %s", err)
	}
	return cfg, nil
}

func (v ConfigVar) OpenDB() (*sql.DB, error) {
	cfg, err := v.Load()
	if err != nil {
		return nil, err
	}
	return database.Open(context.Background(), cfg.DBPath)
}

type ContextVar struct {
}

func Context(fs *flag.FlagSet) *ContextVar {
	v := &ContextVar{}
	return v
}

func (v ContextVar) Context() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	ctx = clog.WithLogger(ctx, log.Default())
	return signal.NotifyContext(ctx, unix.SIGTERM, unix.SIGINT)
}
