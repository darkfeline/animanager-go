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
	"fmt"
	"testing"
)

func TestPrefixLogger(t *testing.T) {
	t.Parallel()
	var s spyLogger
	p := prefixLogger{
		prefix: "mika:",
		logger: &s,
	}
	p.Printf("%s %s", "azusa", "hifumi")
	got := s.msg
	const want = "mika:azusa hifumi"
	if got != want {
		t.Errorf("got log message %q; want %q", got, want)
	}
}

func TestPrefixLogger_nested(t *testing.T) {
	t.Parallel()
	var s spyLogger
	p := prefixLogger{
		prefix: "outer:",
		logger: &s,
	}
	p = prefixLogger{
		prefix: "inner:",
		logger: p,
	}
	p.Printf("%s %s", "azusa", "hifumi")
	got := s.msg
	const want = "outer:inner:azusa hifumi"
	if got != want {
		t.Errorf("got log message %q; want %q", got, want)
	}
}

type spyLogger struct {
	msg string
}

func (l *spyLogger) Printf(format string, a ...any) {
	l.msg = fmt.Sprintf(format, a...)
}
