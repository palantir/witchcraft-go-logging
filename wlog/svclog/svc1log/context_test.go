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
	"io"
	"testing"

	"github.com/palantir/pkg/objmatcher"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/logreader"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	wparams "github.com/palantir/witchcraft-go-params"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
	"github.com/palantir/witchcraft-go-tracing/wzipkin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestLogger(w io.Writer, origin string) svc1log.Logger {
	return svc1log.WithParams(svc1log.NewFromCreator(w, wlog.InfoLevel, wlog.NewJSONMarshalLoggerProvider().NewLeveledLogger), svc1log.Origin(origin))
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

// Tests that the logger returned by svc1log.FromContext has UID, SID and TokenID parameters set on it if the context
// has those values set on it using wlog.
func TestFromContextUsesCommonIDs(t *testing.T) {
	buf, ctx := newBufAndCtxWithLogger()

	ctx = wlog.ContextWithUID(ctx, "test-UID")
	ctx = wlog.ContextWithSID(ctx, "test-SID")
	ctx = wlog.ContextWithTokenID(ctx, "test-TokenID")

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
		"uid":     objmatcher.NewEqualsMatcher("test-UID"),
		"sid":     objmatcher.NewEqualsMatcher("test-SID"),
		"tokenId": objmatcher.NewEqualsMatcher("test-TokenID"),
	})
	err = matcher.Matches(map[string]interface{}(entries[0]))
	assert.NoError(t, err, "%v", err)
}

// Tests that the logger returned by svc1log.FromContext has a TraceID set on it if the context has a wtracing TraceID.
func TestFromContextSetsTraceID(t *testing.T) {
	buf, ctx := newBufAndCtxWithLogger()

	// create a no-op tracer to use for the test
	tracer, err := wzipkin.NewTracer(wtracing.NewNoopReporter())
	require.NoError(t, err)

	createMatcher := func(msg, traceID string) objmatcher.Matcher {
		matcher := objmatcher.MapMatcher(map[string]objmatcher.Matcher{
			"level":   objmatcher.NewEqualsMatcher("INFO"),
			"time":    objmatcher.NewRegExpMatcher(".+"),
			"origin":  objmatcher.NewEqualsMatcher("com.palantir.test"),
			"type":    objmatcher.NewEqualsMatcher("service.1"),
			"message": objmatcher.NewEqualsMatcher(msg),
		})
		if traceID != "" {
			matcher["traceId"] = objmatcher.NewEqualsMatcher(traceID)
		}
		return matcher
	}

	// logger output should have no TraceID (none set as parameter and none exists in context)
	logger := svc1log.FromContext(ctx)
	logger.Info("Message0")

	entries, err := logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)
	assert.Equal(t, 1, len(entries))
	matcher := createMatcher("Message0", "")
	err = matcher.Matches(map[string]interface{}(entries[0]))
	assert.NoError(t, err, "%v", err)
	buf.Reset()

	// logger output should have TraceID set in context (span is set on context)
	spanOne := tracer.StartSpan("spanOne")
	ctx = wtracing.ContextWithSpan(ctx, spanOne)
	logger = svc1log.FromContext(ctx)
	logger.Info("Message1")

	entries, err = logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)
	assert.Equal(t, 1, len(entries))
	matcher = createMatcher("Message1", string(spanOne.Context().TraceID))
	err = matcher.Matches(map[string]interface{}(entries[0]))
	assert.NoError(t, err, "%v", err)
	buf.Reset()

	// manually adding a TraceID parameter will override the TraceID (because it is applied after the context one)
	logger = svc1log.WithParams(logger, svc1log.TraceID("manually-set-trace-id"))
	logger.Info("Message2")

	entries, err = logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)
	assert.Equal(t, 1, len(entries))
	matcher = createMatcher("Message2", "manually-set-trace-id")
	err = matcher.Matches(map[string]interface{}(entries[0]))
	assert.NoError(t, err, "%v", err)
	buf.Reset()
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
