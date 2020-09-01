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

package wlog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/completion"
	"github.com/palantir/pkg/cli/flag"
	"github.com/palantir/witchcraft-go-logging/wlog-tmpl/logentryformatter"
	"github.com/palantir/witchcraft-go-logging/wlog-tmpl/logs"
	"github.com/pkg/errors"
)

const (
	inputFlagName   = "input"
	fileFlagName    = "file"
	strictFlagName  = "strict"
	colorFlagName   = "color"
	tmplFlagName    = "template"
	onlyFlagName    = "only"
	excludeFlagName = "exclude"
	noSubFlagName   = "no-substitution"
)

func Command() cli.Command {
	return cli.Command{
		Name:        "wlog",
		Usage:       "Converts witchcraft-formatted logs to a common, human-readable format on stdout",
		Description: getDescription(logs.OrderedLogTypes(), logs.Formatters()),
		Flags: []flag.Flag{
			flag.StringFlag{
				Name:  inputFlagName,
				Alias: "i",
				Usage: "The witchcraft log line to convert. If this flag is absent, lines are read from stdin.",
			},
			flag.StringFlag{
				Name:  fileFlagName,
				Alias: "f",
				Usage: "The witchcraft log file to convert. If this flag is absent, lines are read from stdin.",
			},
			flag.BoolFlag{
				Name:  strictFlagName,
				Usage: "If set, an error will be printed for any lines that cannot be parsed. If not set, lines that cannot be parsed will be printed as received.",
			},
			flag.BoolFlag{
				Name:  colorFlagName,
				Usage: "Set to true to colorize output based on log level or false to disable color output. If not set, then output is colorized when sent to a terminal and is not otherwise.",
			},
			flag.BoolFlag{
				Name:  noSubFlagName,
				Usage: "if specified, then no substitution will be performed on the log output even if the log type defines a substitution function (for example, substituting path parameters from param maps into the path in request logs)",
			},
			flag.StringFlag{
				Name:  tmplFlagName,
				Usage: "specify custom template to use for a log type. Specified in the form 'type:template'. Flag can be specified multiple times to define templates for multiple log types.",
			},
			flag.StringFlag{
				Name:  onlyFlagName,
				Usage: "if specified, then only logs of this type will be outputted. Flag can be specified multiple times to specify multiple log types.",
			},
			flag.StringFlag{
				Name:  excludeFlagName,
				Usage: "if specified, then logs of this type will not be outputted. Flag can be specified multiple times to specify multiple log types.",
			},
		},
		Action: func(ctx cli.Context) error {
			if ctx.Has(inputFlagName) && ctx.Has(fileFlagName) {
				return fmt.Errorf("%s and %s cannot both be specified", inputFlagName, fileFlagName)
			}

			if ctx.Has(colorFlagName) {
				color.NoColor = !ctx.Bool(colorFlagName)
			}

			var params []logentryformatter.Param
			if ctx.Has(noSubFlagName) {
				params = append(params, logentryformatter.NoSubstitution())
			}

			formatters := logs.Formatters(params...)
			if ctx.Has(tmplFlagName) {
				for _, curr := range ctx.StringSlice(tmplFlagName) {
					colonIdx := strings.Index(curr, ":")
					if colonIdx == -1 {
						return fmt.Errorf("value of %s flag must be of the form 'type:template', but flag did not have ':'", tmplFlagName)
					}
					typ := logentryformatter.LogType(curr[:colonIdx])
					tmpl := curr[colonIdx+1:]

					formatter, err := logs.Formatter(typ, tmpl, params...)
					if err != nil {
						return errors.Wrapf(err, "failed to set formatter for type %v to use template %s", typ, tmpl)
					}
					formatters[typ] = formatter
				}
			}

			var reader io.Reader = os.Stdin
			switch {
			case ctx.Has(inputFlagName):
				reader = strings.NewReader(ctx.String(inputFlagName))
			case ctx.Has(fileFlagName):
				f, err := os.Open(ctx.String(fileFlagName))
				if err != nil {
					return errors.Wrapf(err, "failed to open file")
				}
				defer func() {
					// file is opened for reading, so nothing to be done on error closing
					_ = f.Close()
				}()
				reader = f
			}
			return Convert(reader, ctx.App.Stdout,
				Strict(ctx.Bool(strictFlagName)),
				Formatters(formatters),
				OnlyString(ctx.StringSlice(onlyFlagName)...),
				ExcludeString(ctx.StringSlice(excludeFlagName)...),
			)
		},
	}
}

func Completions() map[string]completion.Provider {
	return map[string]completion.Provider{
		fileFlagName: completion.Filepath,
	}
}

func getDescription(orderedKeys []logentryformatter.LogType, formatters map[logentryformatter.LogType]logentryformatter.Formatter) string {
	buf := &bytes.Buffer{}
	fmt.Fprintln(buf, "Transforms input lines that contain witchcraft log entries. Each line must contain a single valid log entry.")
	if len(formatters) == 0 {
		return strings.TrimRight(buf.String(), "\n")
	}

	fmt.Fprintf(buf, "\t Provides built-in support for the following log types and uses the specified Go templates to print them:\n\n")

	fmt.Fprintf(buf, "\t Custom output templates can be specified for any log type using the '--%s' flag. For the built-in log types,\n", tmplFlagName)
	fmt.Fprintf(buf, "\t the templates use the struct definitions outlined below. For any other log type, the JSON log is read in as a generic\n")
	fmt.Fprintf(buf, "\t map[string]interface{} and this map is provided to the template (so accessors should use the JSON key names directly).\n")
	fmt.Fprintln(buf)

	fmt.Fprintf(buf, "\t The following is an example of how to specify a custom template for logs of type \"service.1\":\n")
	fmt.Fprintln(buf)

	fmt.Fprintf(buf, "\t wlog --template 'service.1:{{.Level}} ({{.Time}}): {{.Message}}'\n")
	fmt.Fprintln(buf)

	fmt.Fprintf(buf, "\t Object structures for built-in log types\n")
	fmt.Fprintf(buf, "\t ----------------------------------------\n")

	for _, k := range orderedKeys {
		fmt.Fprintf(buf, "   %s\n", k)
		fmt.Fprintf(buf, "   %s\n", strings.Repeat("-", len(k)))
		fmt.Fprintf(buf, "   Default output template:\n")
		fmt.Fprintf(buf, "   %s\n\n", formatters[k].RawTemplate())

		body := formatters[k].TemplateObjectDescription()
		parts := strings.Split(body, "\n")
		for i, curr := range parts {
			fmt.Fprintf(buf, "   %s", curr)
			if i != len(parts)-1 {
				fmt.Fprintln(buf)
			}
		}
		fmt.Fprintf(buf, "\n")
		fmt.Fprintf(buf, "   %s\n\n", strings.Repeat("-", len(k)))
	}
	return strings.TrimRight(buf.String(), "\n")
}
