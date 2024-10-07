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
	"io"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit2log"
	"github.com/palantir/witchcraft-go-logging/wlog/diaglog/diag1log"
	"github.com/palantir/witchcraft-go-logging/wlog/evtlog/evt2log"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
	"github.com/palantir/witchcraft-go-logging/wlog/metriclog/metric1log"
	"github.com/palantir/witchcraft-go-logging/wlog/reqlog/req2log"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/palantir/witchcraft-go-logging/wlog/trclog/trc1log"
	"github.com/palantir/witchcraft-go-logging/wtypes"
)

var objectPool = wloginternal.NewPool((*wtypes.WrappedLogV1).Reset)

type defaultLogger struct {
	delegate wlog.Logger[wtypes.WrappedLogV1]
	level    wlog.LogLevel
	params   []Param
}

func (l *defaultLogger) Audit() audit2log.Logger {
	return audit2log.NewFromCreator(nil, wrapPrinter(l.delegate, wtypes.NewWrappedLogV1PayloadFromAuditLogV2, l.params))
}

func (l *defaultLogger) Diagnostic() diag1log.Logger {
	return diag1log.NewFromCreator(nil, wrapPrinter(l.delegate, wtypes.NewWrappedLogV1PayloadFromDiagnosticLogV1, l.params))
}

func (l *defaultLogger) Event() evt2log.Logger {
	return evt2log.NewFromCreator(nil, wrapPrinter(l.delegate, wtypes.NewWrappedLogV1PayloadFromEventLogV2, l.params))
}

func (l *defaultLogger) Metric() metric1log.Logger {
	return metric1log.NewFromCreator(nil, wrapPrinter(l.delegate, wtypes.NewWrappedLogV1PayloadFromMetricLogV1, l.params))
}

func (l *defaultLogger) Request(params ...req2log.LoggerCreatorParam) req2log.Logger {
	return req2log.NewFromCreator(nil, wrapPrinter(l.delegate, wtypes.NewWrappedLogV1PayloadFromRequestLogV2, l.params), params...)
}

func (l *defaultLogger) Service(params ...svc1log.Param) svc1log.Logger {
	return svc1log.NewFromCreator(nil, l.level, wrapPrinter(l.delegate, wtypes.NewWrappedLogV1PayloadFromServiceLogV1, l.params), params...)
}

func (l *defaultLogger) Trace() trc1log.Logger {
	return trc1log.NewFromCreator(nil, wrapPrinter(l.delegate, wtypes.NewWrappedLogV1PayloadFromTraceLogV1, l.params))
}

// wrappedPrinter implements Printer for logs included in the wrapped.1 payload field.
// When an underlying log object is constructed and passed to the Print method,
// the delegate WrappedLogV1 logger is called with a new WrappedLogV1Payload object.
type wrappedPrinter[T wtypes.LogTypes] struct {
	logger     wlog.Logger[wtypes.WrappedLogV1]
	newPayload func(payload T) wtypes.WrappedLogV1Payload
	params     []Param
}

func wrapPrinter[T wtypes.LogTypes](
	logger wlog.Logger[wtypes.WrappedLogV1],
	newPayload func(payload T) wtypes.WrappedLogV1Payload,
	params []Param,
) wlog.LoggerCreator[T] {
	printer := wrappedPrinter[T]{logger: logger, newPayload: newPayload, params: params}
	return func(io.Writer) wlog.Logger[T] { return wlog.NewDefaultLoggerWithPrinter[T](printer) }
}

func (l wrappedPrinter[T]) Print(obj wtypes.LogType) error {
	wloginternal.LogObject(l.logger.Log, objectPool, defaultParam(l.newPayload(obj.(T))), l.params...)
	return nil
}
