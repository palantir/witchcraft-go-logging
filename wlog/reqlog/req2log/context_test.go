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

package req2log_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/palantir/pkg/objmatcher"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/logreader"
	"github.com/palantir/witchcraft-go-logging/wlog/reqlog/req2log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestLogger(w io.Writer) req2log.Logger {
	return req2log.NewFromCreator(w, wlog.NewJSONMarshalLoggerProvider().NewLogger)
}

func TestFromContext(t *testing.T) {
	buf, ctx := newBufAndCtxWithLogger()

	logger := req2log.FromContext(ctx)
	exampleURL, _ := url.Parse("https://example.com/path")
	logger.Request(req2log.Request{
		Request: &http.Request{
			Method: http.MethodGet,
			URL:    exampleURL,
		},
		Duration: time.Second,
		RouteInfo: req2log.RouteInfo{
			Template:   "/{foo/{bar}",
			PathParams: map[string]string{"foo": "FOO", "bar": "BAR"},
		},
		ResponseStatus: 200,
	})

	entries, err := logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)

	assert.Equal(t, 1, len(entries))

	matcher := objmatcher.MapMatcher(map[string]objmatcher.Matcher{
		"time":         objmatcher.NewRegExpMatcher(".+"),
		"type":         objmatcher.NewEqualsMatcher("request.2"),
		"method":       objmatcher.NewEqualsMatcher("GET"),
		"path":         objmatcher.NewEqualsMatcher("/{foo/{bar}"),
		"duration":     objmatcher.NewEqualsMatcher(json.Number("1000000")),
		"protocol":     objmatcher.NewEqualsMatcher(""),
		"requestSize":  objmatcher.NewEqualsMatcher(json.Number("0")),
		"responseSize": objmatcher.NewEqualsMatcher(json.Number("0")),
		"status":       objmatcher.NewEqualsMatcher(json.Number("200")),
		"unsafeParams": objmatcher.NewEqualsMatcher(map[string]any{"foo": "FOO", "bar": "BAR"}),
	})
	err = matcher.Matches(map[string]interface{}(entries[0]))
	assert.NoError(t, err, "%v", err)
}

func newBufAndCtxWithLogger() (*bytes.Buffer, context.Context) {
	buf := &bytes.Buffer{}
	ctx := req2log.WithLogger(context.Background(), newTestLogger(buf))
	return buf, ctx
}
