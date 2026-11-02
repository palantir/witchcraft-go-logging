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
)

const (
	TypeValue = "diagnostic.1"

	DiagnosticKey = "diagnostic"
)

type Param = wloginternal.Param[logging.DiagnosticLogV1]

type paramFunc = wloginternal.ParamFunc[logging.DiagnosticLogV1]

func ApplyParam(p Param, entry wlog.LogEntry) {
	if p == nil {
		return
	}
	log := new(logging.DiagnosticLogV1)
	wloginternal.ApplyParam(log, p)
	applyObjectParams(log, entry)
}

// populateLogObject pre-populates the type, time, and diagnostic fields and applies user params in-place.
func populateLogObject(log *logging.DiagnosticLogV1, diagnostic logging.Diagnostic, params []Param) {
	log.Type = TypeValue
	log.Time = datetime.DateTime(time.Now())
	log.Diagnostic = diagnostic
	wloginternal.ApplyParams(log, params)
}

func objectToParam(log *logging.DiagnosticLogV1) wlog.Param {
	return wlog.NewParam(func(entry wlog.LogEntry) {
		entry.StringValue(wlog.TypeKey, log.Type)
		entry.StringValue(wlog.TimeKey, log.Time.String())
		if dp := diagnosticFromObject(log.Diagnostic); dp != nil {
			wlog.ApplyParams(entry, []wlog.Param{dp})
		}
		applyObjectParams(log, entry)
	})
}

func applyObjectParams(log *logging.DiagnosticLogV1, entry wlog.LogEntry) {
	if len(log.UnsafeParams) > 0 {
		entry.AnyMapValue(wlog.UnsafeParamsKey, log.UnsafeParams)
	}
}

func UnsafeParam(key string, value any) Param {
	return paramFunc(func(log *logging.DiagnosticLogV1) {
		wloginternal.SetMapParam(&log.UnsafeParams, key, value)
	})
}

func UnsafeParams(unsafe map[string]any) Param {
	return paramFunc(func(log *logging.DiagnosticLogV1) {
		wloginternal.SetMapParams(&log.UnsafeParams, unsafe)
	})
}
