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

package wlog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"testing"

	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit2log"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit3log"
	"github.com/palantir/witchcraft-go-logging/wlog/diaglog/diag1log"
	"github.com/palantir/witchcraft-go-logging/wlog/evtlog/evt2log"
	"github.com/palantir/witchcraft-go-logging/wlog/metriclog/metric1log"
	"github.com/palantir/witchcraft-go-logging/wlog/trclog/trc1log"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
	"github.com/stretchr/testify/require"
)

func TestTypedWithLoggerParams(t *testing.T) {
	for name, provider := range jsonProviders() {
		t.Run(name, func(t *testing.T) {
			t.Run("audit2", func(t *testing.T) {
				checkContextParams(t, func(w io.Writer) audit2log.Logger { return audit2log.NewFromCreator(w, provider.NewLogger) },
					audit2log.WithLogger, audit2log.WithLoggerParams, audit2log.FromContext,
					func(value string) audit2log.Param { return audit2log.RequestParam("key", value) },
					func(logger audit2log.Logger) { logger.Audit("test", audit2log.AuditResultSuccess) }, "requestParams")
			})
			t.Run("audit3", func(t *testing.T) {
				checkContextParams(t, func(w io.Writer) audit3log.Logger {
					return audit3log.WithParams(audit3log.NewFromCreator(w, provider.NewLogger), audit3log.EventID("c15487b9-ff6a-4bb1-8c25-2433a185c438"))
				}, audit3log.WithLogger, audit3log.WithLoggerParams, audit3log.FromContext,
					func(value string) audit3log.Param { return audit3log.RequestField("key", value) },
					func(logger audit3log.Logger) { logger.Audit("test", audit3log.AuditResultSuccess) }, "requestFields")
			})
			t.Run("event", func(t *testing.T) {
				checkContextParams(t, func(w io.Writer) evt2log.Logger { return evt2log.NewFromCreator(w, provider.NewLogger) },
					evt2log.WithLogger, evt2log.WithLoggerParams, evt2log.FromContext,
					func(value string) evt2log.Param { return evt2log.Value("key", value) },
					func(logger evt2log.Logger) { logger.Event("test") }, "values")
			})
			t.Run("metric", func(t *testing.T) {
				checkContextParams(t, func(w io.Writer) metric1log.Logger { return metric1log.NewFromCreator(w, provider.NewLogger) },
					metric1log.WithLogger, metric1log.WithLoggerParams, metric1log.FromContext,
					func(value string) metric1log.Param { return metric1log.Value("key", value) },
					func(logger metric1log.Logger) { logger.Metric("test", "gauge") }, "values")
			})
			t.Run("trace", func(t *testing.T) {
				checkContextParams(t, func(w io.Writer) trc1log.Logger { return trc1log.NewFromCreator(w, provider.NewLogger) },
					trc1log.WithLogger, trc1log.WithLoggerParams, trc1log.FromContext,
					func(value string) trc1log.Param { return trc1log.UnsafeParam("key", value) },
					func(logger trc1log.Logger) { logger.Log(wtracing.SpanModel{}) }, "unsafeParams")
			})
			t.Run("diagnostic", func(t *testing.T) {
				checkContextParams(t, func(w io.Writer) diag1log.Logger { return diag1log.NewFromCreator(w, provider.NewLogger) },
					diag1log.WithLogger, diag1log.WithLoggerParams, diag1log.FromContext,
					func(value string) diag1log.Param { return diag1log.UnsafeParam("key", value) },
					func(logger diag1log.Logger) {
						logger.Diagnostic(logging.NewDiagnosticFromGeneric(logging.GenericDiagnostic{DiagnosticType: "test", Value: "value"}))
					}, "unsafeParams")
			})
		})
	}
}

func checkContextParams[L, P any](t *testing.T, create func(io.Writer) L, withLogger func(context.Context, L) context.Context,
	withParams func(context.Context, ...P) context.Context, fromContext func(context.Context) L,
	param func(string) P, emit func(L), field string,
) {
	t.Helper()
	var buf bytes.Buffer
	ctx := context.Background()
	require.True(t, ctx == withParams(ctx))
	ctx = withLogger(ctx, create(&buf))
	ctx = withParams(ctx, param("stored"))
	child := withParams(ctx, param("override"))
	for _, tc := range []struct {
		ctx  context.Context
		want string
	}{{child, "override"}, {ctx, "stored"}, {child, "override"}} {
		buf.Reset()
		emit(fromContext(tc.ctx))
		var record map[string]any
		require.NoError(t, json.Unmarshal(buf.Bytes(), &record))
		require.Equal(t, tc.want, record[field].(map[string]any)["key"])
	}
}
