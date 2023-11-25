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

package server

import (
	"context"
	"errors"
	"time"

	"go.felesatra.moe/animanager/internal/clog"
)

type pingFunc func(context.Context) (port string, _ error)

// keepalive calls the ping function at intervals until canceled.
func keepalive(ctx context.Context, p pingFunc, d time.Duration) error {
	t := time.NewTicker(d)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			ctx, cancel := context.WithTimeoutCause(ctx, 2*time.Second, errors.New("keepalive ping timeout"))
			if _, err := p(ctx); err != nil {
				clog.Printf(ctx, "keepalive ping: %s", err)
			}
			cancel()
		case <-ctx.Done():
			return context.Cause(ctx)
		}
	}
}
