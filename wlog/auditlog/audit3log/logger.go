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

package audit3log

import (
	"io"

	"github.com/palantir/witchcraft-go-logging/wlog"
)

type AuditResultType string
type AuditProducerType string
type AuditSensitivityType string

type AuditOrganizationType struct {
	ID     string `json:"id"`
	Reason string `json:"reason"`
}

type AuditContextualizedUserType struct {
	UID       string   `json:"uid"`
	UserName  string   `json:"userName"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Groups    []string `json:"groups"`
	Realm     string   `json:"realm"`
}

type AuditSensitivityTaggedValueType struct {
	Level   []AuditSensitivityType `json:"level"`
	Payload interface{}            `json:"payload"`
}

const (
	AuditResultSuccess        AuditResultType      = "SUCCESS"
	AuditResultUnauthorized   AuditResultType      = "UNAUTHORIZED"
	AuditResultError          AuditResultType      = "ERROR"
	AuditProducerServer       AuditProducerType    = "SERVER"
	AuditProducerClient       AuditProducerType    = "CLIENT"
	AuditSensitivityMetadata  AuditSensitivityType = "Metadata"
	AuditSensitivityUserInput AuditSensitivityType = "UserInput"
	AuditSensitivityData      AuditSensitivityType = "Data"
)

type Logger interface {
	Audit(name string, result AuditResultType, deployment string, product string, productVersion string, params ...Param)
}

func New(w io.Writer) Logger {
	return NewFromCreator(w, wlog.DefaultLoggerProvider().NewLogger)
}

func NewFromCreator(w io.Writer, creator wlog.LoggerCreator) Logger {
	return &defaultLogger{
		logger: creator(w),
	}
}

func WithParams(logger Logger, params ...Param) Logger {
	if len(params) == 0 {
		return logger
	}

	if innerWrapped, ok := logger.(*wrappedLogger); ok {
		return &wrappedLogger{
			logger: innerWrapped.logger,
			params: append(innerWrapped.params, params...),
		}
	}

	return &wrappedLogger{
		logger: logger,
		params: params,
	}
}
