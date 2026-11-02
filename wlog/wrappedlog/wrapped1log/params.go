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

package wrapped1log

import (
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit2log"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit3log"
	"github.com/palantir/witchcraft-go-logging/wlog/diaglog/diag1log"
	"github.com/palantir/witchcraft-go-logging/wlog/evtlog/evt2log"
	"github.com/palantir/witchcraft-go-logging/wlog/extractor"
	"github.com/palantir/witchcraft-go-logging/wlog/metriclog/metric1log"
	"github.com/palantir/witchcraft-go-logging/wlog/reqlog/req2log"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/palantir/witchcraft-go-logging/wlog/trclog/trc1log"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
)

const (
	TypeValue = "wrapped.1"

	WrappedEntityNameKey    = "entityName"
	WrappedEntityVersionKey = "entityVersion"

	PayloadKey             = "payload"
	PayloadTypeKey         = "type"
	PayloadServiceLogV1    = "serviceLogV1"
	PayloadRequestLogV2    = "requestLogV2"
	PayloadTraceLogV1      = "traceLogV1"
	PayloadEventLogV2      = "eventLogV2"
	PayloadMetricLogV1     = "metricLogV1"
	PayloadAuditLogV2      = "auditLogV2"
	PayloadAuditLogV3      = "auditLogV3"
	PayloadDiagnosticLogV1 = "diagnosticLogV1"
)

type Param interface {
	apply(entry wlog.LogEntry)
}

func ApplyParam(p Param, entry wlog.LogEntry) {
	if p == nil {
		return
	}
	p.apply(entry)
}

type paramFunc func(entry wlog.LogEntry)

func (f paramFunc) apply(entry wlog.LogEntry) {
	f(entry)
}

func audit2PayloadParams(name string, result audit2log.AuditResultType, params []audit2log.Param) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		applyPayload(entry, PayloadAuditLogV2, audit2log.ToParams(name, result, params))
	})
}

func audit3PayloadParams(name string, result audit3log.AuditResultType, params []audit3log.Param) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		applyPayload(entry, PayloadAuditLogV3, audit3log.ToParams(name, result, params))
	})
}

func diag1PayloadParams(diagnostic logging.Diagnostic, params []diag1log.Param) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		applyPayload(entry, PayloadDiagnosticLogV1, diag1log.ToParams(diagnostic, params))
	})
}

func evt2PayloadParams(name string, params []evt2log.Param) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		applyPayload(entry, PayloadEventLogV2, evt2log.ToParams(name, params))
	})
}

func metric1PayloadParams(metricName, metricType string, params []metric1log.Param) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		applyPayload(entry, PayloadMetricLogV1, metric1log.ToParams(metricName, metricType, params))
	})
}

func req2PayloadParams(r req2log.Request, idsExtractor extractor.IDsFromRequest, pathParamPerms, queryParamPerms, headerParamPerms req2log.ParamPerms) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		applyPayload(entry, PayloadRequestLogV2, req2log.ToParams(r, idsExtractor, pathParamPerms, queryParamPerms, headerParamPerms))
	})
}

func svc1PayloadParams(obj logging.ServiceLogV1) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		payloadParams := []wlog.Param{svc1log.ObjectToParam(&obj)}
		// svc1log.ObjectToParam doesn't emit "message" (the direct provider path gets it via the msg arg); emit it
		// here so the embedded service.1 payload includes the message.
		if obj.Message != "" {
			payloadParams = append(payloadParams, wlog.StringParam(svc1log.MessageKey, obj.Message))
		}
		applyPayload(entry, PayloadServiceLogV1, payloadParams)
	})
}

func trc1PayloadParams(span wtracing.SpanModel, params []trc1log.Param) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		applyPayload(entry, PayloadTraceLogV1, trc1log.ToParams(span, params))
	})
}

func applyPayload(entry wlog.LogEntry, payloadType string, params []wlog.Param) {
	inner := wlog.NewMapLogEntry()
	wlog.ApplyParams(inner, params)
	entry.AnyMapValue(PayloadKey, map[string]any{
		PayloadTypeKey: payloadType,
		payloadType:    inner.AllValues(),
	})
}

var wrappedLogTypeParam = wlog.StringParam(wlog.TypeKey, TypeValue)

func wrappedParams(name, version string, payload Param) []wlog.Param {
	return []wlog.Param{
		wrappedLogTypeParam,
		wlog.NewParam(wrappedTypeParams(name, version).apply),
		wlog.NewParam(payload.apply),
	}
}

func wrappedTypeParams(name, version string) Param {
	return paramFunc(func(logger wlog.LogEntry) {
		logger.StringValue(WrappedEntityNameKey, name)
		logger.StringValue(WrappedEntityVersionKey, version)
	})
}
