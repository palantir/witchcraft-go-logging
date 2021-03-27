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
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/extractor"
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

func req2PayloadParams(r req2log.Request, idsExtractor extractor.IDsFromRequest, pathParamPerms, queryParamPerms, headerParamPerms req2log.ParamPerms) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		req2Log := wlog.NewMapLogEntry()
		wlog.ApplyParams(req2Log, req2log.ToParams(r, idsExtractor, pathParamPerms, queryParamPerms, headerParamPerms))
		payload := wlog.NewMapLogEntry()
		payload.StringValue(PayloadTypeKey, PayloadRequestLogV2)
		payload.AnyMapValue(PayloadRequestLogV2, req2Log.AllValues())

		entry.AnyMapValue(PayloadKey, payload.AllValues())
	})
}

func svc1PayloadParams(message string, level wlog.Param, params []svc1log.Param) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		svc1Log := wlog.NewMapLogEntry()
		wlog.ApplyParams(svc1Log, wlog.ParamsWithMessage(message, svc1log.ToParams(level, params)))
		payload := wlog.NewMapLogEntry()
		payload.StringValue(PayloadTypeKey, PayloadServiceLogV1)
		payload.AnyMapValue(PayloadServiceLogV1, svc1Log.AllValues())

		entry.AnyMapValue(PayloadKey, payload.AllValues())
	})
}

func trc1PayloadParams(span wtracing.SpanModel, params ...trc1log.Param) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		trc1Log := wlog.NewMapLogEntry()
		wlog.ApplyParams(trc1Log, trc1log.ToParams(span, params))
		payload := wlog.NewMapLogEntry()
		payload.StringValue(PayloadTypeKey, PayloadTraceLogV1)
		payload.AnyMapValue(PayloadTraceLogV1, trc1Log.AllValues())

		entry.AnyMapValue(PayloadKey, payload.AllValues())
	})
}

func wrappedTypeParams(name, version string) Param {
	return paramFunc(func(logger wlog.LogEntry) {
		logger.StringValue(WrappedEntityNameKey, name)
		logger.StringValue(WrappedEntityVersionKey, version)
	})
}
