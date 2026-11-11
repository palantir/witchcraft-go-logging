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

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/internal/auditloginternal"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
)

type AuditResultType string

const (
	AuditResultSuccess      = AuditResultType(auditloginternal.AuditResultSuccess)
	AuditResultUnauthorized = AuditResultType(auditloginternal.AuditResultUnauthorized)
	AuditResultError        = AuditResultType(auditloginternal.AuditResultError)
)

type Logger interface {
	Audit(name string, result AuditResultType, params ...Param)
}

func New(w io.Writer) Logger {
	return NewFromCreator(w, wlog.DefaultLoggerProvider().NewLogger)
}

func NewDualLogger(audit2Writer, audit3Writer io.Writer) Logger {
	return NewDualLoggerFromCreator(audit2Writer, audit3Writer, wlog.DefaultLoggerProvider().NewLogger)
}

func NewFromCreator(w io.Writer, creator wlog.LoggerCreator) Logger {
	return NewDualLoggerFromCreator(w, nil, creator)
}

func NewDualLoggerFromCreator(audit2Writer, audit3Writer io.Writer, creator wlog.LoggerCreator) Logger {
	return &defaultLogger{logger: auditloginternal.NewLogger(creator(audit2Writer), auditloginternal.CreateIfSet(audit3Writer, creator))}
}

func WithParams(logger Logger, params ...Param) Logger {
	if len(params) == 0 {
		return logger
	}

	if innerWrapped, ok := logger.(*wrappedLogger); ok {
		return &wrappedLogger{
			logger: innerWrapped.logger,
			params: wloginternal.AppendParams(innerWrapped.params, params),
		}
	}

	return &wrappedLogger{
		logger: logger,
		params: params,
	}
}
