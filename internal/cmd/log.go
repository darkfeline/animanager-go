package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

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
		fmt.Fprintln(os.Stderr, err)
	}
}
