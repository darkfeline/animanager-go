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

package vars

import (
	"context"
	"flag"
	"os/signal"

	"golang.org/x/sys/unix"
)

type ContextVar struct {
}

func Context(fs *flag.FlagSet) *ContextVar {
	v := &ContextVar{}
	return v
}

func (v ContextVar) Context() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	return signal.NotifyContext(ctx, unix.SIGTERM, unix.SIGINT)
}
