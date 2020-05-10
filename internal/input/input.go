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

// Package input provides input reading utilities.
package input

import (
	"errors"
	"fmt"
	"strings"
)

// Reader defines the interface required by input reading functions.
// This is usually a bufio.Reader around os.Stdin.
type Reader interface {
	ReadString(byte) (string, error)
}

// ReadLine reads one line of input.
func ReadLine(r Reader) (string, error) {
	return r.ReadString('\n')
}

// ErrInvalid is returned for invalid input.
var ErrInvalid = errors.New("invalid input")

// ReadYN reads a yes or no input.  The provided default value is
// returned for empty inputs.  If an invalid input is provided, the
// returned error is ErrInvalid.
func ReadYN(r Reader, def bool) (bool, error) {
	line, err := ReadLine(r)
	if err != nil {
		return def, err
	}
	s := strings.TrimSpace(line)
	s = strings.ToLower(s)
	switch s {
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	case "":
		return def, nil
	default:
		return def, fmt.Errorf("%w %s", ErrInvalid, s)
	}
}
