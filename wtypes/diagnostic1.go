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

package wtypes

import (
	"github.com/palantir/pkg/datetime"
	"github.com/palantir/pkg/safelong"
)

// diagnostic.1

type DiagnosticLogV1 struct {
	// "diagnostic.1"
	Type string `json:"type"`
	// RFC3339Nano timestamp when the log event was emitted
	Time datetime.DateTime `json:"time"`
	// The diagnostic being logged.
	Diagnostic Diagnostic `json:"diagnostic"`
	// Unredacted parameters
	UnsafeParams map[string]any `json:"unsafeParams,omitempty"`
}

func (log *DiagnosticLogV1) Reset() {
	log.Type = ""
	log.Time = datetime.DateTime{}
	log.Diagnostic = Diagnostic{}
	clear(log.UnsafeParams)
}

type Diagnostic struct {
	Type       string             `json:"type"`
	Generic    *GenericDiagnostic `json:"generic,omitempty"`
	ThreadDump *ThreadDumpV1      `json:"threadDump,omitempty"`
}

func NewDiagnosticFromGeneric(v GenericDiagnostic) Diagnostic {
	return Diagnostic{Type: "generic", Generic: &v}
}

func NewDiagnosticFromThreadDump(v ThreadDumpV1) Diagnostic {
	return Diagnostic{Type: "threadDump", ThreadDump: &v}
}

type GenericDiagnostic struct {
	// An identifier for the type of diagnostic represented.
	DiagnosticType string `json:"diagnosticType"`
	// Observations, measurements and context associated with the diagnostic.
	Value any `json:"value"`
}

type ThreadDumpV1 struct {
	// Information about each of the threads in the thread dump. "Thread" may refer to a userland thread such as a goroutine, or an OS-level thread.
	Threads []ThreadInfoV1 `json:"threads,omitempty"`
}

type ThreadInfoV1 struct {
	// The ID of the thread.
	Id *safelong.SafeLong `json:"id,omitempty"`
	// The name of the thread. Note that thread names may include unsafe information such as the path of the HTTP request being processed. It must be safely redacted.
	Name *string `json:"name,omitempty"`
	// A list of stack frames for the thread, ordered with the current frame first.
	StackTrace []StackFrameV1 `json:"stackTrace,omitempty"`
	// Other thread-level information.
	Params map[string]any `json:"params,omitempty"`
}

type StackFrameV1 struct {
	// The address of the execution point of this stack frame. This is a string because a safelong can't represent the full 64 bit address space.
	Address *string `json:"address,omitempty"`
	// The identifier of the procedure containing the execution point of this stack frame. This is a fully qualified method name in Java and a demangled symbol name in native code, for example. Note that procedure names may include unsafe information if a service is, for exmaple, running user-defined code. It must be safely redacted.
	Procedure *string `json:"procedure,omitempty"`
	// The name of the file containing the source location of the execution point of this stack frame. Note that file names may include unsafe information if a service is, for example, running user-defined code. It must be safely redacted.
	File *string `json:"file,omitempty"`
	// The line number of the source location of the execution point of this stack frame.
	Line *int `json:"line,omitempty"`
	// Other frame-level information.
	Params map[string]any `json:"params,omitempty"`
}
