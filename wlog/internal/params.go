// Copyright (c) 2024 Palantir Technologies. All rights reserved.
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

package wloginternal

import (
	"maps"

	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
)

type LogType interface {
	logging.AuditLogV2 | logging.AuditLogV3 | logging.DiagnosticLogV1 | logging.EventLogV2 | logging.MetricLogV1 | logging.RequestLogV2 | logging.ServiceLogV1 | logging.TraceLogV1 | logging.WrappedLogV1
}

// Param is an interface that modifies a log object, generally used to set struct fields.
//
// Each log type will export its own param constructors that implement this interface.
// Parameters are applied in order to a struct representing a single log line.
//
// It is implemented by ParamFunc[T].
type Param[T LogType] interface {
	apply(*T)
}

// ParamFunc is a function that implements Param and modifies a log struct.
type ParamFunc[T LogType] func(log *T)

// apply makes ParamFunc[T] satisfy Param[T]. Invoked via interface dispatch in ApplyParam/ApplyParams.
//
//nolint:unused // false positive: called through the Param interface, not statically visible to the linter.
func (f ParamFunc[T]) apply(log *T) {
	f(log)
}

// ApplyParam executes the private apply function on the provided log struct.
func ApplyParam[T LogType](log *T, p Param[T]) {
	if p != nil {
		p.apply(log)
	}
}

// ApplyParams executes the private apply functions on the provided log struct.
func ApplyParams[T LogType](log *T, params []Param[T]) {
	for _, p := range params {
		if p != nil {
			p.apply(log)
		}
	}
}

// OmitEmpty clears an optional string that params set to empty. Params keep the empty value so that, applied to a
// wlog.LogEntry, it removes an earlier value like an empty OptionalStringValue.
func OmitEmpty[T ~string](value **T) {
	if *value != nil && **value == "" {
		*value = nil
	}
}

// SetMapParam sets the key:value pair on the target map, taking a pooled map if it is nil.
func SetMapParam[V any](m *map[string]V, key string, value V) {
	if *m == nil {
		*m = getMap[V]()
	}
	(*m)[key] = value
}

// SetMapParams sets the key:value pairs on the target map, taking a pooled map if it is nil.
func SetMapParams[V any](m *map[string]V, params map[string]V) {
	if len(params) == 0 {
		return
	}
	if *m == nil {
		*m = getMap[V]()
	}
	maps.Copy(*m, params)
}
