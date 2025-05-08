// Copyright (c) 2025 Palantir Technologies. All rights reserved.
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

package wlogtests

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/palantir/pkg/objmatcher"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type EntryTestCase struct {
	Name string

	EntryFn func(entry wlog.LogEntry)
	Matcher objmatcher.MapMatcher
}

func JSONTestSuite(t *testing.T, loggerProvider func(w io.Writer) wlog.Logger) {
	// Verifies behavior of basic log entry operations
	runEntryTestCases(t, "basic tests", basicTestCases(), loggerProvider)

	// Verifies that if different parameters are specified using Value and Values params, all of the values are present
	// in the final output (that is, these parameters should be additive)
	runEntryTestCases(t, "set same key tests", setSameKeyTestCases(), loggerProvider)
}

func basicTestCases() []EntryTestCase {
	return []EntryTestCase{
		{
			Name: "OptionalStringValue that is non-empty sets value",
			EntryFn: func(entry wlog.LogEntry) {
				entry.OptionalStringValue("test-key", "test-optional-string-value")
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher("test-optional-string-value"),
			},
		},
		{
			Name: "OptionalStringValue that is empty is a noop",
			EntryFn: func(entry wlog.LogEntry) {
				entry.OptionalStringValue("test-key", "")
			},
			Matcher: map[string]objmatcher.Matcher{},
		},
		{
			Name: "StringList writes array values",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringListValue("test-key", []string{"one", "two"})
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.NewEqualsMatcher("one"),
					objmatcher.NewEqualsMatcher("two"),
				}),
			},
		},
		{
			Name: "StringList that is empty writes empty array",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringListValue("test-key", []string{})
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher([]any{}),
			},
		},
		{
			Name: "StringList that is nil writes empty array",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringListValue("test-key", nil)
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher([]any{}),
			},
		},
		{
			Name: "StringListAppendValue writes array values",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringListAppendValue("test-key", []string{"one", "two"})
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.NewEqualsMatcher("one"),
					objmatcher.NewEqualsMatcher("two"),
				}),
			},
		},
		{
			Name: "StringListAppendValue appends array values",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringListAppendValue("test-key", []string{"one", "two"})
				entry.StringListAppendValue("test-key", []string{"three", "four"})
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.NewEqualsMatcher("one"),
					objmatcher.NewEqualsMatcher("two"),
					objmatcher.NewEqualsMatcher("three"),
					objmatcher.NewEqualsMatcher("four"),
				}),
			},
		},
		{
			Name: "StringListAppendValue that is empty writes empty array",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringListAppendValue("test-key", []string{})
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher([]any{}),
			},
		},
		{
			Name: "StringListAppendValue that is nil writes empty array",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringListAppendValue("test-key", nil)
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher([]any{}),
			},
		},
		{
			Name: "StringMapValue writes map values",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringMapValue("test-key", map[string]string{
					"one": "two",
				})
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.MapMatcher{
					"one": objmatcher.NewEqualsMatcher("two"),
				},
			},
		},
		{
			Name: "StringMapValue that is empty writes empty map",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringMapValue("test-key", map[string]string{})
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher(map[string]any{}),
			},
		},
		{
			Name: "StringMapValue that is nil writes empty map",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringMapValue("test-key", nil)
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher(map[string]any{}),
			},
		},
		{
			Name: "AnyMapValue writes map values",
			EntryFn: func(entry wlog.LogEntry) {
				entry.AnyMapValue("test-key", map[string]any{
					"one": 2,
				})
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.MapMatcher{
					"one": objmatcher.NewEqualsMatcher(json.Number("2")),
				},
			},
		},
		{
			Name: "AnyMapValue that is empty writes empty map",
			EntryFn: func(entry wlog.LogEntry) {
				entry.AnyMapValue("test-key", map[string]any{})
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher(map[string]any{}),
			},
		},
		{
			Name: "AnyMapValue that is nil writes empty map",
			EntryFn: func(entry wlog.LogEntry) {
				entry.AnyMapValue("test-key", nil)
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher(map[string]any{}),
			},
		},
		{
			Name: "ObjectList writes array values",
			EntryFn: func(entry wlog.LogEntry) {
				entry.ObjectListValue("test-key", []any{"one", 2, struct{ Three int }{Three: 3}})
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.NewEqualsMatcher("one"),
					objmatcher.NewEqualsMatcher(json.Number("2")),
					objmatcher.NewEqualsMatcher(map[string]any{"Three": json.Number("3")}),
				}),
			},
		},
		{
			Name: "ObjectList that is empty writes empty array",
			EntryFn: func(entry wlog.LogEntry) {
				entry.ObjectListValue("test-key", []any{})
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher([]any{}),
			},
		},
		{
			Name: "ObjectList that is nil writes empty array",
			EntryFn: func(entry wlog.LogEntry) {
				entry.ObjectListValue("test-key", nil)
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher([]any{}),
			},
		},
		{
			Name: "ObjectListAppendValue writes array values",
			EntryFn: func(entry wlog.LogEntry) {
				entry.ObjectListAppendValue("test-key", []any{"one", 2})
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.NewEqualsMatcher("one"),
					objmatcher.NewEqualsMatcher(json.Number("2")),
				}),
			},
		},
		{
			Name: "ObjectListAppendValue appends array values",
			EntryFn: func(entry wlog.LogEntry) {
				entry.ObjectListAppendValue("test-key", []any{"one", 2})
				entry.ObjectListAppendValue("test-key", []any{struct{ Three int }{Three: 3}, "four"})
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.NewEqualsMatcher("one"),
					objmatcher.NewEqualsMatcher(json.Number("2")),
					objmatcher.NewEqualsMatcher(map[string]any{"Three": json.Number("3")}),
					objmatcher.NewEqualsMatcher("four"),
				}),
			},
		},
		{
			Name: "ObjectListAppendValue that is empty writes empty array",
			EntryFn: func(entry wlog.LogEntry) {
				entry.ObjectListAppendValue("test-key", []any{})
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher([]any{}),
			},
		},
		{
			Name: "ObjectListAppendValue that is nil writes empty array",
			EntryFn: func(entry wlog.LogEntry) {
				entry.ObjectListAppendValue("test-key", nil)
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher([]any{}),
			},
		},
	}
}

