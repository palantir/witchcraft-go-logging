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
	"bytes"
	"fmt"
	"strconv"

	"github.com/palantir/witchcraft-go-logging/conjure/sls/spec/logging"
)

func marshalLoggingDiagnostic(key string, val interface{}) string {
	builder := &bytes.Buffer{}
	_, _ = builder.WriteString("{")
	diagnostic := val.(logging.Diagnostic)
	_ = diagnostic.Accept(&marshalVisitor{
		builder: builder,
	})
	_, _ = builder.WriteString("}")
	return builder.String()
}

type marshalVisitor struct {
	builder *bytes.Buffer
}

func (m *marshalVisitor) VisitGeneric(v logging.GenericDiagnostic) error {
	_, _ = m.builder.WriteString("type: generic")
	_, _ = m.builder.WriteString(separator)

	_, _ = m.builder.WriteString("generic: {")
	_, _ = m.builder.WriteString("diagnosticType: " + v.DiagnosticType)
	_, _ = m.builder.WriteString(separator)
	_, _ = m.builder.WriteString(fmt.Sprintf("value: %v", v.Value))
	_, _ = m.builder.WriteString("}")
	return nil
}

func (m *marshalVisitor) VisitThreadDump(v logging.ThreadDumpV1) error {
	_, _ = m.builder.WriteString("type: threadDump")
	_, _ = m.builder.WriteString(separator)

	_, _ = m.builder.WriteString("threadDump: {")
	_, _ = m.builder.WriteString("threads: [")
	for i, thread := range v.Threads {
		encodeThreadInfoV1(m.builder, thread)
		if i != len(v.Threads)-1 {
			_, _ = m.builder.WriteString(separator)
		}
	}
	_, _ = m.builder.WriteString("]")
	_, _ = m.builder.WriteString("}")
	return nil
}

func encodeThreadInfoV1(builder *bytes.Buffer, threadInfo logging.ThreadInfoV1) {
	needSeparator := false
	if threadInfo.Id != nil {
		_, _ = builder.WriteString("id: " + fmt.Sprint(int64(*threadInfo.Id)))
		needSeparator = true
	}
	encodeNonEmptyStr(builder, "name", threadInfo.Name, &needSeparator)
	if len(threadInfo.StackTrace) > 0 {
		if needSeparator {
			_, _ = builder.WriteString(separator)
		}
		_, _ = builder.WriteString("stackTrace: [")
		for i, stackFrame := range threadInfo.StackTrace {
			encodeStackFrameV1(builder, stackFrame)
			if i != len(threadInfo.StackTrace)-1 {
				_, _ = builder.WriteString(separator)
			}
		}
		_, _ = builder.WriteString("]")
		needSeparator = true
	}
	if len(threadInfo.Params) > 0 {
		if needSeparator {
			_, _ = builder.WriteString(separator)
		}
		_, _ = builder.WriteString(fmt.Sprintf("params: %v", threadInfo.Params))
		needSeparator = true
	}
}

func encodeStackFrameV1(builder *bytes.Buffer, stackFrame logging.StackFrameV1) {
	needSeparator := false
	encodeNonEmptyStr(builder, "address", stackFrame.Address, &needSeparator)
	encodeNonEmptyStr(builder, "procedure", stackFrame.Procedure, &needSeparator)
	encodeNonEmptyStr(builder, "file", stackFrame.File, &needSeparator)
	if stackFrame.Line != nil {
		if needSeparator {
			_, _ = builder.WriteString(separator)
		}
		_, _ = builder.WriteString("line: " + strconv.Itoa(*stackFrame.Line))
		needSeparator = true
	}
	if len(stackFrame.Params) > 0 {
		if needSeparator {
			_, _ = builder.WriteString(separator)
		}
		_, _ = builder.WriteString(fmt.Sprintf("params: %v", stackFrame.Params))
		needSeparator = true
	}
}

func encodeNonEmptyStr(builder *bytes.Buffer, key string, val *string, needSeparator *bool) {
	if val == nil || len(*val) == 0 {
		return
	}
	if *needSeparator {
		_, _ = builder.WriteString(separator)
	}
	_, _ = builder.WriteString(key + ": " + *val)
	*needSeparator = true
}

func (m *marshalVisitor) VisitUnknown(typeName string) error {
	return fmt.Errorf("unknown diagnostic type: %s", typeName)
}
