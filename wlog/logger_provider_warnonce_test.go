// Copyright (c) 2018 Palantir Technologies. All rights reserved.
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
	"testing"

	"github.com/stretchr/testify/assert"
)

// Verifies that a logger created by newWarnOnceLoggerProvider outputs a warning on the first log call but not on
// subsequent calls.
func TestWarnOnceProvider(t *testing.T) {
	const wantOutput = `[WARNING] Logging operation that uses the default logger provider was performed without specifying a logger provider implementation.` + "\n" +
		`          To see logger output, set the global tracer implementation using wlog.SetDefaultLoggerProvider or by importing an implementation.` + "\n" +
		`          This warning can be disabled by setting the global logger provider to be the noop logger provider using wlog.SetDefaultLoggerProvider(wlog.NewNoopLoggerProvider()).` + "\n"
	provider := newWarnOnceLoggerProvider()

	buf := &bytes.Buffer{}
	logger := provider.NewLogger(buf)

	logger.Log()
	assert.Equal(t, wantOutput, buf.String())

	buf.Reset()
	logger.Log()
	assert.Equal(t, "", buf.String())
}
