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

package evt2logtests

import (
	"bytes"
	"io"
	"testing"

	"github.com/palantir/pkg/objmatcher"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/witchcraft-go-logging/wlog/evtlog/evt2log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	Name         string
	EventName    string
	Values       map[string]interface{}
	UID          string
	SID          string
	TokenID      string
	OrgID        string
	UnsafeParams map[string]interface{}
	JSONMatcher  objmatcher.MapMatcher
}

func (tc TestCase) Params() []evt2log.Param {
	return []evt2log.Param{
		evt2log.Values(tc.Values),
		evt2log.UID(tc.UID),
		evt2log.SID(tc.SID),
		evt2log.Tag("tagName", "tagValue"),
		evt2log.TokenID(tc.TokenID),
		evt2log.OrgID(tc.OrgID),
		evt2log.UnsafeParams(tc.UnsafeParams),
	}
}

func TestCases() []TestCase {
	return []TestCase{
		{
			Name:      "basic event log entry",
			EventName: "com.palantir.foundry.build.buildstarted",
			UID:       "user-1",
			SID:       "session-1",
			Values: map[string]interface{}{
				"dataset": "my-cool-dataset",
			},
			TokenID: "X-Y-Z",
			OrgID:   "org-1",
			UnsafeParams: map[string]interface{}{
				"Password": "HelloWorld!",
			},
			JSONMatcher: map[string]objmatcher.Matcher{
				"type":      objmatcher.NewEqualsMatcher("event.2"),
				"eventName": objmatcher.NewEqualsMatcher("com.palantir.foundry.build.buildstarted"),
				"time":      objmatcher.NewRegExpMatcher(".+"),
				"uid":       objmatcher.NewEqualsMatcher("user-1"),
				"sid":       objmatcher.NewEqualsMatcher("session-1"),
				"values": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"dataset": objmatcher.NewEqualsMatcher("my-cool-dataset"),
				}),
				"tags": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"tagName": objmatcher.NewEqualsMatcher("tagValue"),
				}),
				"tokenId": objmatcher.NewEqualsMatcher("X-Y-Z"),
				"orgId":   objmatcher.NewEqualsMatcher("org-1"),
				"unsafeParams": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"Password": objmatcher.NewEqualsMatcher("HelloWorld!"),
				}),
			},
		},
	}
}

func JSONTestSuite(t *testing.T, loggerProvider func(w io.Writer) evt2log.Logger) {
	jsonOutputTests(t, loggerProvider)
	valueIsntOverwrittenByValues(t, loggerProvider)
	extraValuesIndependentAcrossCalls(t, loggerProvider)
}

func jsonOutputTests(t *testing.T, loggerProvider func(w io.Writer) evt2log.Logger) {
	for i, tc := range TestCases() {
		t.Run(tc.Name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := loggerProvider(buf)

			logger.Event(tc.EventName, tc.Params()...)

			gotEventLog := map[string]interface{}{}
			logEntry := buf.Bytes()
			err := safejson.Unmarshal(logEntry, &gotEventLog)
			require.NoError(t, err, "Case %d: %s\nEvent log line is not a valid map: %v", i, tc.Name, string(logEntry))

			assert.NoError(t, tc.JSONMatcher.Matches(gotEventLog), "Case %d: %s", i, tc.Name)
		})
	}
}

// Verifies that if different parameters are specified using Value and Values params, all of the values are present in
// the final output (that is, these parameters should be additive).
func valueIsntOverwrittenByValues(t *testing.T, loggerProvider func(w io.Writer) evt2log.Logger) {
	t.Run("Value and Values params are additive", func(t *testing.T) {
		var buf bytes.Buffer
		logger := loggerProvider(&buf)

		logger.Event("event", evt2log.Value("key", "value"), evt2log.Values(map[string]interface{}{"keys": "values"}))

		gotEventLog := map[string]interface{}{}
		logEntry := buf.Bytes()
		err := safejson.Unmarshal(logEntry, &gotEventLog)
		require.NoError(t, err, "Event log line is not a valid map: %v", string(logEntry))

		assert.NoError(t, objmatcher.MapMatcher(map[string]objmatcher.Matcher{
			"eventName": objmatcher.NewEqualsMatcher("event"),
			"time":      objmatcher.NewRegExpMatcher(".+"),
			"type":      objmatcher.NewEqualsMatcher("event.2"),
			"values": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"key":  objmatcher.NewEqualsMatcher("value"),
				"keys": objmatcher.NewEqualsMatcher("values"),
			}),
		}).Matches(gotEventLog))
	})
}

// Verifies that parameters remain separate between different logger calls (ensures there is not a bug where parameters
// are modified by making a logger call).
func extraValuesIndependentAcrossCalls(t *testing.T, loggerProvider func(w io.Writer) evt2log.Logger) {
	t.Run("Value and Values params stay separate across logger calls", func(t *testing.T) {
		var buf bytes.Buffer
		logger := loggerProvider(&buf)

		reusedParams := evt2log.Values(map[string]interface{}{"keys": "values"})
		logger.Event("event", reusedParams, evt2log.Value("key", "value"))
		gotEventLog := map[string]interface{}{}
		logEntry := buf.Bytes()
		err := safejson.Unmarshal(logEntry, &gotEventLog)
		require.NoError(t, err, "Event log line is not a valid map: %v", string(logEntry))

		assert.NoError(t, objmatcher.MapMatcher(map[string]objmatcher.Matcher{
			"eventName": objmatcher.NewEqualsMatcher("event"),
			"time":      objmatcher.NewRegExpMatcher(".+"),
			"type":      objmatcher.NewEqualsMatcher("event.2"),
			"values": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"key":  objmatcher.NewEqualsMatcher("value"),
				"keys": objmatcher.NewEqualsMatcher("values"),
			}),
		}).Matches(gotEventLog))

		buf.Reset()
		logger.Event("event", reusedParams)

		gotEventLog = map[string]interface{}{}
		logEntry = buf.Bytes()
		err = safejson.Unmarshal(logEntry, &gotEventLog)
		require.NoError(t, err, "Event log line is not a valid map: %v", string(logEntry))

		assert.NoError(t, objmatcher.MapMatcher(map[string]objmatcher.Matcher{
			"eventName": objmatcher.NewEqualsMatcher("event"),
			"time":      objmatcher.NewRegExpMatcher(".+"),
			"type":      objmatcher.NewEqualsMatcher("event.2"),
			"values": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"keys": objmatcher.NewEqualsMatcher("values"),
			}),
		}).Matches(gotEventLog))
	})
}
