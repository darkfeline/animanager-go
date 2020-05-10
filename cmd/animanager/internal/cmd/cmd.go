// Copyright (C) 2019  Allen Li
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/subcommands"
	"go.felesatra.moe/animanager/internal/config"
)

func AddCommands(c *subcommands.Commander) {
	c.Register(wrap(&Add{}), "")
	c.Register(wrap(&FindFiles{}), "")
	c.Register(wrap(&Register{}), "")
	c.Register(wrap(&Show{}), "")
	c.Register(wrap(&ShowFiles{}), "")
	c.Register(wrap(&Search{}), "")
	c.Register(wrap(&SetDone{}), "")
	c.Register(wrap(&Stats{}), "")
	c.Register(wrap(&UpdateTitles{}), "")
	c.Register(wrap(&Unfinished{}), "")
	c.Register(wrap(&Unregister{}), "")
	c.Register(wrap(&Watch{}), "")
	c.Register(wrap(&Watchable{}), "")
}

type command interface {
	Name() string
	Synopsis() string
	Usage() string
	SetFlags(*flag.FlagSet)
	Run(context.Context, *flag.FlagSet, config.Config) error
}

func wrap(c command) subcommands.Command {
	return wrapper{c}
}

type wrapper struct {
	command
}

func (w wrapper) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	cfg := args[0].(config.Config)
	if err := w.command.Run(ctx, f, cfg); err != nil {
		switch err.(type) {
		case usageError:
			fmt.Fprintf(os.Stderr, w.command.Usage())
			return subcommands.ExitUsageError
		default:
			log.Printf("Error: %s", err)
			return subcommands.ExitFailure
		}
	}
	return subcommands.ExitSuccess
}

type usageError struct {
	message string
}

func (e usageError) Error() string {
	return fmt.Sprintf("usage error: %s", e.message)
}
