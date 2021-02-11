// Copyright (c) 2021 Palantir Technologies. All rights reserved.
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
	"io"
	"testing"

	"github.com/palantir/pkg/refreshable"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/logreader"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/stretchr/testify/require"
)

func newOriginLevelTestLogger(w io.Writer, config refreshable.Refreshable, origin string, provider wlog.LoggerProvider) svc1log.Logger {
	return svc1log.WithParams(svc1log.NewOriginLevelFromCreator(w, config, provider.NewLeveledLogger), svc1log.Origin(origin))
}

func newBufAndCtxWithOriginLevelLogger(provider wlog.LoggerProvider, config refreshable.Refreshable) (*bytes.Buffer, context.Context) {
	buf := &bytes.Buffer{}
	ctx := svc1log.WithLogger(context.Background(), newOriginLevelTestLogger(buf, config, "com.palantir.test", provider))
	return buf, ctx
}

func TestOriginLevels(t *testing.T) {
	config := refreshable.NewDefaultRefreshable(svc1log.OriginLevelLoggerConfig{
		Level: wlog.Warn,
		PerOriginLevels: map[string]wlog.Level{
			"foo": wlog.Info,
			"bar": wlog.Debug,
		},
	})
	buf, ctx := newBufAndCtxWithOriginLevelLogger(wlog.NewJSONMarshalLoggerProvider(), config)

	logger := svc1log.FromContext(ctx)
	logger.Info("Test")

	// No entries by default at info level
	entries, err := logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)
	require.Len(t, entries, 0)
	buf.Reset()

	// "foo" origin allows info level logging, but not debug level
	logger.Info("test", svc1log.Origin("foo"))
	logger.Debug("test", svc1log.Origin("foo"))
	entries, err = logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)
	require.Len(t, entries, 1)
	buf.Reset()

	// "bar" origin allows debug level logging
	logger.Info("test", svc1log.Origin("bar"))
	logger.Debug("test", svc1log.Origin("bar"))
	entries, err = logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)
	require.Len(t, entries, 2)
	buf.Reset()

	// refreshing the configuration works
	require.NoError(t, config.Update(svc1log.OriginLevelLoggerConfig{
		Level: wlog.Info,
		PerOriginLevels: map[string]wlog.Level{
			"foo": wlog.Debug,
			"bar": wlog.Debug,
		},
	}))
	logger.Info("Test")
	logger.Debug("test", svc1log.Origin("foo"))
	logger.Debug("test", svc1log.Origin("bar"))
	entries, err = logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)
	require.Len(t, entries, 3)
}
