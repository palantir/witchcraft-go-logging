// Copyright (c) 2022 Palantir Technologies. All rights reserved.
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

package svc1zap

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/palantir/pkg/objmatcher"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	// Use zap as logger implementation
	_ "github.com/palantir/witchcraft-go-logging/wlog-zap"
)

func TestSvc1ZapWrapper(t *testing.T) {

	prefixParamFunc := func(key string, value interface{}) svc1log.Param {
		if strings.HasPrefix(key, "safe") {
			return svc1log.SafeParam(key, value)
		}
		if !strings.HasPrefix(key, "forbidden") {
			return svc1log.UnsafeParam(key, value)
		}
		return nil
	}

	t.Run("defaults to all unsafe params", func(t *testing.T) {
		buf := new(bytes.Buffer)
		logger := svc1log.New(buf, wlog.DebugLevel)
		logr1 := New(logger)
		logr1.Info("logr 1", zap.String("safeString", "string"), zap.String("forbiddenToken", "token"), zap.Int("unsafeInt", 42))
		assertLogLine(t, buf.Bytes(), objmatcher.MapMatcher{
			"level":   objmatcher.NewEqualsMatcher("INFO"),
			"time":    objmatcher.NewRegExpMatcher(".+"),
			"message": objmatcher.NewEqualsMatcher("logr 1"),
			"type":    objmatcher.NewEqualsMatcher(svc1log.TypeValue),
			"unsafeParams": objmatcher.MapMatcher{
				"forbiddenToken": objmatcher.NewEqualsMatcher("token"),
				"safeString":     objmatcher.NewEqualsMatcher("string"),
				"unsafeInt":      objmatcher.NewEqualsMatcher(float64(42)),
			},
		})
	})

	t.Run("caller origin and custom params", func(t *testing.T) {
		buf := new(bytes.Buffer)
		logger := svc1log.New(buf, wlog.DebugLevel)
		logr2 := New(logger, WithOriginFromZapCaller(), WithNewParamFunc(prefixParamFunc)).WithOptions(zap.AddCaller())
		logr2.Info("logr 2", zap.String("safeString", "string"), zap.String("forbiddenToken", "token"), zap.Int("unsafeInt", 42))
		assertLogLine(t, buf.Bytes(), objmatcher.MapMatcher{
			"level":   objmatcher.NewEqualsMatcher("INFO"),
			"time":    objmatcher.NewRegExpMatcher(".+"),
			"message": objmatcher.NewEqualsMatcher("logr 2"),
			"type":    objmatcher.NewEqualsMatcher(svc1log.TypeValue),
			"origin":  objmatcher.NewRegExpMatcher("^github.com/palantir/witchcraft-go-logging/adapters/svc1zap/svc1zap_test.go:\\d+"),
			"params": objmatcher.MapMatcher{
				"safeString": objmatcher.NewEqualsMatcher("string"),
			},
			"unsafeParams": objmatcher.MapMatcher{
				"unsafeInt": objmatcher.NewEqualsMatcher(float64(42)),
			},
		})
	})

	t.Run("logger with attached params", func(t *testing.T) {
		buf := new(bytes.Buffer)
		logger := svc1log.New(buf, wlog.DebugLevel)
		logr3 := New(logger).Named("logr3").With(zap.String("name", "logr3"))
		logr3.Error("logr 3", zap.String("safeString", "string"), zap.String("forbiddenToken", "token"), zap.Int("unsafeInt", 42))
		assertLogLine(t, buf.Bytes(), objmatcher.MapMatcher{
			"level":   objmatcher.NewEqualsMatcher("ERROR"),
			"time":    objmatcher.NewRegExpMatcher(".+"),
			"message": objmatcher.NewEqualsMatcher("logr3: logr 3"),
			"type":    objmatcher.NewEqualsMatcher(svc1log.TypeValue),
			"unsafeParams": objmatcher.MapMatcher{
				"forbiddenToken": objmatcher.NewEqualsMatcher("token"),
				"name":           objmatcher.NewEqualsMatcher("logr3"),
				"safeString":     objmatcher.NewEqualsMatcher("string"),
				"unsafeInt":      objmatcher.NewEqualsMatcher(float64(42)),
			},
		})
	})

	t.Run("logger with disabled level", func(t *testing.T) {
		buf := new(bytes.Buffer)
		logger := svc1log.New(buf, wlog.InfoLevel)
		logr4 := New(logger)
		logr4.Debug("logr 4")
		assert.Empty(t, buf.String())
	})

	t.Run("logger with mutator", func(t *testing.T) {
		buf := new(bytes.Buffer)
		logger := svc1log.New(buf, wlog.DebugLevel)
		logr5 := New(logger, WithEntryMutatorFunc(func(entry zapcore.Entry) (out zapcore.Entry, ok bool) {
			if entry.Message == "skip me" {
				return entry, false
			}
			if entry.Message == "drop to debug" {
				entry.Level = zapcore.DebugLevel
			}
			return entry, true
		}))
		logr5.Info("logr 5")
		logr5.Warn("skip me")
		logr5.Warn("drop to debug")

		lines := bytes.Split(bytes.TrimSpace(buf.Bytes()), []byte("\n"))
		if assert.Len(t, lines, 2) {
			assertLogLine(t, lines[0], objmatcher.MapMatcher{
				"type":    objmatcher.NewEqualsMatcher(svc1log.TypeValue),
				"time":    objmatcher.NewRegExpMatcher(".+"),
				"level":   objmatcher.NewEqualsMatcher("INFO"),
				"message": objmatcher.NewEqualsMatcher("logr 5"),
			})
			assertLogLine(t, lines[1], objmatcher.MapMatcher{
				"type":    objmatcher.NewEqualsMatcher(svc1log.TypeValue),
				"time":    objmatcher.NewRegExpMatcher(".+"),
				"level":   objmatcher.NewEqualsMatcher("DEBUG"),
				"message": objmatcher.NewEqualsMatcher("drop to debug"),
			})
		}
	})
}

func assertLogLine(t *testing.T, logLine []byte, matcher objmatcher.MapMatcher) {
	logEntry := map[string]interface{}{}
	err := json.Unmarshal(logLine, &logEntry)
	assert.NoError(t, err)
	assert.NoError(t, matcher.Matches(logEntry))
}
