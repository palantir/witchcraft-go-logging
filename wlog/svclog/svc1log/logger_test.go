// Copyright (c) 2026 Palantir Technologies. All rights reserved.
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

package svc1log_test

import (
	"bytes"
	"encoding/json"
	"io"
	"maps"
	"slices"
	"testing"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/stretchr/testify/require"
)

func TestEntryProviderDispatch(t *testing.T) {
	backend := &recordingLeveledLogger{}
	logger := svc1log.NewFromCreator(io.Discard, wlog.InfoLevel, func(io.Writer, wlog.LogLevel) wlog.LeveledLogger {
		return backend
	}, svc1log.SafeParam("stored", true))
	for _, tc := range []struct {
		level string
		emit  func(string, ...svc1log.Param)
	}{
		{"DEBUG", logger.Debug},
		{"INFO", logger.Info},
		{"WARN", logger.Warn},
		{"ERROR", logger.Error},
	} {
		t.Run(tc.level, func(t *testing.T) {
			for _, message := range []string{"message", ""} {
				tc.emit(message, svc1log.SafeParam("call", 42))
				require.Equal(t, tc.level, backend.method)
				require.Equal(t, message, backend.message)
				require.Equal(t, tc.level, backend.fields["level"])
				require.Equal(t, map[string]any{"stored": true, "call": 42}, backend.fields["params"])
			}
		})
	}
	logger.SetLevel(wlog.ErrorLevel)
	require.Equal(t, wlog.ErrorLevel, backend.level)
}

func TestPooledLogsOmitUnsetFields(t *testing.T) {
	var buf bytes.Buffer
	logger := svc1log.NewFromCreator(&buf, wlog.InfoLevel, wlog.NewJSONMarshalLoggerProvider().NewLeveledLogger)
	for range 3 {
		logger.Info("set", svc1log.SafeParam("a", 1), svc1log.UnsafeParam("b", 2), svc1log.Tag("c", "3"), svc1log.Origin("origin"), svc1log.UID("uid"))
		buf.Reset()
		logger.Info("unset", svc1log.Origin(""), svc1log.UID(""))
		var got map[string]any
		require.NoError(t, json.Unmarshal(buf.Bytes(), &got))
		require.Equal(t, []string{"level", "message", "time", "type"}, slices.Sorted(maps.Keys(got)))
		buf.Reset()
	}
}

// This provider deliberately has no LevelChecker; it owns filtering its entry calls.
type recordingLeveledLogger struct {
	method  string
	message string
	fields  map[string]any
	level   wlog.LogLevel
}

func (l *recordingLeveledLogger) record(method, message string, params []wlog.Param) {
	l.method, l.message = method, message
	entry := wlog.NewMapLogEntry()
	wlog.ApplyParams(entry, params)
	l.fields = entry.AllValues()
}

func (l *recordingLeveledLogger) Debug(msg string, params ...wlog.Param) {
	l.record("DEBUG", msg, params)
}

func (l *recordingLeveledLogger) Info(msg string, params ...wlog.Param) {
	l.record("INFO", msg, params)
}

func (l *recordingLeveledLogger) Warn(msg string, params ...wlog.Param) {
	l.record("WARN", msg, params)
}

func (l *recordingLeveledLogger) Error(msg string, params ...wlog.Param) {
	l.record("ERROR", msg, params)
}

func (l *recordingLeveledLogger) SetLevel(level wlog.LogLevel) {
	l.level = level
}
