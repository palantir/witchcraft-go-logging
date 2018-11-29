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
	"fmt"

	"github.com/palantir/witchcraft-go-logging/conjure/sls/spec/logging"
	"github.com/rs/zerolog"
)

func marshalLoggingDiagnostic(evt *zerolog.Event, key string, val interface{}) *zerolog.Event {
	return evt.Object(key, logObjectMarshalerFn(func(e *zerolog.Event) {
		diagnostic := val.(logging.Diagnostic)
		_ = diagnostic.Accept(&marshalVisitor{
			evt: evt,
		})
	}))
}

type marshalVisitor struct {
	evt *zerolog.Event
}

func (m *marshalVisitor) VisitGeneric(v logging.GenericDiagnostic) error {
	m.evt.Str("type", "generic")
	m.evt.Object("generic", logObjectMarshalerFn(func(e *zerolog.Event) {
		e.Str("diagnosticType", v.DiagnosticType)
		e.Interface("value", v.Value)

	}))
	return nil
}

func (m *marshalVisitor) VisitThreadDump(v logging.ThreadDumpV1) error {
	m.evt.Str("type", "threadDump")
	m.evt.Object("threadDump", logObjectMarshalerFn(func(e *zerolog.Event) {
		e.Array("threads", logArrayMarshalerFn(func(a *zerolog.Array) {
			for _, currThread := range v.Threads {
				a.Object(threadInfoV1Encoder(currThread))
			}
		}))
	}))
	return nil
}

func threadInfoV1Encoder(threadInfo logging.ThreadInfoV1) zerolog.LogObjectMarshaler {
	return logObjectMarshalerFn(func(e *zerolog.Event) {
		if threadInfo.Id != nil {
			e.Int64("id", int64(*threadInfo.Id))
		}
		encodeNonEmptyString(e, "name", threadInfo.Name)
		if len(threadInfo.StackTrace) > 0 {
			e.Array("stackTrace", logArrayMarshalerFn(func(a *zerolog.Array) {
				for _, stackFrame := range threadInfo.StackTrace {
					a.Object(stackFrameV1Encoder(stackFrame))
				}
			}))
		}
		if len(threadInfo.Params) > 0 {
			e.Interface("params", threadInfo.Params)
		}
	})
}

func stackFrameV1Encoder(stackFrame logging.StackFrameV1) zerolog.LogObjectMarshaler {
	return logObjectMarshalerFn(func(e *zerolog.Event) {
		encodeNonEmptyString(e, "address", stackFrame.Address)
		encodeNonEmptyString(e, "procedure", stackFrame.Procedure)
		encodeNonEmptyString(e, "file", stackFrame.File)
		if stackFrame.Line != nil {
			e.Int("line", *stackFrame.Line)
		}
		if len(stackFrame.Params) > 0 {
			e.Interface("params", stackFrame.Params)
		}
	})
}

func encodeNonEmptyString(e *zerolog.Event, key string, val *string) {
	if val == nil || len(*val) == 0 {
		return
	}
	e.Str(key, *val)
}

func (m *marshalVisitor) VisitUnknown(typeName string) error {
	return fmt.Errorf("unknown diagnostic type: %s", typeName)
}
