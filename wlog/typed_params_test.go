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
	"encoding/json"
	"maps"
	"strconv"
	"testing"
	"time"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit2log"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit3log"
	"github.com/palantir/witchcraft-go-logging/wlog/diaglog/diag1log"
	"github.com/palantir/witchcraft-go-logging/wlog/evtlog/evt2log"
	"github.com/palantir/witchcraft-go-logging/wlog/metriclog/metric1log"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/palantir/witchcraft-go-logging/wlog/trclog/trc1log"
	"github.com/palantir/witchcraft-go-logging/wlog/wrappedlog/wrapped1log"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
	"github.com/stretchr/testify/require"
)

func TestAudit3EventIDOverrides(t *testing.T) {
	const validID = "c15487b9-ff6a-4bb1-8c25-2433a185c438"
	for name, provider := range jsonProviders() {
		for _, wrapped := range []bool{false, true} {
			for _, tc := range []struct{ name, stored, supplied string }{
				{"invalid replaces valid", validID, "invalid"},
				{"partially valid replaces valid", validID, "c15487b9-ff6a-4bb1-8c25-2433a185c43g"},
				{"empty clears valid", validID, ""},
				{"valid replaces valid", validID, "46ac025f-5d70-4a79-8e8e-270da9635a43"},
				{"valid replaces invalid", "invalid", validID},
			} {
				t.Run(name+"/wrapped="+strconv.FormatBool(wrapped)+"/"+tc.name, func(t *testing.T) {
					var buf bytes.Buffer
					var logger audit3log.Logger
					if wrapped {
						logger = wrapped1log.NewFromProvider(&buf, wlog.InfoLevel, provider, "entity", "version").AuditV3()
					} else {
						logger = audit3log.NewFromCreator(&buf, provider.NewLogger)
					}
					logger = audit3log.WithParams(logger, audit3log.EventID(tc.stored))
					for range 2 {
						buf.Reset()
						logger.Audit("test", audit3log.AuditResultSuccess, audit3log.EventID(tc.supplied))
						var record map[string]any
						require.NoError(t, json.Unmarshal(buf.Bytes(), &record))
						if wrapped {
							record = record["payload"].(map[string]any)["auditLogV3"].(map[string]any)
						}
						if tc.supplied == "" {
							require.NotContains(t, record, "eventId")
						} else {
							require.Equal(t, tc.supplied, record["eventId"])
						}
					}
				})
			}
		}
	}
}

func TestApplyTypedParamPreservesMetadata(t *testing.T) {
	for name, apply := range map[string]func(wlog.LogEntry){
		"service": func(e wlog.LogEntry) { svc1log.ApplyParam(nil, e); svc1log.ApplyParam(svc1log.UID("user"), e) },
		"event":   func(e wlog.LogEntry) { evt2log.ApplyParam(nil, e); evt2log.ApplyParam(evt2log.UID("user"), e) },
		"metric":  func(e wlog.LogEntry) { metric1log.ApplyParam(nil, e); metric1log.ApplyParam(metric1log.UID("user"), e) },
		"trace":   func(e wlog.LogEntry) { trc1log.ApplyParam(nil, e); trc1log.ApplyParam(trc1log.UID("user"), e) },
		"audit":   func(e wlog.LogEntry) { audit2log.ApplyParam(nil, e); audit2log.ApplyParam(audit2log.UID("user"), e) },
		"diagnostic": func(e wlog.LogEntry) {
			diag1log.ApplyParam(nil, e)
			diag1log.ApplyParam(diag1log.UnsafeParam("key", "value"), e)
		},
	} {
		t.Run(name, func(t *testing.T) {
			entry := wlog.NewMapLogEntry()
			metadata := map[string]any{}
			for _, key := range []string{"type", "time", "level", "eventName", "metricName", "metricType", "name", "result", "span", "diagnostic"} {
				entry.StringValue(key, "existing")
				metadata[key] = "existing"
			}
			apply(entry)
			want := maps.Clone(metadata)
			if name == "diagnostic" {
				want["unsafeParams"] = map[string]any{"key": "value"}
			} else {
				want["uid"] = "user"
			}
			require.Equal(t, want, entry.AllValues())
		})
	}
}

func TestAuditNilParams(t *testing.T) {
	for name, provider := range jsonProviders() {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			audit2 := audit2log.NewFromCreator(&buf, provider.NewLogger)
			audit2log.WithParams(audit2, nil, audit2log.UID("stored"), nil).Audit("test", audit2log.AuditResultSuccess, nil, audit2log.UID("supplied"), nil)
			var record map[string]any
			require.NoError(t, json.Unmarshal(buf.Bytes(), &record))
			require.Equal(t, "supplied", record["uid"])
			buf.Reset()
			audit3 := audit3log.NewFromCreator(&buf, provider.NewLogger)
			audit3log.WithParams(audit3, nil, audit3log.UID("stored"), nil).Audit("test", audit3log.AuditResultSuccess,
				nil, audit3log.EventID("c15487b9-ff6a-4bb1-8c25-2433a185c438"), audit3log.UID("supplied"), nil)
			require.NoError(t, json.Unmarshal(buf.Bytes(), &record))
			require.Equal(t, "supplied", record["uid"])
		})
	}
}

func TestApplyMetricTimeParam(t *testing.T) {
	for _, timestamp := range []time.Time{{}, time.Date(2026, 10, 23, 1, 2, 3, 0, time.UTC)} {
		entry := wlog.NewMapLogEntry()
		entry.StringValue("type", "metric.1")
		entry.StringValue("time", "existing")
		metric1log.ApplyParam(metric1log.Time(timestamp), entry)
		require.Equal(t, map[string]any{"type": "metric.1", "time": timestamp.Format(time.RFC3339Nano)}, entry.AllValues())
	}
}

