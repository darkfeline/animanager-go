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

package cmd

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/google/subcommands"
	"go.felesatra.moe/animanager/internal/config"
	"go.felesatra.moe/go2/errors"
	"golang.org/x/xerrors"
)

// Logger is used by this package for logging.
var Logger = log.New(ioutil.Discard, "cmd: ", log.LstdFlags)

// PrintError is used by this package for printing user facing errors.
var PrintError func(error) = func(err error) {
	var err2 userError
	if xerrors.As(err, &err2) {
		fmt.Fprintln(os.Stderr, err2.UserError())
	} else {
		fmt.Fprintln(os.Stderr, errors.Format(err, false))
	}
}

// userError is the interface implemented by errors to provide a
// user-friendly error string used by the default PrintError function.
type userError interface {
	error
	UserError() string
}

// getConfig gets the Config passed into a subcommand.
func getConfig(x []interface{}) config.Config {
	return x[0].(config.Config)
}

func executeInner(e innerExecutor, ctx context.Context, f *flag.FlagSet, x []interface{}) subcommands.ExitStatus {
	err := e.innerExecute(ctx, f, x...)
	var err2 usageError
	if xerrors.As(err, &err2) {
		PrintError(err)
		return subcommands.ExitUsageError
	}
	if err != nil {
		PrintError(err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

type innerExecutor interface {
	innerExecute(ctx context.Context, f *flag.FlagSet, x ...interface{}) error
}

type usageError struct {
	error
}
