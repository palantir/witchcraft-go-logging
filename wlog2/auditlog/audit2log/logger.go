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

package audit2log

import (
	"io"
	"os"

	"github.com/palantir/witchcraft-go-logging/wapi/logging"
	wlog "github.com/palantir/witchcraft-go-logging/wlog2"
)

var (
	DefaultOutput = os.Stdout
)

type AuditResultType logging.AuditResult_Value

const (
	AuditResultSuccess      = AuditResultType(logging.AuditResult_SUCCESS)
	AuditResultUnauthorized = AuditResultType(logging.AuditResult_UNAUTHORIZED)
	AuditResultError        = AuditResultType(logging.AuditResult_ERROR)
)

type Logger interface {
	Audit(name string, result AuditResultType, params ...Param)
}

func New(w io.Writer, params ...Param) Logger {
	return &wrappedLogger{
		logger: &defaultLogger{
			logger: wlog.NewDefaultLogger(w, Type(), TimeNow()),
		},
		params: params,
	}
}

type defaultLogger struct {
	logger wlog.ConjureLogger[logging.AuditLogV2]
}

func (l *defaultLogger) Audit(name string, result AuditResultType, params ...Param) {
	l.logger.Log(append([]Param{Name(name), Result(result)}, params...)...)
}

func WithParams(logger Logger, params ...Param) Logger {
	switch logger := logger.(type) {
	case *defaultLogger:
		return &wrappedLogger{
			logger: logger,
			params: params,
		}
	case *wrappedLogger:
		return &wrappedLogger{
			logger: logger.logger,
			params: append(append([]Param{}, logger.params...), params...),
		}
	default:
		return &wrappedLogger{
			logger: logger,
			params: params,
		}
	}
}

type wrappedLogger struct {
	logger Logger
	params []Param
}

func (l *wrappedLogger) Audit(name string, result AuditResultType, params ...Param) {
	l.logger.Audit(name, result, append(append([]Param{}, l.params...), params...)...)
}
