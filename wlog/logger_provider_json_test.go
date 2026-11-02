// Copyright (c) 2025 Palantir Technologies. All rights reserved.
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
	"io"
	"testing"

	"github.com/palantir/witchcraft-go-logging/wlog"
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
	"github.com/palantir/witchcraft-go-logging/wlog/wlogtests"
	"github.com/palantir/witchcraft-go-logging/wlog/wrappedlog/wrapped1log"
	"github.com/palantir/witchcraft-go-logging/wlog/wrappedlog/wrapped1log/wrapped1logtests"
)

func jsonProviders() map[string]wlog.LoggerProvider {
	return map[string]wlog.LoggerProvider{
		"object": wlog.NewJSONMarshalLoggerProvider(),
		"map":    wlogtests.MapJSONProvider(),
	}
}

func TestJSONProviders(t *testing.T) {
	for name, provider := range jsonProviders() {
		t.Run(name, func(t *testing.T) {
			t.Run("wlog", func(t *testing.T) {
				wlogtests.JSONTestSuite(t, provider.NewLogger)
			})
			t.Run("service", func(t *testing.T) {
				svc1logtests.JSONTestSuite(t, func(w io.Writer, level wlog.LogLevel, origin string) svc1log.Logger {
					return svc1log.NewFromCreator(w, level, provider.NewLeveledLogger, svc1log.Origin(origin))
				})
			})
			t.Run("request", func(t *testing.T) {
				req2logtests.JSONTestSuite(t, func(w io.Writer, params ...req2log.LoggerCreatorParam) req2log.Logger {
					return req2log.NewFromCreator(w, provider.NewLogger, params...)
				})
			})
			t.Run("event", func(t *testing.T) {
				evt2logtests.JSONTestSuite(t, func(w io.Writer) evt2log.Logger {
					return evt2log.NewFromCreator(w, provider.NewLogger)
				})
			})
			t.Run("trace", func(t *testing.T) {
				trc1logtests.JSONTestSuite(t, func(w io.Writer) trc1log.Logger {
					return trc1log.NewFromCreator(w, provider.NewLogger)
				})
			})
			t.Run("metric", func(t *testing.T) {
				metric1logtests.JSONTestSuite(t, func(w io.Writer) metric1log.Logger {
					return metric1log.NewFromCreator(w, provider.NewLogger)
				})
			})
			t.Run("audit2", func(t *testing.T) {
				audit2logtests.JSONTestSuite(t, func(w io.Writer) audit2log.Logger {
					return audit2log.NewFromCreator(w, provider.NewLogger)
				})
			})
			t.Run("diagnostic", func(t *testing.T) {
				diag1logtests.JSONTestSuite(t, func(w io.Writer) diag1log.Logger {
					return diag1log.NewFromCreator(w, provider.NewLogger)
				})
			})
			const entityName, entityVersion = "entity", "version"
			t.Run("wrapped/audit2", func(t *testing.T) {
				wrapped1logtests.Audit2LogJSONTestSuite(t, entityName, entityVersion, func(w io.Writer) audit2log.Logger {
					return wrapped1log.NewFromProvider(w, wlog.InfoLevel, provider, entityName, entityVersion).Audit()
				})
			})
			t.Run("wrapped/diagnostic", func(t *testing.T) {
				wrapped1logtests.Diag1LogJSONTestSuite(t, entityName, entityVersion, func(w io.Writer) diag1log.Logger {
					return wrapped1log.NewFromProvider(w, wlog.InfoLevel, provider, entityName, entityVersion).Diagnostic()
				})
			})
			t.Run("wrapped/event", func(t *testing.T) {
				wrapped1logtests.Evt2LogJSONTestSuite(t, entityName, entityVersion, func(w io.Writer) evt2log.Logger {
					return wrapped1log.NewFromProvider(w, wlog.InfoLevel, provider, entityName, entityVersion).Event()
				})
			})
			t.Run("wrapped/metric", func(t *testing.T) {
				wrapped1logtests.Metric1LogJSONTestSuite(t, entityName, entityVersion, func(w io.Writer) metric1log.Logger {
					return wrapped1log.NewFromProvider(w, wlog.InfoLevel, provider, entityName, entityVersion).Metric()
				})
			})
			t.Run("wrapped/request", func(t *testing.T) {
				wrapped1logtests.Req2LogJSONTestSuite(t, entityName, entityVersion, func(w io.Writer, params ...req2log.LoggerCreatorParam) req2log.Logger {
					return wrapped1log.NewFromProvider(w, wlog.InfoLevel, provider, entityName, entityVersion).Request(params...)
				})
			})
			t.Run("wrapped/service", func(t *testing.T) {
				wrapped1logtests.Svc1LogJSONTestSuite(t, entityName, entityVersion, func(w io.Writer, level wlog.LogLevel, origin string) svc1log.Logger {
					return wrapped1log.NewFromProvider(w, level, provider, entityName, entityVersion).Service(svc1log.Origin(origin))
				})
			})
			t.Run("wrapped/trace", func(t *testing.T) {
				wrapped1logtests.Trc1LogJSONTestSuite(t, entityName, entityVersion, func(w io.Writer) trc1log.Logger {
					return wrapped1log.NewFromProvider(w, wlog.InfoLevel, provider, entityName, entityVersion).Trace()
				})
			})
		})
	}
}
