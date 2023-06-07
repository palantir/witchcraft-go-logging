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

package wlogglog_test

import (
	"flag"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/palantir/witchcraft-go-logging/wlog"
	wlogglog "github.com/palantir/witchcraft-go-logging/wlog-glog"
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
	"github.com/palantir/witchcraft-go-logging/wlog/wrappedlog/wrapped1log"
	"github.com/palantir/witchcraft-go-logging/wlog/wrappedlog/wrapped1log/wrapped1logtests"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
	"github.com/palantir/witchcraft-go-tracing/wzipkin"
	"github.com/stretchr/testify/require"
)

func TestSvc1Log(t *testing.T) {
	os.Args = []string{
		os.Args[0],
		"-logtostderr=true",
	}
	flag.Parse()

	for _, tc := range svc1logtests.TestCases() {
		// TODO: test output
		logger := svc1log.NewFromCreator(
			os.Stdout,
			wlog.DebugLevel,
			wlogglog.LoggerProvider().NewLeveledLogger,
			svc1log.Origin(tc.Origin),
		)
		logger.Debug(tc.Message, tc.LogParams...)
	}
}

func TestReq2Log(t *testing.T) {
	os.Args = []string{
		os.Args[0],
		"-logtostderr=true",
	}
	flag.Parse()

	for _, tc := range req2logtests.TestCases() {
		// TODO: test output
		logger := req2log.New(
			os.Stdout,
			req2log.Creator(wlogglog.LoggerProvider().NewLogger),
		)
		req := req2logtests.GenerateRequest(map[string]string{
			"Authorization":      "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ2cDlrWFZMZ1NlbTZNZHN5a25ZVjJ3PT0iLCJzaWQiOiJyVTFLNW1XdlRpcVJvODlBR3NzZFRBPT0iLCJqdGkiOiJrbmY1cjQyWlFJcVU3L1VlZ3I0ditBPT0iLCJvcmciOiJDWmpsY3pIWFNabUwrUXZGOUZrdHVRPT0ifQ.GMqKu_zrkgNR5I-jAWdR6x0G2gObVYRbqw7iJJatI4A",
			"FooHeaderParamName": "fooHeaderParamVal",
		}, tc.ExtraQueryParams, tc.ExtraHeaderParams)

		logger.Request(req2log.Request{
			Request:        req,
			RouteInfo:      req2log.RouteInfo{},
			ResponseStatus: http.StatusOK,
			ResponseSize:   int64(100),
			Duration:       1 * time.Second,
		})
	}
}

func TestEvt2Log(t *testing.T) {
	os.Args = []string{
		os.Args[0],
		"-logtostderr=true",
	}
	flag.Parse()

	for _, tc := range evt2logtests.TestCases() {
		// TODO: test output
		logger := evt2log.NewFromCreator(
			os.Stdout,
			wlogglog.LoggerProvider().NewLogger,
		)

		logger.Event(tc.EventName, tc.Params()...)
	}
}

func TestTrc1Log(t *testing.T) {
	os.Args = []string{
		os.Args[0],
		"-logtostderr=true",
	}
	flag.Parse()

	tracer, err := wzipkin.NewTracer(wtracing.NewNoopReporter())
	require.NoError(t, err)
	clientSpan := tracer.StartSpan("testOp", wtracing.WithKind(wtracing.Client))
	defer clientSpan.Finish()

	for _, tc := range trc1logtests.TestCases(clientSpan) {
		// TODO: test output
		logger := trc1log.NewFromCreator(
			os.Stdout,
			wlogglog.LoggerProvider().NewLogger,
		)

		tracer, err := wzipkin.NewTracer(
			logger,
			wtracing.WithLocalEndpoint(&wtracing.Endpoint{
				ServiceName: "testService",
				IPv4:        net.IPv4(127, 0, 0, 1),
				Port:        1234,
			}),
		)
		require.NoError(t, err)
		span := tracer.StartSpan("testOp", append([]wtracing.SpanOption{wtracing.WithParent(clientSpan)}, tc.SpanOptions...)...)
		// Finish() triggers logging
		span.Finish()
	}
}

