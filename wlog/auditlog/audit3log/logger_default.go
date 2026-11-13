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
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/internal/auditloginternal"
)

type defaultLogger struct {
	logger *auditloginternal.Logger
}

func (l *defaultLogger) Audit(name string, result AuditResultType, params ...Param) {
	auditloginternal.Audit3(l.logger, name, auditloginternal.AuditResultType(result), params, internalParam)
}

// ToParams builds parameters backed by a fresh object that can outlive the log call.
func ToParams(name string, result AuditResultType, inParams []Param) []wlog.Param {
	return auditloginternal.Audit3ToParams(name, auditloginternal.AuditResultType(result), inParams, internalParam)
}