func TestTypedEmptyStringsClearDefaults(t *testing.T) {
	for name, provider := range jsonProviders() {
		t.Run(name, func(t *testing.T) {
			t.Run("service", func(t *testing.T) {
				var buf bytes.Buffer
				logger := svc1log.NewFromCreator(&buf, wlog.InfoLevel, provider.NewLeveledLogger)
				checkEmptyStrings(t, &buf, []func(string) svc1log.Param{svc1log.UID, svc1log.SID, svc1log.TokenID, svc1log.OrgID, svc1log.TraceID, svc1log.Origin}, func(base, params []svc1log.Param) {
					svc1log.WithParams(logger, base...).Info("message", params...)
				})
			})
			t.Run("event", func(t *testing.T) {
				var buf bytes.Buffer
				logger := evt2log.NewFromCreator(&buf, provider.NewLogger)
				checkEmptyStrings(t, &buf, []func(string) evt2log.Param{evt2log.UID, evt2log.SID, evt2log.TokenID, evt2log.OrgID, evt2log.TraceID}, func(base, params []evt2log.Param) {
					evt2log.WithParams(logger, base...).Event("event", params...)
				})
			})
			t.Run("metric", func(t *testing.T) {
				var buf bytes.Buffer
				logger := metric1log.NewFromCreator(&buf, provider.NewLogger)
				checkEmptyStrings(t, &buf, []func(string) metric1log.Param{metric1log.UID, metric1log.SID, metric1log.TokenID, metric1log.OrgID}, func(base, params []metric1log.Param) {
					metric1log.WithParams(logger, base...).Metric("metric", "gauge", params...)
				})
			})
			t.Run("trace", func(t *testing.T) {
				var buf bytes.Buffer
				logger := trc1log.NewFromCreator(&buf, provider.NewLogger)
				checkEmptyStrings(t, &buf, []func(string) trc1log.Param{trc1log.UID, trc1log.SID, trc1log.TokenID, trc1log.OrgID}, func(base, params []trc1log.Param) {
					trc1log.WithParams(logger, base...).Log(wtracing.SpanModel{}, params...)
				})
			})
			t.Run("audit2", func(t *testing.T) {
				var buf bytes.Buffer
				logger := audit2log.NewFromCreator(&buf, provider.NewLogger)
				checkEmptyStrings(t, &buf, []func(string) audit2log.Param{audit2log.UID, audit2log.SID, audit2log.TokenID, audit2log.OrgID, audit2log.TraceID, audit2log.Origin}, func(base, params []audit2log.Param) {
					audit2log.WithParams(logger, base...).Audit("audit", audit2log.AuditResultSuccess, params...)
				})
			})
			t.Run("audit3", func(t *testing.T) {
				var buf bytes.Buffer
				logger := audit3log.NewFromCreator(&buf, provider.NewLogger)
				checkEmptyStrings(t, &buf, []func(string) audit3log.Param{audit3log.UID, audit3log.SID, audit3log.TokenID, audit3log.OrgID, audit3log.TraceID, audit3log.Origin}, func(base, params []audit3log.Param) {
					logger := audit3log.WithParams(logger, audit3log.EventID("c15487b9-ff6a-4bb1-8c25-2433a185c438"))
					audit3log.WithParams(logger, base...).Audit("audit", audit3log.AuditResultSuccess, params...)
				})
			})
		})
	}
}

func checkEmptyStrings[P any](t *testing.T, buf *bytes.Buffer, constructors []func(string) P, emit func(base, params []P)) {
	t.Helper()
	var base, empty []P
	for _, constructor := range constructors {
		base = append(base, constructor("base"))
		empty = append(empty, constructor(""))
	}
	for _, clear := range []bool{false, true, false} {
		buf.Reset()
		if clear {
			emit(base, empty)
		} else {
			emit(base, nil)
		}
		var record map[string]any
		require.NoError(t, json.Unmarshal(buf.Bytes(), &record))
		keys := []string{"uid", "sid", "tokenId", "orgId", "traceId", "origin"}
		for _, key := range keys[:len(constructors)] {
			if clear {
				require.NotContains(t, record, key)
			} else {
				require.Equal(t, "base", record[key], key)
			}
		}
	}
}

func TestServiceLegacyToParams(t *testing.T) {
	for _, tc := range []struct {
		level func() wlog.Param
		want  string
	}{
		{svc1log.DebugLevelParam, svc1log.LevelDebugValue},
		{svc1log.InfoLevelParam, svc1log.LevelInfoValue},
		{svc1log.WarnLevelParam, svc1log.LevelWarnValue},
		{svc1log.ErrorLevelParam, svc1log.LevelErrorValue},
	} {
		t.Run(tc.want, func(t *testing.T) {
			entry := wlog.NewMapLogEntry()
			params := svc1log.ToParams(tc.level(), []svc1log.Param{svc1log.UID("user"), svc1log.UID("")})
			wlog.ApplyParams(entry, params)
			require.Equal(t, tc.want, entry.AllValues()["level"])
			require.Equal(t, "service.1", entry.AllValues()["type"])
			require.NotContains(t, entry.AllValues(), "uid")
			require.NotContains(t, entry.AllValues(), "message")
			_, err := time.Parse(time.RFC3339Nano, entry.AllValues()["time"].(string))
			require.NoError(t, err)
		})
	}
}
