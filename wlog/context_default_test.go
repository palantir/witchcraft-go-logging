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
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit2log"
	"github.com/palantir/witchcraft-go-logging/wlog/evtlog/evt2log"
	"github.com/palantir/witchcraft-go-logging/wlog/metriclog/metric1log"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests the behavior of logging output using the logger retrieved from a context.Context when no logger was set on the
// context. Each logger type has its own functions for retrieving the logger from the context, so this test is defined
// as a set of tests run on for each logger that supports being stored in and retrieved from a context.
//
// There are 3 cases that are tested:
//   * No default logger is set for the logger type and no global logger provider is set
//   * No default logger is set for the logger type and a JSON-writing global logger provider is set
//   * A no-op default logger is set for the logger type
//
// See the comments on the test cases in testFromContextFromEmptyContextForSingleLogger for more details on the tests
// and the expected behavior that is tested.
//
// Note that, because these tests are testing behavior based on global state, they should not be run in parallel.
func TestOutputFromContextEmptyContext(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	// set os.Stderr to be the temporary file and defer restoration
	prevStderr := os.Stderr
	defer func() {
		os.Stderr = prevStderr
	}()

	for _, tc := range []loggerTestCase{
		{
			loggerPkg: "audit2log",
			performLogging: func() {
				logger := audit2log.FromContext(context.Background())
				logger.Audit("TEST_EVT", audit2log.AuditResultSuccess)
			},
			validateJSON: func(bytes []byte) {
				var logEntry logging.AuditLogV2
				require.NoError(t, json.Unmarshal(bytes, &logEntry))
				assert.Equal(t, "TEST_EVT", logEntry.Name)
			},
			setEmptyLoggerCreator: func() {
				// set the default logger creator
				audit2log.SetDefaultLoggerCreator(func() audit2log.Logger {
					// technically, a truly no-op writer would implement the logger interface and do nothing, but it is
					// easier to just return a new logger that writes to a no-op writer.
					return audit2log.New(ioutil.Discard)
				})
			},
		},
		{
			loggerPkg: "evt2log",
			performLogging: func() {
				logger := evt2log.FromContext(context.Background())
				logger.Event("TEST_EVT")
			},
			validateJSON: func(bytes []byte) {
				var logEntry logging.EventLogV2
				require.NoError(t, json.Unmarshal(bytes, &logEntry))
				assert.Equal(t, "TEST_EVT", logEntry.EventName)
			},
			setEmptyLoggerCreator: func() {
				// set the default logger creator
				evt2log.SetDefaultLoggerCreator(func() evt2log.Logger {
					// technically, a truly no-op writer would implement the logger interface and do nothing, but it is
					// easier to just return a new logger that writes to a no-op writer.
					return evt2log.New(ioutil.Discard)
				})
			},
		},
		{
			loggerPkg: "metric1log",
			performLogging: func() {
				logger := metric1log.FromContext(context.Background())
				logger.Metric("com.palantir.metric", "gauge")
			},
			validateJSON: func(bytes []byte) {
				var logEntry logging.MetricLogV1
				require.NoError(t, json.Unmarshal(bytes, &logEntry))
				assert.Equal(t, "com.palantir.metric", logEntry.MetricName)
				assert.Equal(t, "gauge", logEntry.MetricType)
			},
			setEmptyLoggerCreator: func() {
				// set the default logger creator
				metric1log.SetDefaultLoggerCreator(func() metric1log.Logger {
					// technically, a truly no-op writer would implement the logger interface and do nothing, but it is
					// easier to just return a new logger that writes to a no-op writer.
					return metric1log.New(ioutil.Discard)
				})
			},
		},
		{
			loggerPkg: "svc1log",
			performLogging: func() {
				logger := svc1log.FromContext(context.Background())
				logger.Info("Test message")
			},
			validateJSON: func(bytes []byte) {
				var logEntry logging.ServiceLogV1
				require.NoError(t, json.Unmarshal(bytes, &logEntry))
				assert.Equal(t, "Test message", logEntry.Message)
				assert.Equal(t, logging.LogLevel("INFO"), logEntry.Level)
			},
			setEmptyLoggerCreator: func() {
				// set the default logger creator
				svc1log.SetDefaultLoggerCreator(func() svc1log.Logger {
					// technically, a truly no-op writer would implement the logger interface and do nothing, but it is
					// easier to just return a new logger that writes to a no-op writer.
					return svc1log.New(ioutil.Discard, wlog.InfoLevel)
				})
			},
		},
	} {
		t.Run(tc.loggerPkg, func(t *testing.T) {
			testFromContextFromEmptyContextForSingleLogger(t, tmpDir, tc)
		})
	}
}

type loggerTestCase struct {
	loggerPkg             string
	performLogging        func()
	validateJSON          func([]byte)
	setEmptyLoggerCreator func()
}

