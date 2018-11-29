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

package svc1log_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/palantir/pkg/objmatcher"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/logreader"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/palantir/witchcraft-go-params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestLogger(w io.Writer, origin string) svc1log.Logger {
	return &testSvc1Logger{
		w:      w,
		origin: origin,
	}
}

type testSvc1Logger struct {
	origin string
	w      io.Writer
}

func (l *testSvc1Logger) Debug(msg string, params ...svc1log.Param) {
	l.log(msg, "DEBUG", params...)
}

func (l *testSvc1Logger) Info(msg string, params ...svc1log.Param) {
	l.log(msg, "INFO", params...)
}

func (l *testSvc1Logger) Warn(msg string, params ...svc1log.Param) {
	l.log(msg, "WARN", params...)
}

func (l *testSvc1Logger) Error(msg string, params ...svc1log.Param) {
	l.log(msg, "ERROR", params...)
}

func (l *testSvc1Logger) SetLevel(level wlog.LogLevel) {}

func (l *testSvc1Logger) log(msg, lvl string, params ...svc1log.Param) {
	entry := wlog.NewMapLogEntry()
	entry.StringValue(wlog.TypeKey, svc1log.TypeValue)
	entry.StringValue(svc1log.MessageKey, msg)
	entry.StringValue(svc1log.LevelKey, lvl)
	entry.StringValue(svc1log.OriginKey, l.origin)
	entry.StringValue(wlog.TimeKey, time.Now().Format(time.RFC3339Nano))
	for _, p := range params {
		svc1log.ApplyParam(p, entry)
	}
	jsonBytes, _ := json.Marshal(entry.AllValues())
	fmt.Fprintln(l.w, string(jsonBytes))
}

func TestFromContext(t *testing.T) {
	buf, ctx := newBufAndCtxWithLogger()

	logger := svc1log.FromContext(ctx)
	logger.Info("Test")

	entries, err := logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)

	assert.Equal(t, 1, len(entries))

	matcher := objmatcher.MapMatcher(map[string]objmatcher.Matcher{
		"level":   objmatcher.NewEqualsMatcher("INFO"),
		"time":    objmatcher.NewRegExpMatcher(".+"),
		"origin":  objmatcher.NewEqualsMatcher("com.palantir.test"),
		"type":    objmatcher.NewEqualsMatcher("service.1"),
		"message": objmatcher.NewEqualsMatcher("Test"),
	})
	err = matcher.Matches(map[string]interface{}(entries[0]))
	assert.NoError(t, err, "%v", err)
}

func TestWithLoggerParams(t *testing.T) {
	buf, ctx := newBufAndCtxWithLogger()

	ctx = svc1log.WithLoggerParams(ctx, svc1log.SafeParam("foo", "bar"))
	ctx = svc1log.WithLoggerParams(ctx, svc1log.SafeParam("ten", 10))

	logger := svc1log.FromContext(ctx)
	logger.Info("Test")

	entries, err := logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)

	assert.Equal(t, 1, len(entries))

	matcher := objmatcher.MapMatcher(map[string]objmatcher.Matcher{
		"level":   objmatcher.NewEqualsMatcher("INFO"),
		"time":    objmatcher.NewRegExpMatcher(".+"),
		"origin":  objmatcher.NewEqualsMatcher("com.palantir.test"),
		"type":    objmatcher.NewEqualsMatcher("service.1"),
		"message": objmatcher.NewEqualsMatcher("Test"),
		"params": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
			"foo": objmatcher.NewEqualsMatcher("bar"),
			"ten": objmatcher.NewEqualsMatcher(json.Number("10")),
		}),
		"unsafeParams": objmatcher.MapMatcher(map[string]objmatcher.Matcher{}),
	})
	err = matcher.Matches(map[string]interface{}(entries[0]))
	assert.NoError(t, err, "%v", err)
}

func TestWParamsSafeAndUnsafeParamsUsed(t *testing.T) {
	buf, ctx := newBufAndCtxWithLogger()

	ctx = wparams.ContextWithSafeParam(ctx, "foo", "bar")
	ctx = wparams.ContextWithSafeParam(ctx, "ten", 10)
	ctx = wparams.ContextWithUnsafeParam(ctx, "unsafe", "secret")

	logger := svc1log.FromContext(ctx)
	logger.Info("Test")

	entries, err := logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)

	assert.Equal(t, 1, len(entries))

	matcher := objmatcher.MapMatcher(map[string]objmatcher.Matcher{
		"level":   objmatcher.NewEqualsMatcher("INFO"),
		"time":    objmatcher.NewRegExpMatcher(".+"),
		"origin":  objmatcher.NewEqualsMatcher("com.palantir.test"),
		"type":    objmatcher.NewEqualsMatcher("service.1"),
		"message": objmatcher.NewEqualsMatcher("Test"),
		"params": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
			"foo": objmatcher.NewEqualsMatcher("bar"),
			"ten": objmatcher.NewEqualsMatcher(json.Number("10")),
		}),
		"unsafeParams": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
			"unsafe": objmatcher.NewEqualsMatcher("secret"),
		}),
	})
	err = matcher.Matches(map[string]interface{}(entries[0]))
	assert.NoError(t, err, "%v", err)
}

func TestWithLoggerParamsSetsWParamsSafeAndUnsafeParams(t *testing.T) {
	_, ctx := newBufAndCtxWithLogger()

	ctx = svc1log.WithLoggerParams(ctx, svc1log.SafeParam("foo", "bar"))
	ctx = svc1log.WithLoggerParams(ctx, svc1log.UnsafeParam("ten", 10))

	safe, unsafe := wparams.SafeAndUnsafeParamsFromContext(ctx)
	assert.Equal(t, map[string]interface{}{
		"foo": "bar",
	}, safe)
	assert.Equal(t, map[string]interface{}{
		"ten": 10,
	}, unsafe)
}

func newBufAndCtxWithLogger() (*bytes.Buffer, context.Context) {
	buf := &bytes.Buffer{}
	ctx := svc1log.WithLogger(context.Background(), newTestLogger(buf, "com.palantir.test"))
	return buf, ctx
}
