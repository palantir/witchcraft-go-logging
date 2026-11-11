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
	"io"
	"time"

	"github.com/palantir/pkg/uuid"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logentry"
)

type AuditResultType string

const (
	AuditResultSuccess      AuditResultType = "SUCCESS"
	AuditResultUnauthorized AuditResultType = "UNAUTHORIZED"
	AuditResultError        AuditResultType = "ERROR"
)

// Logger writes audit.2 and audit.3 logs. Audit2 and Audit3 write their own version, then the other version if its
// logger is set.
type Logger struct {
	audit2Logger logentry.Logger
	audit3Logger logentry.Logger
}

// NewLogger returns a Logger that writes to the provided loggers. Either logger may be nil.
func NewLogger(audit2Logger, audit3Logger wlog.Logger) *Logger {
	return &Logger{audit2Logger: newObjectLogger(audit2Logger), audit3Logger: newObjectLogger(audit3Logger)}
}

// CreateIfSet returns a logger for w, or nil if w is nil.
func CreateIfSet(w io.Writer, creator wlog.LoggerCreator) wlog.Logger {
	if w == nil {
		return nil
	}
	return creator(w)
}

func newObjectLogger(logger wlog.Logger) logentry.Logger {
	if logger == nil {
		return nil
	}
	return logentry.NewLogger(logger)
}

// Audit2 writes an audit.2 log, and an audit.3 log if the audit.3 logger is set. get returns the AuditParam of a
// public param.
func Audit2[P any](l *Logger, name string, result AuditResultType, params []P, get func(P) AuditParam) {
	var values atLogTimeValues
	writeAuditLogV2(l.audit2Logger, &values, name, result, params, get)
	if l.audit3Logger != nil {
		writeAuditLogV3(l.audit3Logger, &values, name, result, params, get)
	}
}

// Audit3 writes an audit.3 log, and an audit.2 log if the audit.2 logger is set. get returns the AuditParam of a
// public param.
func Audit3[P any](l *Logger, name string, result AuditResultType, params []P, get func(P) AuditParam) {
	var values atLogTimeValues
	writeAuditLogV3(l.audit3Logger, &values, name, result, params, get)
	if l.audit2Logger != nil {
		writeAuditLogV2(l.audit2Logger, &values, name, result, params, get)
	}
}

func writeAuditLogV2[P any](logger logentry.Logger, values *atLogTimeValues, name string, result AuditResultType, params []P, get func(P) AuditParam) {
	log := poolAuditLogV2.Get()
	populateAuditLogV2(log, values, name, result, params, get)
	logger.LogObject(log)
	poolAuditLogV2.Put(log)
}

func writeAuditLogV3[P any](logger logentry.Logger, values *atLogTimeValues, name string, result AuditResultType, params []P, get func(P) AuditParam) {
	log := poolAuditLogV3.Get()
	populateAuditLogV3(log, values, name, result, params, get)
	logger.LogObject(log)
	poolAuditLogV3.Put(log)
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
	timeValue  time.Time
	logEntryID string
}

func (v *atLogTimeValues) getOrComputeTimeValue() time.Time {
	if v.timeValue.IsZero() {
		v.timeValue = time.Now()
	}
	return v.timeValue
}

func (v *atLogTimeValues) getOrGenerateLogEntryID() string {
	if v.logEntryID == "" {
		v.logEntryID = uuid.NewUUID().String()
	}
	return v.logEntryID
}
