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

package audit2log_test

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
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit2log"
	"github.com/palantir/witchcraft-go-logging/wlog/logreader"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
	"github.com/palantir/witchcraft-go-tracing/wzipkin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestLogger(w io.Writer) audit2log.Logger {
	return &testAudit2Logger{
		w: w,
	}
}

type testAudit2Logger struct {
	w io.Writer
}

func (l *testAudit2Logger) Audit(name string, result audit2log.AuditResultType, params ...audit2log.Param) {
	entry := wlog.NewMapLogEntry()
	entry.StringValue(wlog.TypeKey, audit2log.TypeValue)
	entry.StringValue(wlog.TimeKey, time.Now().Format(time.RFC3339Nano))
	entry.StringValue(audit2log.NameKey, name)
	entry.StringValue(audit2log.ResultKey, string(result))
	for _, p := range params {
		audit2log.ApplyParam(p, entry)
	}
	jsonBytes, _ := json.Marshal(entry.AllValues())
	fmt.Fprintln(l.w, string(jsonBytes))
}

func TestFromContext(t *testing.T) {
	buf, ctx := newBufAndCtxWithLogger()

	logger := audit2log.FromContext(ctx)
	logger.Audit("TEST_ENTRY", audit2log.AuditResultSuccess)

	entries, err := logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)

	assert.Equal(t, 1, len(entries))

	matcher := objmatcher.MapMatcher(map[string]objmatcher.Matcher{
		"time":   objmatcher.NewRegExpMatcher(".+"),
		"type":   objmatcher.NewEqualsMatcher("audit.2"),
		"name":   objmatcher.NewEqualsMatcher("TEST_ENTRY"),
		"result": objmatcher.NewEqualsMatcher("SUCCESS"),
	})
	err = matcher.Matches(map[string]interface{}(entries[0]))
	assert.NoError(t, err, "%v", err)
}

// Tests that the logger returned by audit2log.FromContext has a TraceID set on it if the context has a wtracing
// TraceID.
func TestFromContextSetsTraceID(t *testing.T) {
	buf, ctx := newBufAndCtxWithLogger()

	// create a no-op tracer to use for the test
	tracer, err := wzipkin.NewTracer(wtracing.NewNoopReporter())
	require.NoError(t, err)

	createMatcher := func(name, traceID string) objmatcher.Matcher {
		matcher := objmatcher.MapMatcher(map[string]objmatcher.Matcher{
			"time":   objmatcher.NewRegExpMatcher(".+"),
			"type":   objmatcher.NewEqualsMatcher("audit.2"),
			"name":   objmatcher.NewEqualsMatcher(name),
			"result": objmatcher.NewEqualsMatcher("SUCCESS"),
		})
		if traceID != "" {
			matcher["traceId"] = objmatcher.NewEqualsMatcher(traceID)
		}
		return matcher
	}

	// logger output should have no TraceID (none set as parameter and none exists in context)
	logger := audit2log.FromContext(ctx)
	logger.Audit("EVENT_0", audit2log.AuditResultSuccess)

	entries, err := logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)
	assert.Equal(t, 1, len(entries))
	matcher := createMatcher("EVENT_0", "")
	err = matcher.Matches(map[string]interface{}(entries[0]))
	assert.NoError(t, err, "%v", err)
	buf.Reset()

	// logger output should have TraceID set in context (span is set on context)
	spanOne := tracer.StartSpan("spanOne")
	ctx = wtracing.ContextWithSpan(ctx, spanOne)
	logger = audit2log.FromContext(ctx)
	logger.Audit("EVENT_1", audit2log.AuditResultSuccess)

	entries, err = logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)
	assert.Equal(t, 1, len(entries))
	matcher = createMatcher("EVENT_1", string(spanOne.Context().TraceID))
	err = matcher.Matches(map[string]interface{}(entries[0]))
	assert.NoError(t, err, "%v", err)
	buf.Reset()

	// manually adding a TraceID parameter will override the TraceID (because it is applied after the context one)
	logger = audit2log.WithParams(logger, audit2log.TraceID("manually-set-trace-id"))
	logger.Audit("EVENT_2", audit2log.AuditResultSuccess)

	entries, err = logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)
	assert.Equal(t, 1, len(entries))
	matcher = createMatcher("EVENT_2", "manually-set-trace-id")
	err = matcher.Matches(map[string]interface{}(entries[0]))
	assert.NoError(t, err, "%v", err)
	buf.Reset()
}

func newBufAndCtxWithLogger() (*bytes.Buffer, context.Context) {
	buf := &bytes.Buffer{}
	ctx := audit2log.WithLogger(context.Background(), newTestLogger(buf))
	return buf, ctx
}