func setSameKeyTestCases() []EntryTestCase {
	return []EntryTestCase{
		{
			Name: "IntValue overrides StringValue",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringValue("test-key", "test-string-value")

				entry.IntValue("test-key", 1)
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher(json.Number("1")),
			},
		},
		{
			Name: "SafeLongValue overrides StringValue",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringValue("test-key", "test-string-value")

				entry.SafeLongValue("test-key", 1)
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher(json.Number("1")),
			},
		},
		{
			Name: "StringValue overrides IntValue",
			EntryFn: func(entry wlog.LogEntry) {
				entry.IntValue("test-key", 1)

				entry.StringValue("test-key", "test-string-value")
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher("test-string-value"),
			},
		},
		{
			Name: "StringValue overrides StringMapValue",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringMapValue("test-key", map[string]string{
					"string-map-key-1": "string-map-value-1",
					"string-map-key-2": "string-map-value-2",
				})

				entry.StringValue("test-key", "test-string-value")
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher("test-string-value"),
			},
		},
		{
			Name: "StringValue overrides AnyMapValue",
			EntryFn: func(entry wlog.LogEntry) {
				entry.AnyMapValue("test-key", map[string]any{
					"any-map-key-1": "any-map-value-1",
					"any-map-key-2": 2,
				})

				entry.StringValue("test-key", "test-string-value")
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher("test-string-value"),
			},
		},
		{
			Name: "OptionalStringValue overrides StringValue",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringValue("test-key", "test-string-value")

				entry.OptionalStringValue("test-key", "test-optional-string-value")
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher("test-optional-string-value"),
			},
		},
		{
			Name: "OptionalStringValue that is empty deletes key",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringValue("test-key", "test-string-value")

				entry.OptionalStringValue("test-key", "")
			},
			Matcher: map[string]objmatcher.Matcher{},
		},
		{
			Name: "StringListValue overrides StringValue",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringValue("test-key", "test-string-value")

				entry.StringListValue("test-key", []string{"test-string-list-value"})
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.SliceMatcher([]objmatcher.Matcher{
					objmatcher.NewEqualsMatcher("test-string-list-value"),
				}),
			},
		},
		{
			Name: "StringListValue that is empty overrides StringValue",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringValue("test-key", "test-string-value")

				entry.StringListValue("test-key", []string{})
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher([]any{}),
			},
		},
		{
			Name: "StringListValue that is nil overrides StringValue",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringValue("test-key", "test-string-value")

				entry.StringListValue("test-key", nil)
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher([]any{}),
			},
		},
		{
			Name: "StringMapValue overrides StringValue",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringValue("test-key", "test-string-value")

				entry.StringMapValue("test-key", map[string]string{
					"string-map-key-1": "string-map-value-1",
					"string-map-key-2": "string-map-value-2",
				})
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"string-map-key-1": objmatcher.NewEqualsMatcher("string-map-value-1"),
					"string-map-key-2": objmatcher.NewEqualsMatcher("string-map-value-2"),
				}),
			},
		},
		{
			Name: "StringMapValue that is empty overrides StringValue",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringValue("test-key", "test-string-value")

				entry.StringMapValue("test-key", map[string]string{})
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher(map[string]any{}),
			},
		},
		{
			Name: "StringMapValue that is nil overrides StringValue",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringValue("test-key", "test-string-value")

				entry.StringMapValue("test-key", nil)
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher(map[string]any{}),
			},
		},
		{
			Name: "StringMapValue overrides AnyMapValue",
			EntryFn: func(entry wlog.LogEntry) {
				entry.AnyMapValue("test-key", map[string]any{
					"shared-key-1":       1,
					"shared-key-2":       2,
					"any-map-key-unique": 3,
				})

				entry.StringMapValue("test-key", map[string]string{
					"shared-key-1":          "string-map-value-1",
					"shared-key-2":          "string-map-value-2",
					"string-map-key-unique": "string-map-value-3",
				})
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"shared-key-1":          objmatcher.NewEqualsMatcher("string-map-value-1"),
					"shared-key-2":          objmatcher.NewEqualsMatcher("string-map-value-2"),
					"string-map-key-unique": objmatcher.NewEqualsMatcher("string-map-value-3"),
				}),
			},
		},
		{
			Name: "AnyMapValue overrides StringValue",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringValue("test-key", "test-string-value")

				entry.AnyMapValue("test-key", map[string]any{
					"string-map-key-1": 1,
					"string-map-key-2": 2,
				})
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"string-map-key-1": objmatcher.NewEqualsMatcher(json.Number("1")),
					"string-map-key-2": objmatcher.NewEqualsMatcher(json.Number("2")),
				}),
			},
		},
		{
			Name: "AnyMapValue that is empty overrides StringValue",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringValue("test-key", "test-string-value")

				entry.AnyMapValue("test-key", map[string]any{})
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher(map[string]any{}),
			},
		},
		{
			Name: "AnyMapValue that is nil overrides StringValue",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringValue("test-key", "test-string-value")

				entry.AnyMapValue("test-key", nil)
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher(map[string]any{}),
			},
		},
		{
			Name: "AnyMapValue overrides StringMapValue",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringMapValue("test-key", map[string]string{
					"shared-key-1":          "string-map-value-1",
					"shared-key-2":          "string-map-value-2",
					"string-map-key-unique": "string-map-value-3",
				})

				entry.AnyMapValue("test-key", map[string]any{
					"shared-key-1":       1,
					"shared-key-2":       2,
					"any-map-key-unique": 3,
				})
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"shared-key-1":       objmatcher.NewEqualsMatcher(json.Number("1")),
					"shared-key-2":       objmatcher.NewEqualsMatcher(json.Number("2")),
					"any-map-key-unique": objmatcher.NewEqualsMatcher(json.Number("3")),
				}),
			},
		},
		{
			Name: "ObjectValue overrides StringValue",
			EntryFn: func(entry wlog.LogEntry) {
				entry.StringValue("test-key", "test-string-value")

				entry.ObjectValue("test-key", 13, nil)
			},
			Matcher: map[string]objmatcher.Matcher{
				"test-key": objmatcher.NewEqualsMatcher(json.Number("13")),
			},
		},
	}
}

func runEntryTestCases(t *testing.T, suiteName string, testCases []EntryTestCase, loggerProvider func(w io.Writer) wlog.Logger) {
	t.Run(suiteName, func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				var buf bytes.Buffer
				logger := loggerProvider(&buf)

				param := wlog.NewParam(tc.EntryFn)

				logger.Log(param)

				gotLog := map[string]interface{}{}
				logEntry := buf.Bytes()
				err := safejson.Unmarshal(logEntry, &gotLog)
				require.NoError(t, err, "Log line is not a valid map: %v", string(logEntry))

				assert.NoError(t, tc.Matcher.Matches(gotLog))
			})
		}
	})
}
