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
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/palantir/witchcraft-go-logging/wlog-tmpl/logentryformatter"
	"github.com/palantir/witchcraft-go-logging/wlog-tmpl/logs"
	"github.com/pkg/errors"
)

// Convert scans input from the provided reader and converts any lines that can be parsed as witchcraft log entries
// into a human-readable form and writes the result to the provided writer using the provided parameters.
func Convert(in io.Reader, out io.Writer, params ...Param) error {
	reader := bufio.NewReader(in)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return errors.Wrapf(err, "failed while reading input")
		}
		if len(line) > 0 {
			// ReadString includes delimiter in returned value, so trim it if present
			writeFormattedLine(out, strings.TrimSuffix(line, "\n"), getConvertParams(params...))
		}
		if err == io.EOF {
			break
		}
	}
	return nil
}

func writeFormattedLine(out io.Writer, text string, params convertParams) {
	output, err := logentryformatter.FormatLogLine(text, logs.Unwrappers, params.formatters, params.only, params.exclude)
	switch {
	case err == nil && output == "":
		// no error and no output: skip output
		return
	case err != nil && params.strict:
		// error was encountered and strict mode is true: output is error
		output = err.Error()
	case err != nil:
		// error was encountered and strict mode is false: output is raw input
		output = text
	}
	_, _ = fmt.Fprintln(out, output)
}
