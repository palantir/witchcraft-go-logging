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
	"github.com/palantir/witchcraft-go-logging/wapi/logging"
)

// Param is an interface that modifies a log object, generally used to set struct fields.
// It is implemented by ParamFunc[T].
type Param[T logging.LogTypes] interface {
	apply(*T)
}

// ParamFunc is a function that implements Param and modifies a Conjure log object.
type ParamFunc[T logging.LogTypes] func(log *T)

func (f ParamFunc[T]) apply(log *T) {
	f(log)
}

// ApplyParams executes the private apply functions on the provided log object.
func ApplyParams[T logging.LogTypes](log *T, params ...Param[T]) {
	for _, p := range params {
		if p != nil {
			p.apply(log)
		}
	}
}

// SetMapParam initializes the map, if nil, and sets the key:value pair on the target map.
func SetMapParam[V any](m *map[string]V, key string, value V) {
	if *m == nil {
		*m = make(map[string]V)
	}
	(*m)[key] = value
}

// SetMapParams initializes the map, if nil, and sets the key:value pairs on the target map.
func SetMapParams[V any](m *map[string]V, params map[string]V) {
	if *m == nil {
		*m = make(map[string]V)
	}
	for key, value := range params {
		(*m)[key] = value
	}
}
