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

package marshalers

import (
	"reflect"

	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft/api/logging"
	"github.com/rs/zerolog"
)

type encoderFunc func(evt *zerolog.Event, key string, val interface{}) *zerolog.Event

type logObjectMarshalerFn func(e *zerolog.Event)

func (f logObjectMarshalerFn) MarshalZerologObject(e *zerolog.Event) {
	f(e)
}

type logArrayMarshalerFn func(a *zerolog.Array)

func (f logArrayMarshalerFn) MarshalZerologArray(a *zerolog.Array) {
	f(a)
}

var encoders = map[reflect.Type]encoderFunc{
	reflect.TypeOf(logging.Diagnostic{}): marshalLoggingDiagnostic,
}

func EncodeType(evt *zerolog.Event, typ reflect.Type, key string, val interface{}) bool {
	fn, ok := encoders[typ]
	if !ok {
		return false
	}
	fn(evt, key, val)
	return true
}