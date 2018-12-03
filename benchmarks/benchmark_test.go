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
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog-glog"
	"github.com/palantir/witchcraft-go-logging/wlog-zap"
	"github.com/palantir/witchcraft-go-logging/wlog-zerolog"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit2log"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit2log/audit2logtests"
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
				logger := audit2log.NewFromCreator(ioutil.Discard, provider.NewLogger)
				for n := 0; n < b.N; n++ {
					logger.Audit(tc.AuditName, tc.AuditResult, params...)
				}
			})
		})
	}
}

func BenchmarkDiag1Log(b *testing.B) {
	for _, tc := range diag1logtests.TestCases() {
		b.Run(tc.Name, func(b *testing.B) {
			RunBenchmarks(b, func(b *testing.B, provider wlog.LoggerProvider) {
				b.ReportAllocs()
				logger := diag1log.NewFromCreator(ioutil.Discard, provider.NewLogger)
				for n := 0; n < b.N; n++ {
					logger.Diagnostic(tc.Diagnostic, diag1log.UnsafeParams(tc.UnsafeParams))
				}
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
				logger := evt2log.NewFromCreator(ioutil.Discard, provider.NewLogger)
				for n := 0; n < b.N; n++ {
					logger.Event(tc.Name, params...)
				}
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
				logger := metric1log.NewFromCreator(ioutil.Discard, provider.NewLogger)
				for n := 0; n < b.N; n++ {
					logger.Metric(tc.MetricName, tc.MetricType, params...)
				}
			})
		})
	}
}

func BenchmarkReq2Log(b *testing.B) {
	for _, tc := range req2logtests.TestCases() {
		b.Run(tc.Name, func(b *testing.B) {
			req := req2log.Request{
				Request: req2logtests.GenerateRequest(map[string]string{
					"Authorization":      "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ2cDlrWFZMZ1NlbTZNZHN5a25ZVjJ3PT0iLCJzaWQiOiJyVTFLNW1XdlRpcVJvODlBR3NzZFRBPT0iLCJqdGkiOiJrbmY1cjQyWlFJcVU3L1VlZ3I0ditBPT0ifQ.JTD36MhcwmSuvfdCkfSYc-LHOGNA1UQ-0FKLKqdXbF4",
					"FooHeaderParamName": "fooHeaderParamVal",
				}, tc.ExtraQueryParams, tc.ExtraHeaderParams),
				RouteInfo:      req2log.RouteInfo{},
				ResponseStatus: http.StatusOK,
				ResponseSize:   int64(100),
				Duration:       1 * time.Second,
			}
			RunBenchmarks(b, func(b *testing.B, provider wlog.LoggerProvider) {
				b.ReportAllocs()
				logger := req2log.New(ioutil.Discard, req2log.Creator(provider.NewLogger))
				for n := 0; n < b.N; n++ {
					logger.Request(req)
				}
			})
		})
	}
}

func BenchmarkSvc1Log(b *testing.B) {
	for _, tc := range svc1logtests.TestCases() {
		b.Run(tc.Name, func(b *testing.B) {
			RunBenchmarks(b, func(b *testing.B, provider wlog.LoggerProvider) {
				b.ReportAllocs()
				logger := svc1log.NewFromCreator(ioutil.Discard, wlog.InfoLevel, provider.NewLeveledLogger)
				for n := 0; n < b.N; n++ {
					logger.Info(tc.Message, tc.LogParams...)
				}
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
				logger := trc1log.NewFromCreator(ioutil.Discard, provider.NewLogger)
				tracer, err := wzipkin.NewTracer(
					logger,
					wtracing.WithLocalEndpoint(&wtracing.Endpoint{
						ServiceName: "testService",
						IPv4:        net.IPv4(127, 0, 0, 1),
						Port:        1234,
					}),
				)
				require.NoError(b, err)
				for n := 0; n < b.N; n++ {
					// Finish() triggers logging
					tracer.StartSpan("testOp", wtracing.WithParent(clientSpan)).Finish()
				}
			})
		})
	}
}

func RunBenchmarks(b *testing.B, benchmark func(*testing.B, wlog.LoggerProvider)) {
	b.Run("noop", func(b *testing.B) { benchmark(b, wlog.NewNoopLoggerProvider()) })
	b.Run("glog", func(b *testing.B) { benchmark(b, wlogglog.LoggerProvider()) })
	b.Run("zap", func(b *testing.B) { benchmark(b, wlogzap.LoggerProvider()) })
	b.Run("zerolog", func(b *testing.B) { benchmark(b, wlogzerolog.LoggerProvider()) })
}