func TestMetric1Log(t *testing.T) {
	os.Args = []string{
		os.Args[0],
		"-logtostderr=true",
	}
	flag.Parse()

	for _, tc := range metric1logtests.TestCases() {
		// TODO: test output
		logger := metric1log.NewFromCreator(
			os.Stdout,
			wlogglog.LoggerProvider().NewLogger,
		)

		logger.Metric(tc.MetricName, tc.MetricType, tc.Params()...)
	}
}

func TestAudit2Log(t *testing.T) {
	os.Args = []string{
		os.Args[0],
		"-logtostderr=true",
	}
	flag.Parse()

	for _, tc := range audit2logtests.TestCases() {
		// TODO: test output
		logger := audit2log.NewFromCreator(
			os.Stdout,
			wlogglog.LoggerProvider().NewLogger,
		)

		logger.Audit(tc.AuditName, tc.AuditResult, tc.Params()...)
	}
}

func TestDiag1Log(t *testing.T) {
	os.Args = []string{
		os.Args[0],
		"-logtostderr=true",
	}
	flag.Parse()

	for _, tc := range diag1logtests.TestCases() {
		// TODO: test output
		logger := diag1log.NewFromCreator(
			os.Stdout,
			wlogglog.LoggerProvider().NewLogger,
		)

		logger.Diagnostic(
			tc.Diagnostic,
			diag1log.UnsafeParams(tc.UnsafeParams),
		)
	}
}
func TestWrapped1Diag1Log(t *testing.T) {
	os.Args = []string{
		os.Args[0],
		"-logtostderr=true",
	}
	flag.Parse()

	entityName := "entity"
	entityVersion := "version"
	for _, tc := range wrapped1logtests.Diag1TestCases(entityName, entityVersion) {
		// TODO: test output
		logger := wrapped1log.NewFromProvider(
			os.Stdout,
			wlog.InfoLevel,
			wlogglog.LoggerProvider(),
			entityName,
			entityVersion,
		).Diagnostic()
		logger.Diagnostic(
			tc.Diagnostic,
			diag1log.UnsafeParams(tc.UnsafeParams),
		)
	}
}

func TestWrapped1Evt2Log(t *testing.T) {
	os.Args = []string{
		os.Args[0],
		"-logtostderr=true",
	}
	flag.Parse()

	entityName := "entity"
	entityVersion := "version"
	for _, tc := range wrapped1logtests.Evt2TestCases(entityName, entityVersion) {
		// TODO: test output
		logger := wrapped1log.NewFromProvider(
			os.Stdout,
			wlog.InfoLevel,
			wlogglog.LoggerProvider(),
			entityName,
			entityVersion,
		).Event()

		logger.Event(tc.EventName, tc.Params()...)
	}
}

func TestWrapped1Audit2Log(t *testing.T) {
	os.Args = []string{
		os.Args[0],
		"-logtostderr=true",
	}
	flag.Parse()

	entityName := "entity"
	entityVersion := "version"
	for _, tc := range wrapped1logtests.Audit2TestCases(entityName, entityVersion) {
		// TODO: test output
		logger := wrapped1log.NewFromProvider(
			os.Stdout,
			wlog.InfoLevel,
			wlogglog.LoggerProvider(),
			entityName,
			entityVersion,
		).Audit()

		logger.Audit(tc.AuditName, tc.AuditResult, tc.Params()...)
	}
}

