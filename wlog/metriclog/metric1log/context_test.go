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

package metric1log_test

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/palantir/pkg/objmatcher"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/logreader"
	"github.com/palantir/witchcraft-go-logging/wlog/metriclog/metric1log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestLogger(w io.Writer) metric1log.Logger {
	return metric1log.NewFromCreator(w, wlog.NewJSONMarshalLoggerProvider().NewLogger)
}

func TestFromContext(t *testing.T) {
	buf, ctx := newBufAndCtxWithLogger()

	logger := metric1log.FromContext(ctx)
	logger.Metric("com.palantir.metric", "gauge")

	entries, err := logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)

	assert.Equal(t, 1, len(entries))

	matcher := objmatcher.MapMatcher(map[string]objmatcher.Matcher{
		"time":       objmatcher.NewRegExpMatcher(".+"),
		"type":       objmatcher.NewEqualsMatcher("metric.1"),
		"metricName": objmatcher.NewEqualsMatcher("com.palantir.metric"),
		"metricType": objmatcher.NewEqualsMatcher("gauge"),
	})
	err = matcher.Matches(map[string]any(entries[0]))
	assert.NoError(t, err, "%v", err)
}

func TestExplicitTime(t *testing.T) {
	buf := &bytes.Buffer{}
	timestamp := time.Date(2026, time.September, 2, 12, 34, 56, 789123456, time.UTC)

	newTestLogger(buf).Metric("com.palantir.metric", "gauge", metric1log.Time(timestamp))

	entries, err := logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, timestamp.Format(time.RFC3339Nano), entries[0]["time"])
}

// Tests that the logger returned by metric1log.FromContext has UID, SID, TokenID, and OrgID parameters set on it if the context
// has those values set on it using wlog.
func TestFromContextUsesCommonIDs(t *testing.T) {
	buf, ctx := newBufAndCtxWithLogger()

	ctx = wlog.ContextWithUID(ctx, "test-UID")
	ctx = wlog.ContextWithSID(ctx, "test-SID")
	ctx = wlog.ContextWithTokenID(ctx, "test-TokenID")
	ctx = wlog.ContextWithOrgID(ctx, "test-OrgID")

	logger := metric1log.FromContext(ctx)
	logger.Metric("com.palantir.metric", "gauge")

	entries, err := logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)

	assert.Equal(t, 1, len(entries))

	matcher := objmatcher.MapMatcher(map[string]objmatcher.Matcher{
		"time":       objmatcher.NewRegExpMatcher(".+"),
		"type":       objmatcher.NewEqualsMatcher("metric.1"),
		"metricName": objmatcher.NewEqualsMatcher("com.palantir.metric"),
		"metricType": objmatcher.NewEqualsMatcher("gauge"),
		"uid":        objmatcher.NewEqualsMatcher("test-UID"),
		"sid":        objmatcher.NewEqualsMatcher("test-SID"),
		"tokenId":    objmatcher.NewEqualsMatcher("test-TokenID"),
		"orgId":      objmatcher.NewEqualsMatcher("test-OrgID"),
	})
	err = matcher.Matches(map[string]any(entries[0]))
	assert.NoError(t, err, "%v", err)
}

func newBufAndCtxWithLogger() (*bytes.Buffer, context.Context) {
	buf := &bytes.Buffer{}
	ctx := metric1log.WithLogger(context.Background(), newTestLogger(buf))
	return buf, ctx
}
