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

package benchmarks

import (
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/palantir/witchcraft-go-logging/wlog"
	wlogzap "github.com/palantir/witchcraft-go-logging/wlog-zap"
	wlogzerolog "github.com/palantir/witchcraft-go-logging/wlog-zerolog"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit2log"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit2log/audit2logtests"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit3log"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit3log/audit3logtests"
	"github.com/palantir/witchcraft-go-logging/wlog/diaglog/diag1log"
	"github.com/palantir/witchcraft-go-logging/wlog/diaglog/diag1log/diag1logtests"
	"github.com/palantir/witchcraft-go-logging/wlog/evtlog/evt2log"
	"github.com/palantir/witchcraft-go-logging/wlog/evtlog/evt2log/evt2logtests"
	"github.com/palantir/witchcraft-go-logging/wlog/metriclog/metric1log"
	"github.com/palantir/witchcraft-go-logging/wlog/metriclog/metric1log/metric1logtests"
	"github.com/palantir/witchcraft-go-logging/wlog/reqlog/req2log"
	"github.com/palantir/witchcraft-go-logging/wlog/reqlog/req2log/req2logtests"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log/svc1logtests"
	"github.com/palantir/witchcraft-go-logging/wlog/trclog/trc1log"
	"github.com/palantir/witchcraft-go-logging/wlog/trclog/trc1log/trc1logtests"
	"github.com/palantir/witchcraft-go-logging/wlog/wrappedlog/wrapped1log"
	"github.com/palantir/witchcraft-go-logging/wlog/wrappedlog/wrapped1log/wrapped1logtests"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
	"github.com/palantir/witchcraft-go-tracing/wzipkin"
	"github.com/stretchr/testify/require"
)

func BenchmarkAudit2Log(b *testing.B) {
	for _, tc := range audit2logtests.TestCases() {
		b.Run(tc.Name, func(b *testing.B) {
			params := tc.Params()
			RunBenchmarks(b, func(b *testing.B, provider wlog.LoggerProvider) {
				b.ReportAllocs()
				logger := audit2log.NewFromCreator(io.Discard, provider.NewLogger)
				benchmarkLogCalls(b, func() {
					logger.Audit(tc.AuditName, tc.AuditResult, params...)
				})
			})
		})
	}
}

func BenchmarkAudit3Log(b *testing.B) {
	for _, tc := range audit3logtests.TestCases() {
		b.Run(tc.TestCaseName, func(b *testing.B) {
			params := tc.Params()
			RunBenchmarks(b, func(b *testing.B, provider wlog.LoggerProvider) {
				b.ReportAllocs()
				logger := audit3log.NewFromCreator(io.Discard, provider.NewLogger)
				benchmarkLogCalls(b, func() {
					logger.Audit(tc.Name, tc.Result, params...)
				})
			})
		})
	}
}

func BenchmarkDiag1Log(b *testing.B) {
	for _, tc := range diag1logtests.TestCases() {
		b.Run(tc.Name, func(b *testing.B) {
			RunBenchmarks(b, func(b *testing.B, provider wlog.LoggerProvider) {
				b.ReportAllocs()
				logger := diag1log.NewFromCreator(io.Discard, provider.NewLogger)
				benchmarkLogCalls(b, func() {
					logger.Diagnostic(tc.Diagnostic, diag1log.UnsafeParams(tc.UnsafeParams))
				})
			})
		})
	}
}

func BenchmarkEvt2Log(b *testing.B) {
	for _, tc := range evt2logtests.TestCases() {
		params := tc.Params()
		b.Run(tc.Name, func(b *testing.B) {
			RunBenchmarks(b, func(b *testing.B, provider wlog.LoggerProvider) {
				b.ReportAllocs()
				logger := evt2log.NewFromCreator(io.Discard, provider.NewLogger)
				benchmarkLogCalls(b, func() {
					logger.Event(tc.EventName, params...)
				})
			})
		})
	}
}

func BenchmarkMetric1Log(b *testing.B) {
	for _, tc := range metric1logtests.TestCases() {
		params := tc.Params()
		b.Run(tc.Name, func(b *testing.B) {
			RunBenchmarks(b, func(b *testing.B, provider wlog.LoggerProvider) {
				b.ReportAllocs()
				logger := metric1log.NewFromCreator(io.Discard, provider.NewLogger)
				benchmarkLogCalls(b, func() {
					logger.Metric(tc.MetricName, tc.MetricType, params...)
				})
			})
		})
	}
}

