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

package wlog_test

import (
	"bytes"
	"testing"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/stretchr/testify/assert"
)

// Verifies that a logger created using the default logger provider uses the newWarnOnceLoggerProvider.
//
// Note that this test depends on the global state of defaultLoggerProvider: if the default logger is set before this
// test is run in the same process or if the default logger provider is changed in the future, the test will fail.
func TestDefaultProviderIsWarnOnceProvider(t *testing.T) {
	// create logger that writes to buffer using the default logger provider
	buf := &bytes.Buffer{}
	logger := svc1log.New(buf, wlog.DebugLevel) // uses default provider

	// verify that output provides warning that no logger provider was specified
	logger.Info("Test output 1")
	const wantOutput = `[WARNING] Logging operation that uses the default logger provider was performed without specifying a logger provider implementation. To see logger output, set the global logger provider implementation using wlog.SetDefaultLoggerProvider or by importing an implementation. This warning can be disabled by setting the global logger provider to be the noop logger provider using wlog.SetDefaultLoggerProvider(wlog.NewNoopLoggerProvider()).` + "\n"
	got := buf.String()
	assert.Equal(t, wantOutput, got)

	// verify that warning is only written on first call to logger
	logger.Info("Test output 2")
	buf.Reset()
	got = buf.String()
	assert.Equal(t, "", got)
}
