// Copyright (c) 2018 Palantir Technologies. All rights reserved.
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

package diag1log

import (
	"time"

	"github.com/palantir/pkg/datetime"
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logentry"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logtypes"
)

const (
	TypeValue = "diagnostic.1"

	DiagnosticKey = "diagnostic"
)

type Param = wloginternal.Param[logtypes.DiagnosticLogV1]

type paramFunc = wloginternal.ParamFunc[logtypes.DiagnosticLogV1]

func ApplyParam(p Param, entry wlog.LogEntry) {
	if p == nil {
		return
	}
	log := new(logtypes.DiagnosticLogV1)
	wloginternal.ApplyParam(log, p)
	logentry.DiagnosticLogV1Params(entry, log)
}

// populateLogObject pre-populates the type, time, and diagnostic fields and applies user params in-place.
func populateLogObject(log *logtypes.DiagnosticLogV1, diagnostic logging.Diagnostic, params []Param) {
	log.Type = TypeValue
	log.Time = datetime.DateTime(time.Now())
	_ = diagnostic.Accept(diagnosticVisitor{&log.Diagnostic})
	wloginternal.ApplyParams(log, params)
}

// diagnosticVisitor copies the variant of a generated union into the log object's union.
type diagnosticVisitor struct {
	out *logtypes.DiagnosticLogV1Diagnostic
}

func (d diagnosticVisitor) VisitGeneric(v logging.GenericDiagnostic) error {
	*d.out = logtypes.DiagnosticLogV1Diagnostic{Type: "generic", Generic: &v}
	return nil
}

func (d diagnosticVisitor) VisitThreadDump(v logging.ThreadDumpV1) error {
	*d.out = logtypes.DiagnosticLogV1Diagnostic{Type: "threadDump", ThreadDump: &v}
	return nil
}

func (d diagnosticVisitor) VisitUnknown(string) error {
	return nil
}

func UnsafeParam(key string, value any) Param {
	return paramFunc(func(log *logtypes.DiagnosticLogV1) {
		wloginternal.SetMapParam(&log.UnsafeParams, key, value)
	})
}

func UnsafeParams(unsafe map[string]any) Param {
	return paramFunc(func(log *logtypes.DiagnosticLogV1) {
		wloginternal.SetMapParams(&log.UnsafeParams, unsafe)
	})
}
