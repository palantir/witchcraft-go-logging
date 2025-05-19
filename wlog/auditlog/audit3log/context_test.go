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

package audit3log

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/palantir/pkg/objmatcher"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/logreader"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
	"github.com/palantir/witchcraft-go-tracing/wzipkin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	typeValue       = "audit.3"
	nameValue       = "TEST_ENTRY"
	resultValue     = "SUCCESS"
	uuidRegexpValue = "^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[8|9|a|b][a-f0-9]{3}-[a-f0-9]{12}$"
	uidValue        = "test-UID"
	sidValue        = "test-SID"
	tokenIDValue    = "test-TokenID"
	orgIDValue      = "test-OrgID"
)

func newTestLogger(w io.Writer) Logger {
	return NewFromCreator(w, wlog.NewJSONMarshalLoggerProvider().NewLogger)
}

func TestFromContext(t *testing.T) {
	buf, ctx := newBufAndCtxWithLogger()

	logger := FromContext(ctx)
	logger.Audit("TEST_ENTRY", AuditResultSuccess)

	entries, err := logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)

	assert.Equal(t, 1, len(entries))

	matcher := objmatcher.MapMatcher(map[string]objmatcher.Matcher{
		"time":       objmatcher.NewRegExpMatcher(".+"),
		"type":       objmatcher.NewEqualsMatcher(typeValue),
		"name":       objmatcher.NewEqualsMatcher(nameValue),
		"result":     objmatcher.NewEqualsMatcher(resultValue),
		"logEntryId": objmatcher.NewRegExpMatcher(uuidRegexpValue),
	})
	err = matcher.Matches(map[string]interface{}(entries[0]))
	assert.NoError(t, err, "%v", err)
}

// Tests that the logger returned by FromContext has UID, SID, TokenID, and OrgID parameters set on it if the context
// has those values set on it using wlog.
func TestFromContextUsesCommonIDs(t *testing.T) {
	buf, ctx := newBufAndCtxWithLogger()

	ctx = wlog.ContextWithUID(ctx, uidValue)
	ctx = wlog.ContextWithSID(ctx, sidValue)
	ctx = wlog.ContextWithTokenID(ctx, tokenIDValue)
	ctx = wlog.ContextWithOrgID(ctx, orgIDValue)

	logger := FromContext(ctx)
	logger.Audit(nameValue, AuditResultSuccess)

	entries, err := logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)

	assert.Equal(t, 1, len(entries))

	matcher := objmatcher.MapMatcher(map[string]objmatcher.Matcher{
		"time":       objmatcher.NewRegExpMatcher(".+"),
		"type":       objmatcher.NewEqualsMatcher(typeValue),
		"name":       objmatcher.NewEqualsMatcher(nameValue),
		"result":     objmatcher.NewEqualsMatcher(resultValue),
		"uid":        objmatcher.NewEqualsMatcher(uidValue),
		"sid":        objmatcher.NewEqualsMatcher(sidValue),
		"tokenId":    objmatcher.NewEqualsMatcher(tokenIDValue),
		"orgId":      objmatcher.NewEqualsMatcher(orgIDValue),
		"logEntryId": objmatcher.NewRegExpMatcher(uuidRegexpValue),
	})
	err = matcher.Matches(map[string]interface{}(entries[0]))
	assert.NoError(t, err, "%v", err)
}

// Tests that the logger returned by FromContext has a TraceID set on it if the context has a wtracing
// TraceID.
func TestFromContextSetsTraceID(t *testing.T) {
	buf, ctx := newBufAndCtxWithLogger()

	// create a no-op tracer to use for the test
	tracer, err := wzipkin.NewTracer(wtracing.NewNoopReporter())
	require.NoError(t, err)

	createMatcher := func(name, traceID string) objmatcher.Matcher {
		matcher := objmatcher.MapMatcher(map[string]objmatcher.Matcher{
			"time":       objmatcher.NewRegExpMatcher(".+"),
			"type":       objmatcher.NewEqualsMatcher(typeValue),
			"name":       objmatcher.NewEqualsMatcher(name),
			"result":     objmatcher.NewEqualsMatcher(resultValue),
			"logEntryId": objmatcher.NewRegExpMatcher(uuidRegexpValue),
		})
		if traceID != "" {
			matcher["traceId"] = objmatcher.NewEqualsMatcher(traceID)
		}
		return matcher
	}

	// logger output should have no TraceID (none set as parameter and none exists in context)
	logger := FromContext(ctx)
	logger.Audit("EVENT_0", AuditResultSuccess)

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
	logger = FromContext(ctx)
	logger.Audit("EVENT_1", AuditResultSuccess)

	entries, err = logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)
	assert.Equal(t, 1, len(entries))
	matcher = createMatcher("EVENT_1", string(spanOne.Context().TraceID))
	err = matcher.Matches(map[string]interface{}(entries[0]))
	assert.NoError(t, err, "%v", err)
	buf.Reset()

	// manually adding a TraceID parameter will override the TraceID (because it is applied after the context one)
	logger = WithParams(logger, TraceID("manually-set-trace-id"))
	logger.Audit("EVENT_2", AuditResultSuccess)

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
	ctx := WithLogger(context.Background(), newTestLogger(buf))
	return buf, ctx
}
