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

package wrapped1log_test

import (
	"bytes"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/logreader"
	"github.com/palantir/witchcraft-go-logging/wlog/reqlog/req2log"
	"github.com/palantir/witchcraft-go-logging/wlog/wrappedlog/wrapped1log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestCreatorParamKeepsWrapping(t *testing.T) {
	provider := wlog.NewJSONMarshalLoggerProvider()
	var defaultBuf, creatorBuf bytes.Buffer
	logger := wrapped1log.NewFromProvider(&defaultBuf, wlog.InfoLevel, provider, "entity", "1.0").Request(
		req2log.Creator(func(io.Writer) wlog.Logger { return provider.NewLogger(&creatorBuf) }),
	)
	logger.Request(req2log.Request{Request: httptest.NewRequest("GET", "/path", nil)})

	assert.Empty(t, defaultBuf.String())
	entries, err := logreader.EntriesFromContent(creatorBuf.Bytes())
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, "wrapped.1", entries[0]["type"])
	assert.Equal(t, "entity", entries[0]["entityName"])
	assert.Equal(t, "requestLogV2", entries[0]["payload"].(map[string]any)["type"])
}
