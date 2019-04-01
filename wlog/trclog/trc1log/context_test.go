// Copyright (c) 2019 Palantir Technologies. All rights reserved.
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

package trc1log_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/palantir/pkg/objmatcher"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/logreader"
	"github.com/palantir/witchcraft-go-logging/wlog/trclog/trc1log"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests that the logger returned by trc1log.FromContext has UID, SID and TokenID parameters set on it if the context
// has those values set on it using wlog.
func TestFromContextUsesCommonIDs(t *testing.T) {
	buf := &bytes.Buffer{}
	ctx := trc1log.WithLogger(context.Background(), trc1log.NewFromCreator(buf, wlog.NewJSONMarshalLoggerProvider().NewLogger))

	ctx = wlog.ContextWithUID(ctx, "test-UID")
	ctx = wlog.ContextWithSID(ctx, "test-SID")
	ctx = wlog.ContextWithTokenID(ctx, "test-TokenID")

	logger := trc1log.FromContext(ctx)
	logger.Log(wtracing.SpanModel{})

	entries, err := logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)

	assert.Equal(t, 1, len(entries))

	matcher := objmatcher.MapMatcher(map[string]objmatcher.Matcher{
		"time":    objmatcher.NewRegExpMatcher(".+"),
		"span":    objmatcher.NewAnyMatcher(),
		"type":    objmatcher.NewEqualsMatcher("trace.1"),
		"uid":     objmatcher.NewEqualsMatcher("test-UID"),
		"sid":     objmatcher.NewEqualsMatcher("test-SID"),
		"tokenId": objmatcher.NewEqualsMatcher("test-TokenID"),
	})
	err = matcher.Matches(map[string]interface{}(entries[0]))
	assert.NoError(t, err, "%v", err)
}
