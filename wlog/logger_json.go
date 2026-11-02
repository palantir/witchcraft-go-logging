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

package wlog

import (
	"log"
	"time"

	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
)

type jsonLogger struct {
	output *jsonOutput
	*AtomicLogLevel
}

func (l *jsonLogger) Log(params ...Param) {
	l.logOutput(params)
}

func (l *jsonLogger) Debug(msg string, params ...Param) {
	if l.Enabled(DebugLevel) {
		l.logOutput(ParamsWithMessage(msg, params))
	}
}

func (l *jsonLogger) Info(msg string, params ...Param) {
	if l.Enabled(InfoLevel) {
		l.logOutput(ParamsWithMessage(msg, params))
	}
}

func (l *jsonLogger) Warn(msg string, params ...Param) {
	if l.Enabled(WarnLevel) {
		l.logOutput(ParamsWithMessage(msg, params))
	}
}

func (l *jsonLogger) Error(msg string, params ...Param) {
	if l.Enabled(ErrorLevel) {
		l.logOutput(ParamsWithMessage(msg, params))
	}
}

func (l *jsonLogger) logOutput(params []Param) {
	entry := NewMapLogEntry()
	entry.StringValue(TimeKey, time.Now().Format(time.RFC3339Nano))
	ApplyParams(entry, params)
	l.writeObject("json", entry.AllValues())
}

func (l *jsonLogger) writeObject(logType string, obj any) {
	if err := l.output.write(obj); err != nil {
		log.Printf("%s logger failed: %v", logType, err)
	}
}

func (l *jsonLogger) LogAuditLogV2(obj *logging.AuditLogV2) {
	l.writeObject("audit.2", obj)
}

func (l *jsonLogger) LogAuditLogV3(obj *logging.AuditLogV3) {
	l.writeObject("audit.3", obj)
}

func (l *jsonLogger) LogDiagnosticLogV1(obj *logging.DiagnosticLogV1) {
	l.writeObject("diagnostic.1", obj)
}

func (l *jsonLogger) LogEventLogV2(obj *logging.EventLogV2) {
	l.writeObject("event.2", obj)
}

func (l *jsonLogger) LogMetricLogV1(obj *logging.MetricLogV1) {
	l.writeObject("metric.1", obj)
}

func (l *jsonLogger) LogRequestLogV2(obj *logging.RequestLogV2) {
	l.writeObject("request.2", obj)
}

func (l *jsonLogger) LogServiceLogV1(obj *logging.ServiceLogV1) {
	l.writeObject("service.1", obj)
}

func (l *jsonLogger) LogTraceLogV1(obj *logging.TraceLogV1) {
	l.writeObject("trace.1", obj)
}

func (l *jsonLogger) LogWrappedLogV1(obj *logging.WrappedLogV1) {
	l.writeObject("wrapped.1", obj)
}
