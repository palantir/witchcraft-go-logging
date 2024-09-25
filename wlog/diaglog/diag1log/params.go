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
	"github.com/palantir/pkg/datetime"
	"github.com/palantir/witchcraft-go-logging/wapi/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"maps"
	"time"
)

const (
	TypeValue = "diagnostic.1"
)

type Param = wlog.ConjureLogParam[logging.DiagnosticLogV1]

func Type() Param {
	return func(l *logging.DiagnosticLogV1) {
		l.Type = TypeValue
	}
}

func Time(time time.Time) Param {
	return func(l *logging.DiagnosticLogV1) {
		l.Time = datetime.DateTime(time)
	}
}

func TimeNow() Param {
	// Defer execution of time.Now() until the log is actually written
	return func(l *logging.DiagnosticLogV1) {
		l.Time = datetime.DateTime(time.Now())
	}
}

func Diagnostic(diagnostic logging.Diagnostic) Param {
	return func(l *logging.DiagnosticLogV1) {
		l.Diagnostic = diagnostic
	}
}

func GenericDiagnostic(genericDiagnostic logging.GenericDiagnostic) Param {
	return func(l *logging.DiagnosticLogV1) {
		l.Diagnostic = logging.NewDiagnosticFromGeneric(genericDiagnostic)
	}
}

func ThreadDump(threadDumpV1 logging.ThreadDumpV1) Param {
	return func(l *logging.DiagnosticLogV1) {
		l.Diagnostic = logging.NewDiagnosticFromThreadDump(threadDumpV1)
	}
}

func UnsafeParams(unsafeParams map[string]any) Param {
	return func(l *logging.DiagnosticLogV1) {
		if l.UnsafeParams == nil {
			l.UnsafeParams = maps.Clone(unsafeParams)
		} else {
			for k, v := range unsafeParams {
				l.UnsafeParams[k] = v
			}
		}
	}
}

func UnsafeParam(key string, value any) Param {
	return func(l *logging.DiagnosticLogV1) {
		if l.UnsafeParams == nil {
			l.UnsafeParams = map[string]any{key: value}
		} else {
			l.UnsafeParams[key] = value
		}
	}
}
