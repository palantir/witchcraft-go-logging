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
	"fmt"
	"io"
	"strconv"
	"sync"
	"testing"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/logreader"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// lockedBuffer is a bytes.Buffer that is safe for concurrent writes.
type lockedBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *lockedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *lockedBuffer) Bytes() []byte {
	b.mu.Lock()
	defer b.mu.Unlock()
	return bytes.Clone(b.buf.Bytes())
}

// newLoggerWithSpareParamCapacity returns a logger that has been wrapped by repeated calls to svc1log.WithParams. A
// slice built up by repeated appends has more capacity than length, so a logger created this way is the case in which
// derived loggers and log calls must not write into the capacity that they share with their parent.
func newLoggerWithSpareParamCapacity(w io.Writer) svc1log.Logger {
	logger := svc1log.NewFromCreator(w, wlog.InfoLevel, wlog.NewJSONMarshalLoggerProvider().NewLeveledLogger)
	for i := range 5 {
		logger = svc1log.WithParams(logger, svc1log.SafeParam(fmt.Sprintf("base%d", i), i))
	}
	return logger
}

// Tests that loggers derived from the same parent logger keep their own params: logging with one derived logger must
// not emit the params that were added to a sibling.
func TestWithParamsDerivedLoggersDoNotShareParams(t *testing.T) {
	var buf bytes.Buffer
	parent := newLoggerWithSpareParamCapacity(&buf)

	first := svc1log.WithParams(parent, svc1log.SafeParam("derived", "first"))
	second := svc1log.WithParams(parent, svc1log.UnsafeParam("derived", "second"))

	first.Info("from first")
	second.Info("from second")
	parent.Info("from parent")

	entries, err := logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)
	require.Len(t, entries, 3)

	assert.Equal(t, "from first", entries[0]["message"])
	assert.Equal(t, "first", entries[0]["params"].(map[string]any)["derived"])
	assert.NotContains(t, entries[0], "unsafeParams", "an unsafe param of a sibling logger was written to this logger's entry")

	assert.Equal(t, "from second", entries[1]["message"])
	assert.NotContains(t, entries[1]["params"], "derived")
	assert.Equal(t, "second", entries[1]["unsafeParams"].(map[string]any)["derived"])

	assert.Equal(t, "from parent", entries[2]["message"])
	assert.NotContains(t, entries[2]["params"], "derived")
	assert.NotContains(t, entries[2], "unsafeParams")
}

// Tests that concurrent log calls on a single logger each emit only their own call-site params
// (https://github.com/palantir/witchcraft-go-logging/issues/56).
func TestConcurrentLogCallsDoNotShareParams(t *testing.T) {
	buf := &lockedBuffer{}
	logger := newLoggerWithSpareParamCapacity(buf)

	const numCalls = 200
	var wg sync.WaitGroup
	for i := range numCalls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			logger.Info(fmt.Sprintf("call %d", i), svc1log.SafeParam("call", i), svc1log.UnsafeParam("callUnsafe", i))
		}()
	}
	wg.Wait()

	entries, err := logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)
	require.Len(t, entries, numCalls)

	for _, entry := range entries {
		var call int
		_, err := fmt.Sscanf(entry["message"].(string), "call %d", &call)
		require.NoError(t, err)

		params := entry["params"].(map[string]any)
		unsafeParams := entry["unsafeParams"].(map[string]any)
		// the log reader decodes numbers as json.Number, so compare the printed values
		assert.Equal(t, strconv.Itoa(call), fmt.Sprint(params["call"]), "message %q carries params of another call", entry["message"])
		assert.Equal(t, strconv.Itoa(call), fmt.Sprint(unsafeParams["callUnsafe"]), "message %q carries unsafe params of another call", entry["message"])
	}
}