func testFromContextFromEmptyContextForSingleLogger(t *testing.T, tmpDir string, loggerTestCaseInfo loggerTestCase) {
	for _, tc := range []struct {
		name   string
		before func(loggerTestCase)
		verify func(loggerTestCase, []byte)
	}{
		// Because the context has no logger set on it, the defaultLoggerCreator is used. The defaultLoggerCreator has
		// not been set, so the default implementation writes a warning to os.Stderr followed by the logger output for
		// a logger created by using wlog.DefaultLoggerProvider. Because wlog.DefaultLoggerProvider has not been set,
		// the logger created by that function outputs a warning stating that the global logger should be set. As a
		// result, the over-all logger output to os.Stderr is:
		// "[WARNING] <logger not set in context>: [WARNING] <global logger provider not set>"
		{
			name: "Context with no logger set and no logger provider set returns default stderr warning logger that uses warn-once logger",
			verify: func(loggerTestCaseInfo loggerTestCase, bytes []byte) {
				logOutput := string(bytes)

				firstPortionRegexp, err := regexp.Compile(
					regexp.QuoteMeta(`[WARNING]`) + ".*" + regexp.QuoteMeta(`github.com/`) + ".+" + regexp.QuoteMeta(`/witchcraft-go-logging/wlog_test.TestOutputFromContextEmptyContext`) + ".+" + regexp.QuoteMeta(`/github.com/`) + ".+" + regexp.QuoteMeta(`/witchcraft-go-logging/wlog/context_default_test.go:`) + "[0-9]+" + regexp.QuoteMeta(`]: usage of `+loggerTestCaseInfo.loggerPkg+`.Logger from FromContext that did not have that logger set: `))
				require.NoError(t, err, "Unexpected error compiling regex")
				loc := firstPortionRegexp.FindStringIndex(logOutput)
				require.NotNil(t, loc, "Unexpected log output: %s regex %s", logOutput, firstPortionRegexp.String())

				got := strings.TrimSuffix(logOutput[loc[1]:], "\n")
				assert.Equal(t, `[WARNING] Logging operation that uses the default logger provider was performed without specifying a logger provider implementation. To see logger output, set the global logger provider implementation using wlog.SetDefaultLoggerProvider or by importing an implementation. This warning can be disabled by setting the global logger provider to be the noop logger provider using wlog.SetDefaultLoggerProvider(wlog.NewNoopLoggerProvider()).`, got)
			},
		},
		// Because the context has no logger set on it, the defaultLoggerCreator is used. The defaultLoggerCreator has
		// not been set, so the default implementation writes a warning to os.Stderr followed by the logger output for
		// a logger created by using wlog.DefaultLoggerProvider. wlog.DefaultLoggerProvider has been set to
		// wlog.NewJSONMarshalLoggerProvider(), so the over-all logger output to os.Stderr is:
		// "[WARNING] <logger not set in context>: <JSON logger output>"
		{
			name: "Context with no logger set and no logger provider set returns default stderr warning logger that uses set logger provider",
			before: func(loggerTestCaseInfo loggerTestCase) {
				// set the default logger provider to be the JSON marshal logger provider
				wlog.SetDefaultLoggerProvider(wlog.NewJSONMarshalLoggerProvider())
			},
			verify: func(loggerTestCaseInfo loggerTestCase, bytes []byte) {
				logOutput := string(bytes)

				firstPortionRegexp := regexp.MustCompile(
					regexp.QuoteMeta(`[WARNING]`) + ".*" + regexp.QuoteMeta(`github.com/`) + ".+" + regexp.QuoteMeta(`/witchcraft-go-logging/wlog_test.TestOutputFromContextEmptyContext`) + ".+" + regexp.QuoteMeta(`/github.com/`) + ".+" + regexp.QuoteMeta(`/witchcraft-go-logging/wlog/context_default_test.go:`) + "[0-9]+" + regexp.QuoteMeta(`]: usage of `+loggerTestCaseInfo.loggerPkg+`.Logger from FromContext that did not have that logger set: `))
				loc := firstPortionRegexp.FindStringIndex(logOutput)
				require.NotNil(t, loc, "Unexpected log output: %s", logOutput)

				loggerTestCaseInfo.validateJSON([]byte(logOutput[loc[1]:]))
			},
		},
		// Because the context has no logger set on it, the defaultLoggerCreator is used. The defaultLoggerCreator has
		// been set to be a no-op one, so the result is that nothing is logged.
		{
			name: "Context with no logger set and no logger provider set returns default stderr warning logger that uses set logger provider",
			before: func(loggerTestCaseInfo loggerTestCase) {
				loggerTestCaseInfo.setEmptyLoggerCreator()
			},
			verify: func(loggerTestCaseInfo loggerTestCase, bytes []byte) {
				assert.Equal(t, "", string(bytes))
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// store and restore global default logger provider so that tests can assume that default one is set
			originalProvider := wlog.DefaultLoggerProvider()
			defer func() {
				wlog.SetDefaultLoggerProvider(originalProvider)
			}()

			f, err := ioutil.TempFile(tmpDir, "")
			require.NoError(t, err)
			defer func() {
				_ = f.Close()
			}()
			os.Stderr = f

			if tc.before != nil {
				tc.before(loggerTestCaseInfo)
			}

			// get logger from context and perform logging
			loggerTestCaseInfo.performLogging()

			err = f.Close()
			require.NoError(t, err)

			bytes, err := ioutil.ReadFile(f.Name())
			require.NoError(t, err)

			// verify logger output
			tc.verify(loggerTestCaseInfo, bytes)
		})
	}
}
