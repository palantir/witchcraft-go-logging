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

package wlogzerolog_test

import (
	"io"
	"testing"

	"github.com/palantir/witchcraft-go-logging/wlog"
	wlogzerolog "github.com/palantir/witchcraft-go-logging/wlog-zerolog"
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
)

func TestSvc1Log(t *testing.T) {
	svc1logtests.JSONTestSuite(t, func(w io.Writer, level wlog.LogLevel, origin string) svc1log.Logger {
		return svc1log.NewFromCreator(
			w,
			level,
			wlogzerolog.LoggerProvider().NewLeveledLogger,
			svc1log.Origin(origin),
		)
	})
}

func TestReq2Log(t *testing.T) {
	req2logtests.JSONTestSuite(t, func(w io.Writer, params ...req2log.LoggerCreatorParam) req2log.Logger {
		allParams := append([]req2log.LoggerCreatorParam{
			req2log.Creator(wlogzerolog.LoggerProvider().NewLogger),
		}, params...)
		return req2log.New(
			w,
			allParams...,
		)
	})
}

func TestEvt2Log(t *testing.T) {
	evt2logtests.JSONTestSuite(t, func(w io.Writer) evt2log.Logger {
		return evt2log.NewFromCreator(
			w,
			wlogzerolog.LoggerProvider().NewLogger,
		)
	})
}

func TestTrc1Log(t *testing.T) {
	trc1logtests.JSONTestSuite(t, func(w io.Writer) trc1log.Logger {
		return trc1log.NewFromCreator(
			w,
			wlogzerolog.LoggerProvider().NewLogger,
		)
	})
}

func TestMetric1Log(t *testing.T) {
	metric1logtests.JSONTestSuite(t, func(w io.Writer) metric1log.Logger {
		return metric1log.NewFromCreator(
			w,
			wlogzerolog.LoggerProvider().NewLogger,
		)
	})
}

func TestAudit2Log(t *testing.T) {
	audit2logtests.JSONTestSuite(t, func(w io.Writer) audit2log.Logger {
		return audit2log.NewFromCreator(
			w,
			wlogzerolog.LoggerProvider().NewLogger,
		)
	})
}

func TestDiag1Log(t *testing.T) {
	diag1logtests.JSONTestSuite(t, func(w io.Writer) diag1log.Logger {
		return diag1log.NewFromCreator(
			w,
			wlogzerolog.LoggerProvider().NewLogger,
		)
	})
}

func TestWrapped1Req2Log(t *testing.T) {
	entityName := "entity"
	entityVersion := "version"
	wrapped1logtests.Req2LogJSONTestSuite(t, entityName, entityVersion, func(w io.Writer, params ...req2log.LoggerCreatorParam) req2log.Logger {
		allParams := append([]req2log.LoggerCreatorParam{
			req2log.Creator(wlogzerolog.LoggerProvider().NewLogger),
		}, params...)
		return wrapped1log.NewFromProvider(w, wlog.InfoLevel, wlogzerolog.LoggerProvider(), entityName, entityVersion).Request(allParams...)
	})
}

func TestWrapped1Svc1Log(t *testing.T) {
	entityName := "entity"
	entityVersion := "version"
	wrapped1logtests.Svc1LogJSONTestSuite(
		t,
		entityName,
		entityVersion,
		func(w io.Writer, level wlog.LogLevel, origin string) svc1log.Logger {
			return wrapped1log.NewFromProvider(w, level, wlogzerolog.LoggerProvider(), entityName, entityVersion).Service(svc1log.Origin(origin))
		})
}

func TestWrapped1Trc1Log(t *testing.T) {
	entityName := "entity"
	entityVersion := "version"
	wrapped1logtests.Trc1LogJSONTestSuite(
		t,
		entityName,
		entityVersion,
		func(w io.Writer) trc1log.Logger {
			return wrapped1log.NewFromProvider(w, wlog.InfoLevel, wlogzerolog.LoggerProvider(), entityName, entityVersion).Trace()
		})
}
