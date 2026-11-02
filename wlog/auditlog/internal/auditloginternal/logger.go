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

package auditloginternal

import (
	"time"

	"github.com/palantir/pkg/uuid"
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
)

type AuditResultType string

// audit2ConjureLogger consumes borrowed log objects synchronously. Implementations must copy
// any object or nested data they retain after the method returns.
type audit2ConjureLogger interface {
	LogAuditLogV2(*logging.AuditLogV2)
}

// audit3ConjureLogger consumes borrowed log objects synchronously. Implementations must copy
// any object or nested data they retain after the method returns.
type audit3ConjureLogger interface {
	LogAuditLogV3(*logging.AuditLogV3)
}

const (
	AuditResultSuccess      AuditResultType = "SUCCESS"
	AuditResultUnauthorized AuditResultType = "UNAUTHORIZED"
	AuditResultError        AuditResultType = "ERROR"
)

type logger struct {
	audit2Logger audit2ConjureLogger
	audit3Logger wlog.Logger
}

// newAudit2Logger returns nil when logger is nil so callers can skip dual-writing.
func newAudit2Logger(logger wlog.Logger) audit2ConjureLogger {
	if logger == nil {
		return nil
	}
	if direct, ok := logger.(audit2ConjureLogger); ok {
		return direct
	}
	return audit2EntryLogger{logger: logger}
}

type audit2EntryLogger struct {
	logger wlog.Logger
}

func (l audit2EntryLogger) LogAuditLogV2(log *logging.AuditLogV2) {
	l.logger.Log(audit2ObjectToParam(log))
}

// atLogTimeValues stores values that are computed once when first logged. Exists to support dual-writing audit v2 and
// v3 logs: when dual-writing, there are certain fields that should be computed when they are first logged, but should
// then be logged with same value when written to the other logger -- for example, "time". This struct exists so that
// such values can be computed when the logging first occurs, but subsequent invocations return the same value.
//
// Note that this structure is not safe to be used concurrently. Although it is straightforward to make it safe, it is
// on a possible hot path (audit logging) and is an internal struct that is currently only used in a non-concurrent
// manner, at this point the performance cost of making it safe does not seem to be worth the benefit.
type atLogTimeValues struct {
	timeValue           *time.Time
	logEntryIDValue     *string
	generatedLogEntryID *uuid.UUID
}

func (v *atLogTimeValues) getOrComputeTimeValue() time.Time {
	if v.timeValue == nil {
		v.timeValue = new(time.Now())
	}
	return *v.timeValue
}

func (v *atLogTimeValues) getLogEntryIDValue() *string {
	if v.logEntryIDValue == nil && v.generatedLogEntryID != nil {
		v.logEntryIDValue = new(v.generatedLogEntryID.String())
	}
	return v.logEntryIDValue
}

func (l *logger) Audit2(name string, result AuditResultType, params ...Audit2Param) {
	atLogTimeValuesProvider := &atLogTimeValues{}

	v2Log := poolAuditLogV2.Get()
	populateAuditLogV2(v2Log, atLogTimeValuesProvider, name, result, params)
	l.audit2Logger.LogAuditLogV2(v2Log)
	poolAuditLogV2.Put(v2Log)

	// dual-log if audit3Logger specified
	if l.audit3Logger != nil {
		audit3Params := make([]Audit3Param, len(params))
		for i, param := range params {
			audit3Params[i] = Audit3Param{
				Param: param.Param,
			}
		}
		v3Log := poolAuditLogV3.Get()
		populateAuditLogV3(v3Log, atLogTimeValuesProvider, name, result, audit3Params)
		emitAudit3(l.audit3Logger, v3Log, atLogTimeValuesProvider)
		poolAuditLogV3.Put(v3Log)
	}
}

func (l *logger) Audit3(name string, result AuditResultType, params ...Audit3Param) {
	atLogTimeValuesProvider := &atLogTimeValues{}

	v3Log := poolAuditLogV3.Get()
	populateAuditLogV3(v3Log, atLogTimeValuesProvider, name, result, params)
	emitAudit3(l.audit3Logger, v3Log, atLogTimeValuesProvider)
	poolAuditLogV3.Put(v3Log)

	// dual-log if audit2Logger specified
	if l.audit2Logger != nil {
		audit2Params := make([]Audit2Param, len(params))
		for i, param := range params {
			audit2Params[i] = Audit2Param{
				Param: param.Param,
			}
		}
		v2Log := poolAuditLogV2.Get()
		populateAuditLogV2(v2Log, atLogTimeValuesProvider, name, result, audit2Params)
		l.audit2Logger.LogAuditLogV2(v2Log)
		poolAuditLogV2.Put(v2Log)
	}
}

func emitAudit3(target wlog.Logger, log *logging.AuditLogV3, values *atLogTimeValues) {
	if direct, ok := target.(audit3ConjureLogger); ok && CanEmitAudit3Direct(log, values.logEntryIDValue) {
		direct.LogAuditLogV3(log)
		return
	}
	target.Log(audit3ObjectToParam(log, values.getLogEntryIDValue()))
}

func CanEmitAudit3Direct(log *logging.AuditLogV3, suppliedLogEntryID *string) bool {
	// The generated struct requires an event ID. Keep the sparse path when the caller omitted it or supplied an invalid
	// UUID so direct marshaling cannot invent the zero UUID as an audit event ID. A non-canonical log entry ID also uses
	// this path so the supplied string is preserved exactly.
	if log.EventId == (uuid.UUID{}) || log.LogEntryId == nil {
		return false
	}
	return suppliedLogEntryID == nil || *suppliedLogEntryID == log.LogEntryId.String()
}
