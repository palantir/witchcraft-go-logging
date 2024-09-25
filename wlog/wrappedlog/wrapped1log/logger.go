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

	"github.com/palantir/witchcraft-go-logging/wapi/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit2log"
	"github.com/palantir/witchcraft-go-logging/wlog/diaglog/diag1log"
	"github.com/palantir/witchcraft-go-logging/wlog/evtlog/evt2log"
	"github.com/palantir/witchcraft-go-logging/wlog/metriclog/metric1log"
	"github.com/palantir/witchcraft-go-logging/wlog/reqlog/req2log"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/palantir/witchcraft-go-logging/wlog/trclog/trc1log"
)

type Logger interface {
	Audit() audit2log.Logger
	Diagnostic() diag1log.Logger
	Event() evt2log.Logger
	Metric() metric1log.Logger
	Request(params ...req2log.LoggerCreatorParam) req2log.Logger
	Service(params ...svc1log.Param) svc1log.Logger
	Trace() trc1log.Logger
}

func New(w io.Writer, level wlog.LogLevel, name, version string) Logger {
	delegate := wlog.NewDefaultLogger(w, Type(), EntityName(name), EntityVersion(version))
	return newDelegateLogger(delegate, level)
}

func NewWithPrinter(printer wlog.LogPrinter[logging.WrappedLogV1], level wlog.LogLevel, name, version string) Logger {
	delegate := wlog.NewDefaultLoggerWithPrinter(printer, Type(), EntityName(name), EntityVersion(version))
	return newDelegateLogger(delegate, level)
}

func NewFromProvider(w io.Writer, level wlog.LogLevel, creator wlog.LoggerProvider, name, version string) Logger {
	delegate := creator.NewLeveledLogger(w, level)
	// The second return value is ignored because 'level: nil' is a valid state handled in the implementation.
	levelChecker, _ := delegate.(wlog.LevelChecker)
	return &defaultLogger{
		name:        name,
		version:     version,
		creator:     creator.NewLogger,
		writer:      w,
		logger:      creator.NewLogger(w),
		levellogger: delegate,
		level:       levelChecker,
	}
}

type delegateLogger struct {
	Audit2      audit2log.Logger
	Diagnostic1 diag1log.Logger
	Event2      evt2log.Logger
	Metric1     metric1log.Logger
	Request2    req2log.Logger
	Service1    svc1log.Logger
	Trace1      trc1log.Logger
}

func newDelegateLogger(delegate wlog.Logger[logging.WrappedLogV1], level wlog.LogLevel) *delegateLogger {
	return &delegateLogger{
		Audit2:      audit2log.NewWithPrinter(wrapPrinter(delegate, logging.NewWrappedLogV1PayloadFromAuditLogV2)),
		Diagnostic1: diag1log.NewWithPrinter(wrapPrinter(delegate, logging.NewWrappedLogV1PayloadFromDiagnosticLogV1)),
		Event2:      evt2log.NewWithPrinter(wrapPrinter(delegate, logging.NewWrappedLogV1PayloadFromEventLogV2)),
		Metric1:     metric1log.NewWithPrinter(wrapPrinter(delegate, logging.NewWrappedLogV1PayloadFromMetricLogV1)),
		Request2:    req2log.NewWithPrinter(wrapPrinter(delegate, logging.NewWrappedLogV1PayloadFromRequestLogV2)),
		Service1:    svc1log.NewWithPrinter(wrapPrinter(delegate, logging.NewWrappedLogV1PayloadFromServiceLogV1), level),
		Trace1:      trc1log.NewWithPrinter(wrapPrinter(delegate, logging.NewWrappedLogV1PayloadFromTraceLogV1)),
	}
}

func (l *delegateLogger) Audit() audit2log.Logger {
	return l.Audit2
}

func (l *delegateLogger) Diagnostic() diag1log.Logger {
	return l.Diagnostic1
}

func (l *delegateLogger) Event() evt2log.Logger {
	return l.Event2
}

func (l *delegateLogger) Metric() metric1log.Logger {
	return l.Metric1
}

func (l *delegateLogger) Request(params ...req2log.LoggerCreatorParam) req2log.Logger {
	return l.Request2
}

func (l *delegateLogger) Service(params ...svc1log.Param) svc1log.Logger {
	return svc1log.WithParams(l.Service1, params...)
}

func (l *delegateLogger) Trace() trc1log.Logger {
	return l.Trace1
}
