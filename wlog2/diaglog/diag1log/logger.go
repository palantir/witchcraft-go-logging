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
	"io"
	"os"

	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft/api/logging"
	wlog "github.com/palantir/witchcraft-go-logging/wlog2"
)

var (
	DefaultOutput = os.Stdout
)

type Logger interface {
	Diagnostic(diagnostic logging.Diagnostic, params ...Param)

	WithParams(params ...Param) Logger
}

func New(w io.Writer, params ...Param) Logger {
	return &defaultLogger{
		logger: wlog.NewDefaultLogger(w, Type(), TimeNow()).WithParams(params...),
	}
}

type defaultLogger struct {
	logger wlog.ConjureLogger[logging.DiagnosticLogV1]
}

func (l *defaultLogger) Diagnostic(diagnostic logging.Diagnostic, params ...Param) {
	_ = diagnostic.AcceptFuncs(
		func(genericDiagnostic logging.GenericDiagnostic) error {
			l.logger.Log(append([]Param{GenericDiagnostic(genericDiagnostic)}, params...)...)
			return nil
		},
		func(threadDumpV1 logging.ThreadDumpV1) error {
			l.logger.Log(append([]Param{ThreadDump(threadDumpV1)}, params...)...)
			return nil
		},
		diagnostic.ErrorOnUnknown,
	)
}

func (l *defaultLogger) WithParams(params ...Param) Logger {
	return &defaultLogger{
		logger: l.logger.WithParams(params...),
	}
}
