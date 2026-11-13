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

package wrapped1log

import (
	"github.com/palantir/witchcraft-go-logging/wlog"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logentry"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logtypes"
)

var poolWrappedLogV1 = wloginternal.NewPool(func(log *logtypes.WrappedLogV1) { *log = logtypes.WrappedLogV1{} })

// sink wraps each log object with the entity name and version. Family loggers send it only log objects, so the
// wlog.Logger and wlog.LeveledLogger methods exist only to fit the creator types and drop their input.
type sink struct {
	name    string
	version string
	logger  logentry.Logger

	// leveled and level are set only for service loggers.
	leveled wlog.LeveledLogger
	level   wlog.LevelChecker
}

func newSink(name, version string, logger wlog.Logger) *sink {
	return &sink{name: name, version: version, logger: logentry.NewLogger(logger)}
}

func newLeveledSink(name, version string, logger wlog.LeveledLogger) *sink {
	level, _ := logger.(wlog.LevelChecker)
	return &sink{name: name, version: version, logger: logentry.NewLeveledLogger(logger), leveled: logger, level: level}
}

func (s *sink) LogObject(obj logtypes.Log) {
	payload, ok := payloadOf(obj)
	if !ok {
		return
	}
	log := poolWrappedLogV1.Get()
	log.Type = TypeValue
	log.Payload = payload
	log.EntityName = s.name
	log.EntityVersion = s.version
	s.logger.LogObject(log)
	poolWrappedLogV1.Put(log)
}

func (s *sink) Enabled(level wlog.LogLevel) bool {
	return s.level == nil || s.level.Enabled(level)
}

func (s *sink) SetLevel(level wlog.LogLevel) {
	s.leveled.SetLevel(level)
}

func (*sink) Log(...wlog.Param)           {}
func (*sink) Debug(string, ...wlog.Param) {}
func (*sink) Info(string, ...wlog.Param)  {}
func (*sink) Warn(string, ...wlog.Param)  {}
func (*sink) Error(string, ...wlog.Param) {}

func payloadOf(obj logtypes.Log) (logtypes.WrappedLogV1Payload, bool) {
	switch obj := obj.(type) {
	case *logtypes.ServiceLogV1:
		return logtypes.WrappedLogV1Payload{Type: PayloadServiceLogV1, ServiceLogV1: obj}, true
	case *logtypes.RequestLogV2:
		return logtypes.WrappedLogV1Payload{Type: PayloadRequestLogV2, RequestLogV2: obj}, true
	case *logtypes.TraceLogV1:
		return logtypes.WrappedLogV1Payload{Type: PayloadTraceLogV1, TraceLogV1: obj}, true
	case *logtypes.EventLogV2:
		return logtypes.WrappedLogV1Payload{Type: PayloadEventLogV2, EventLogV2: obj}, true
	case *logtypes.MetricLogV1:
		return logtypes.WrappedLogV1Payload{Type: PayloadMetricLogV1, MetricLogV1: obj}, true
	case *logtypes.AuditLogV2:
		return logtypes.WrappedLogV1Payload{Type: PayloadAuditLogV2, AuditLogV2: obj}, true
	case *logtypes.AuditLogV3:
		return logtypes.WrappedLogV1Payload{Type: PayloadAuditLogV3, AuditLogV3: obj}, true
	case *logtypes.DiagnosticLogV1:
		return logtypes.WrappedLogV1Payload{Type: PayloadDiagnosticLogV1, DiagnosticLogV1: obj}, true
	}
	return logtypes.WrappedLogV1Payload{}, false
}
