package cmd

import (
	"io/ioutil"
	"log"
)

// Logger is used by this package for logging.
var Logger = log.New(ioutil.Discard, "cmd: ", log.LstdFlags)
