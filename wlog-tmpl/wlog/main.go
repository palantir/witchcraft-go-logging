// Copyright (c) 2020 Palantir Technologies. All rights reserved.
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

package main

import (
	"os"

	"github.com/nmiyake/pkg/errorstringer"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/witchcraft-go-logging/wlog-tmpl/wlog/cmd/wlog"
)

var version = "unspecified"

func main() {
	os.Exit(createApp().Run(os.Args))
}

func createApp() *cli.App {
	app := cli.NewApp(cli.DebugHandler(errorstringer.StackWithInterleavedMessages))
	app.Name = "wlog"
	app.Usage = "A simple CLI for interacting with witchcraft log files."
	app.Version = version
	app.Completion = wlog.Completions()

	cmd := wlog.Command()
	cmd.Flags = append(app.Flags, cmd.Flags...)
	app.Command = cmd

	return app
}
