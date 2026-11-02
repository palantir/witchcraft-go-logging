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

package wlog

import (
	"io"
)

// NewJSONMarshalLoggerProvider encodes typed log objects directly with JSON v2.
// Generic Logger and LeveledLogger calls use a MapLogEntry.
func NewJSONMarshalLoggerProvider() LoggerProvider {
	return &jsonMarshalLoggerProvider{}
}

type jsonMarshalLoggerProvider struct{}

func (*jsonMarshalLoggerProvider) NewLogger(w io.Writer) Logger {
	return newJSONLogger(w, InfoLevel)
}

func (*jsonMarshalLoggerProvider) NewLeveledLogger(w io.Writer, level LogLevel) LeveledLogger {
	return newJSONLogger(w, level)
}

func newJSONLogger(w io.Writer, level LogLevel) *jsonLogger {
	return &jsonLogger{
		output:         newJSONOutput(w),
		AtomicLogLevel: NewAtomicLogLevel(level),
	}
}