func TestWrapped1Metric1Log(t *testing.T) {
	os.Args = []string{
		os.Args[0],
		"-logtostderr=true",
	}
	flag.Parse()

	entityName := "entity"
	entityVersion := "version"
	for _, tc := range wrapped1logtests.Metric1TestCases(entityName, entityVersion) {
		// TODO: test output
		logger := wrapped1log.NewFromProvider(
			os.Stdout,
			wlog.InfoLevel,
			wlogglog.LoggerProvider(),
			entityName,
			entityVersion,
		).Metric()

		logger.Metric(tc.MetricName, tc.MetricType, tc.Params()...)
	}
}

func TestWrapped1Req2Log(t *testing.T) {
	os.Args = []string{
		os.Args[0],
		"-logtostderr=true",
	}
	flag.Parse()

	entityName := "entity"
	entityVersion := "version"
	for _, tc := range wrapped1logtests.Req2TestCases(entityName, entityVersion) {
		// TODO: test output
		logger := wrapped1log.NewFromProvider(
			os.Stdout,
			wlog.InfoLevel,
			wlogglog.LoggerProvider(),
			entityName,
			entityVersion,
		).Request()
		req := req2logtests.GenerateRequest(map[string]string{
			"Authorization":      "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ2cDlrWFZMZ1NlbTZNZHN5a25ZVjJ3PT0iLCJzaWQiOiJyVTFLNW1XdlRpcVJvODlBR3NzZFRBPT0iLCJqdGkiOiJrbmY1cjQyWlFJcVU3L1VlZ3I0ditBPT0iLCJvcmciOiJDWmpsY3pIWFNabUwrUXZGOUZrdHVRPT0ifQ.GMqKu_zrkgNR5I-jAWdR6x0G2gObVYRbqw7iJJatI4A",
			"FooHeaderParamName": "fooHeaderParamVal",
		}, tc.ExtraQueryParams, tc.ExtraHeaderParams)

		logger.Request(req2log.Request{
			Request:        req,
			RouteInfo:      req2log.RouteInfo{},
			ResponseStatus: http.StatusOK,
			ResponseSize:   int64(100),
			Duration:       1 * time.Second,
		})
	}
}

func TestWrapped1Svc1Log(t *testing.T) {
	os.Args = []string{
		os.Args[0],
		"-logtostderr=true",
	}
	flag.Parse()

	entityName := "entity"
	entityVersion := "version"
	for _, tc := range wrapped1logtests.Svc1TestCases(entityName, entityVersion) {
		// TODO: test output
		logger := wrapped1log.NewFromProvider(
			os.Stdout,
			wlog.DebugLevel,
			wlogglog.LoggerProvider(),
			entityName,
			entityVersion,
		).Service(svc1log.Origin(tc.Origin))

		logger.Info(tc.Message, tc.LogParams...)
	}
}

func TestWrapped1Trc1Log(t *testing.T) {
	os.Args = []string{
		os.Args[0],
		"-logtostderr=true",
	}
	flag.Parse()

	tracer, err := wzipkin.NewTracer(wtracing.NewNoopReporter())
	require.NoError(t, err)
	clientSpan := tracer.StartSpan("testOp", wtracing.WithKind(wtracing.Client))
	defer clientSpan.Finish()

	entityName := "entity"
	entityVersion := "version"
	for _, tc := range wrapped1logtests.Trc1TestCases(entityName, entityVersion, clientSpan) {
		// TODO: test output
		logger := wrapped1log.NewFromProvider(
			os.Stdout,
			wlog.DebugLevel,
			wlogglog.LoggerProvider(),
			entityName,
			entityVersion,
		).Trace()

		tracer, err := wzipkin.NewTracer(
			logger,
			wtracing.WithLocalEndpoint(&wtracing.Endpoint{
				ServiceName: "testService",
				IPv4:        net.IPv4(127, 0, 0, 1),
				Port:        1234,
			}),
		)
		require.NoError(t, err)
		span := tracer.StartSpan("testOp", append([]wtracing.SpanOption{wtracing.WithParent(clientSpan)}, tc.SpanOptions...)...)
		// Finish() triggers logging
		span.Finish()
	}
}
