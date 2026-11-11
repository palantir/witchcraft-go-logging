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

package logentry_test

import (
	"testing"

	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logentry"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logtypes"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/objecttest"
)

func TestEntryParity(t *testing.T) {
	t.Run("service.1", func(t *testing.T) {
		requireParity(t, &logtypes.ServiceLogV1{Level: logtypes.LogLevelInfo})
	})
	t.Run("event.2", func(t *testing.T) { requireParity(t, &logtypes.EventLogV2{}) })
	t.Run("metric.1", func(t *testing.T) { requireParity(t, &logtypes.MetricLogV1{}) })
	t.Run("trace.1", func(t *testing.T) { requireParity(t, &logtypes.TraceLogV1{}) })
	t.Run("request.2", func(t *testing.T) { requireParity(t, &logtypes.RequestLogV2{}) })
	t.Run("audit.2", func(t *testing.T) { requireParity(t, &logtypes.AuditLogV2{}) })
	t.Run("audit.3", func(t *testing.T) { requireParity(t, &logtypes.AuditLogV3{}) })

	var generic logging.GenericDiagnostic
	objecttest.Fill(t, &generic)
	var threadDump logging.ThreadDumpV1
	objecttest.Fill(t, &threadDump)
	for name, diagnostic := range map[string]logtypes.DiagnosticLogV1Diagnostic{
		"generic":    {Type: "generic", Generic: &generic},
		"threadDump": {Type: "threadDump", ThreadDump: &threadDump},
	} {
		t.Run("diagnostic.1/"+name, func(t *testing.T) {
			requireParity(t, &logtypes.DiagnosticLogV1{}, diagnostic)
		})
	}
}

// requireParity checks that a filled log and empty encode the same through the JSON provider and through entries.
func requireParity[T any, P interface {
	*T
	logtypes.Log
}](t *testing.T, empty P, overrides ...any) {
	filled := P(new(T))
	objecttest.Fill(t, filled, overrides...)
	objecttest.RequireParity(t, (*T)(filled), (*T)(empty),
		func(l wlog.LeveledLogger, log *T) { l.(logentry.Logger).LogObject(P(log)) },
		func(l wlog.LeveledLogger, log *T) { entryLogger(l, P(log)).LogObject(P(log)) },
	)
}

func entryLogger(l wlog.LeveledLogger, log logtypes.Log) logentry.Logger {
	if _, ok := log.(*logtypes.ServiceLogV1); ok {
		return logentry.LeveledEntryLogger{Logger: l}
	}
	return logentry.EntryLogger{Logger: l.(wlog.Logger)}
}
