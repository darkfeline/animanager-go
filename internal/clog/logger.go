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

// Package clog implements attaching loggers to context.
package clog

import "context"

type keyType int

const (
	_ keyType = iota
	loggerKey
	slogKey
)

// A Logger can be used for logging.
// A Logger must be safe to use concurrently.
type Logger interface {
	Printf(string, ...any)
}

// WithLogger returns a new context with the logger attached.
func WithLogger(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// Printf prints to the context logger.
// If there is no context logger, the message is discarded.
func Printf(ctx context.Context, format string, a ...any) {
	getLogger(ctx).Printf(format, a...)
}

// getLogger returns the Logger attached to the context.
func getLogger(ctx context.Context) Logger {
	v := ctx.Value(loggerKey)
	if v == nil {
		return nullLogger{}
	}
	return v.(Logger)
}

type nullLogger struct{}

func (nullLogger) Printf(string, ...any) {}
