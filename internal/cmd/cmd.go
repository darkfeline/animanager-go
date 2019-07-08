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

import "github.com/google/subcommands"

func AddCommands(c *subcommands.Commander) {
	c.Register(wrap(&Add{}), "")
	c.Register(wrap(&FindFiles{}), "")
	c.Register(wrap(&Register{}), "")
	c.Register(wrap(&Show{}), "")
	c.Register(wrap(&ShowFiles{}), "")
	c.Register(wrap(&Search{}), "")
	c.Register(wrap(&SetDone{}), "")
	c.Register(wrap(&Stats{}), "")
	c.Register(wrap2(&UpdateTitles{}), "")
	c.Register(wrap(&Unfinished{}), "")
	c.Register(wrap2(&Unregister{}), "")
	c.Register(wrap2(&Watch{}), "")
	c.Register(wrap2(&Watchable{}), "")
}
