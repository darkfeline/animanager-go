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

package clog

import (
	"context"
	"io"
	"log/slog"
	"math"
)

// WithSlog returns a new context with the logger attached.
func WithSlog(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, slogKey, l)
}

// Slog returns the [log/slog.Logger] attached to the context.
// If there is no context logger, a no-op one is returned.
func Slog(ctx context.Context) *slog.Logger {
	v := ctx.Value(slogKey)
	if v == nil {
		return nullSlog
	}
	return v.(*slog.Logger)
}

var nullHandler = slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
	Level: slog.Level(math.MaxInt),
})
var nullSlog = slog.New(nullHandler)
