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
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"go.felesatra.moe/animanager/internal/config"
	"go.felesatra.moe/go2/errors"
)

// Logger is used by this package for logging.
var Logger = log.New(ioutil.Discard, "cmd: ", log.LstdFlags)

type userError interface {
	error
	UserError() string
}

// PrintError is used by this package for printing user facing errors.
var PrintError func(error) = func(err error) {
	var err2 userError
	if errors.AsValue(err2, err) {
		fmt.Fprintln(os.Stderr, err2.UserError())
	} else {
		fmt.Fprintln(os.Stderr, errors.Format(err, false))
	}
}

// getConfig gets the Config passed into a subcommand.
func getConfig(x []interface{}) config.Config {
	return x[0].(config.Config)
}
