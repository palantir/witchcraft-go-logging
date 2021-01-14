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

package logs_test

import (
	"fmt"
	"testing"

	"github.com/palantir/witchcraft-go-logging/wlog-tmpl/logentryformatter"
	"github.com/palantir/witchcraft-go-logging/wlog-tmpl/logs"
	"github.com/stretchr/testify/assert"
)

type LogTest struct {
	name   string
	input  []string
	output []string
	err    error
}

func RunLogTests(t *testing.T, tests []LogTest) {
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d_%s", i, test.name), test.Test)
	}
}

func (test LogTest) Test(t *testing.T) {
	for inputIdx := range test.input {
		line, err := logentryformatter.FormatLogLine(test.input[inputIdx], logs.Unwrappers, logs.Formatters(), nil, nil)
		if test.err != nil {
			assert.EqualError(t, err, test.err.Error(), "unexpected error on line %d", inputIdx+1)
		} else {
			assert.NoError(t, err, "error on line %d", inputIdx+1)
		}
		expected := test.output[inputIdx]
		assert.Equal(t, expected, line, "unexpected output on line %d", inputIdx+1)
	}
}
