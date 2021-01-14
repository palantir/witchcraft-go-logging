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

package logentryformatter_test

import (
	"encoding/json"
	"testing"

	"github.com/palantir/witchcraft-go-logging/wlog-tmpl/logentryformatter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatLogLine(t *testing.T) {
	for i, currCase := range []struct {
		name          string
		input         string
		unwrappersMap map[logentryformatter.LogType]logentryformatter.Unwrapper
		formattersMap map[logentryformatter.LogType]logentryformatter.Formatter
		only          map[logentryformatter.LogType]struct{}
		exclude       map[logentryformatter.LogType]struct{}
		wantOutput    string
		wantErr       string
	}{
		{
			name:    "non-JSON input",
			input:   `zzzz`,
			wantErr: `Log line "zzzz" is not valid JSON`,
		},
		{
			name:    "JSON input with no type key",
			input:   `{"key":"val"}`,
			wantErr: `Log line JSON "{\"key\":\"val\"}" does not have a "type" key so its log type cannot be determined`,
		},
		{
			name:    "formatter matches on type",
			input:   `{"type":"myType","key":"val"}`,
			wantErr: `Skipping unknown log line type: myType`,
		},
		{
			name:          "formatter matches on type",
			input:         `{"type":"myType","key":"val"}`,
			formattersMap: testMyTypeFormatter,
			wantOutput:    `Output of key: val`,
		},
		{
			name:  "no output if log type does not match",
			input: `{"type":"myType","key":"val"}`,
			only: map[logentryformatter.LogType]struct{}{
				logentryformatter.LogType("myType-2"): {},
			},
		},
		{
			name:  "no output if log type excluded",
			input: `{"type":"myType","key":"val"}`,
			exclude: map[logentryformatter.LogType]struct{}{
				logentryformatter.LogType("myType"): {},
			},
		},
		{
			name:          "output if log type matched",
			input:         `{"type":"myType","key":"val"}`,
			formattersMap: testMyTypeFormatter,
			only: map[logentryformatter.LogType]struct{}{
				logentryformatter.LogType("myType"): {},
			},
			wantOutput: "Output of key: val",
		},
		{
			name:  "no output if log type matched and excluded",
			input: `{"type":"myType","key":"val"}`,
			only: map[logentryformatter.LogType]struct{}{
				logentryformatter.LogType("myType"): {},
			},
			exclude: map[logentryformatter.LogType]struct{}{
				logentryformatter.LogType("myType"): {},
			},
		},
		{
			name:          "Log lines with non-JSON prefix (e.g., from docker-compose logs)",
			input:         `build2.palantir.dev       | {"type":"myType","key":"val"}`,
			formattersMap: testMyTypeFormatter,
			wantOutput:    "Output of key: val",
		},
	} {
		output, err := logentryformatter.FormatLogLine(currCase.input, currCase.unwrappersMap, currCase.formattersMap, currCase.only, currCase.exclude)
		if currCase.wantErr == "" {
			require.NoError(t, err, "Case %d: %s", i, currCase.name)
			assert.Equal(t, currCase.wantOutput, output, "Case %d: %s", i, currCase.name)
		} else {
			assert.EqualError(t, err, currCase.wantErr, "Case %d: %s", i, currCase.name)
		}
	}
}

var testMyTypeFormatter = map[logentryformatter.LogType]logentryformatter.Formatter{
	logentryformatter.LogType("myType"): mustNew(func(lineJSON []byte, substitute bool) (interface{}, error) {
		var m map[string]interface{}
		err := json.Unmarshal(lineJSON, &m)
		return m, err
	}, `Output of key: {{.key}}`),
}

func mustNew(entryParser func([]byte, bool) (interface{}, error), tmplString string) logentryformatter.Formatter {
	f, err := logentryformatter.New(entryParser, tmplString)
	if err != nil {
		panic(err)
	}
	return f
}
