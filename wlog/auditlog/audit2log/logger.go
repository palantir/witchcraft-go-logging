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
	"slices"

	"github.com/palantir/witchcraft-go-logging/wapi/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
)

type AuditResultType = logging.AuditResult

const (
	AuditResultSuccess      = logging.AuditResultSUCCESS
	AuditResultUnauthorized = logging.AuditResultUNAUTHORIZED
	AuditResultError        = logging.AuditResultERROR
)

type Logger interface {
	Audit(name string, result AuditResultType, params ...Param)
}

func New(w io.Writer, params ...Param) Logger {
	return &wrappedLogger{
		logger: &defaultLogger{logger: wlog.NewDefaultLogger(w, Type(), TimeNow())},
		params: params,
	}
}

func NewWithPrinter(printer wlog.ConjureLogPrinter[logging.AuditLogV2], params ...Param) Logger {
	return &wrappedLogger{
		logger: &defaultLogger{logger: wlog.NewDefaultLoggerWithPrinter(printer, Type(), TimeNow())},
		params: params,
	}
}

func WithParams(logger Logger, params ...Param) Logger {
	switch logger := logger.(type) {
	case *defaultLogger:
		return &wrappedLogger{
			logger: logger,
			params: slices.Clone(params),
		}
	case *wrappedLogger:
		return &wrappedLogger{
			logger: logger.logger,
			params: append(slices.Clone(logger.params), params...),
		}
	default:
		return &wrappedLogger{
			logger: logger,
			params: slices.Clone(params),
		}
	}
}