func BenchmarkReq2Log(b *testing.B) {
	for _, tc := range req2logtests.TestCases() {
		b.Run(tc.Name, func(b *testing.B) {
			req := req2log.Request{
				Request: req2logtests.GenerateRequest(map[string]string{
					"Authorization":      "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ2cDlrWFZMZ1NlbTZNZHN5a25ZVjJ3PT0iLCJzaWQiOiJyVTFLNW1XdlRpcVJvODlBR3NzZFRBPT0iLCJqdGkiOiJrbmY1cjQyWlFJcVU3L1VlZ3I0ditBPT0iLCJvcmciOiJDWmpsY3pIWFNabUwrUXZGOUZrdHVRPT0ifQ.GMqKu_zrkgNR5I-jAWdR6x0G2gObVYRbqw7iJJatI4A",
					"FooHeaderParamName": "fooHeaderParamVal",
				}, tc.ExtraQueryParams, tc.ExtraHeaderParams),
				RouteInfo:      req2log.RouteInfo{},
				ResponseStatus: http.StatusOK,
				ResponseSize:   int64(100),
				Duration:       1 * time.Second,
			}
			RunBenchmarks(b, func(b *testing.B, provider wlog.LoggerProvider) {
				b.ReportAllocs()
				logger := req2log.New(io.Discard,
					req2log.Creator(provider.NewLogger),
					req2log.SafeQueryParams(tc.SafeQueryParams...),
					req2log.SafeHeaderParams(tc.SafeHeaderParams...),
					req2log.ForbiddenHeaderParams(tc.ForbiddenHeaderParams...),
				)
				benchmarkLogCalls(b, func() {
					logger.Request(req)
				})
			})
		})
	}
}

func BenchmarkSvc1Log(b *testing.B) {
	for _, tc := range svc1logtests.TestCases() {
		b.Run(tc.Name, func(b *testing.B) {
			RunBenchmarks(b, func(b *testing.B, provider wlog.LoggerProvider) {
				b.ReportAllocs()
				logger := svc1log.NewFromCreator(io.Discard, wlog.InfoLevel, provider.NewLeveledLogger)
				if tc.Origin != "" {
					logger = svc1log.WithParams(logger, svc1log.Origin(tc.Origin))
				}
				benchmarkLogCalls(b, func() {
					logger.Info(tc.Message, tc.LogParams...)
				})
			})
		})
	}
}

func BenchmarkTrc1Log(b *testing.B) {
	tracer, err := wzipkin.NewTracer(wtracing.NewNoopReporter())
	if err != nil {
		b.Fatal(err)
	}
	clientSpan := tracer.StartSpan("testOp", wtracing.WithKind(wtracing.Client))
	defer clientSpan.Finish()

	for _, tc := range trc1logtests.TestCases(clientSpan) {
		b.Run(tc.Name, func(b *testing.B) {
			RunBenchmarks(b, func(b *testing.B, provider wlog.LoggerProvider) {
				b.ReportAllocs()
				logger := trc1log.NewFromCreator(io.Discard, provider.NewLogger)
				tracer, err := wzipkin.NewTracer(
					logger,
					wtracing.WithLocalEndpoint(&wtracing.Endpoint{
						ServiceName: "testService",
						IPv4:        net.IPv4(127, 0, 0, 1),
						Port:        1234,
					}),
				)
				require.NoError(b, err)
				spanOptions := append([]wtracing.SpanOption{wtracing.WithParent(clientSpan)}, tc.SpanOptions...)
				benchmarkLogCalls(b, func() {
					span := tracer.StartSpan("testOp", spanOptions...)
					for key, value := range tc.Tags {
						span.Tag(key, value)
					}
					// Finish() triggers logging
					span.Finish()
				})
			})
		})
	}
}

func BenchmarkWrapped1Log(b *testing.B) {
	// Providers with wrapped object support emit directly; others use the map path.
	tc := wrapped1logtests.Svc1TestCases("bench-entity", "1")[0]
	b.Run(tc.Name, func(b *testing.B) {
		RunBenchmarks(b, func(b *testing.B, provider wlog.LoggerProvider) {
			b.ReportAllocs()
			logger := wrapped1log.NewFromProvider(io.Discard, wlog.InfoLevel, provider, "bench-entity", "1").Service()
			benchmarkLogCalls(b, func() {
				logger.Info(tc.Message, tc.LogParams...)
			})
		})
	})
	b.Run("audit.3", func(b *testing.B) {
		fixture := wrapped1logtests.Audit3TestCases("bench-entity", "1")[0]
		// Typed emission requires a valid event ID.
		params := append(fixture.Params(), audit3log.EventID("c15487b9-ff6a-4bb1-8c25-2433a185c438"))
		RunBenchmarks(b, func(b *testing.B, provider wlog.LoggerProvider) {
			logger := wrapped1log.NewFromProvider(io.Discard, wlog.InfoLevel, provider, "bench-entity", "1").AuditV3()
			benchmarkLogCalls(b, func() {
				logger.Audit(fixture.AuditName, fixture.AuditResult, params...)
			})
		})
	})
}

func RunBenchmarks(b *testing.B, benchmark func(*testing.B, wlog.LoggerProvider)) {
	b.Run("noop", func(b *testing.B) { benchmark(b, wlog.NewNoopLoggerProvider()) })
	b.Run("json", func(b *testing.B) { benchmark(b, wlog.NewJSONMarshalLoggerProvider()) })
	b.Run("zap", func(b *testing.B) { benchmark(b, wlogzap.LoggerProvider()) })
	b.Run("zerolog", func(b *testing.B) { benchmark(b, wlogzerolog.LoggerProvider()) })
}

// Parallel calls share one logger and writer to include provider contention.
func benchmarkLogCalls(b *testing.B, emit func()) {
	b.Run("serial", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			emit()
		}
	})
	b.Run("parallel", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				emit()
			}
		})
	})
}
