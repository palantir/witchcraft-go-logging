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
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/internal/auditloginternal"
)

// ToParams builds parameters backed by a fresh object that can outlive the log call.
func ToParams(name string, result AuditResultType, inParams []Param) []wlog.Param {
	return auditloginternal.Audit3ToParams(nil, name, auditloginternal.AuditResultType(result), convertExternalParamsToInternalParams(inParams))
}

// ToDirectObject builds an audit object. A non-nil fallback parameter must be emitted instead of the object
// to preserve missing event IDs and noncanonical log-entry IDs. Both outputs share the same populated data.
func ToDirectObject(name string, result AuditResultType, inParams []Param) (logging.AuditLogV3, wlog.Param) {
	return auditloginternal.Audit3ToDirectObject(name, auditloginternal.AuditResultType(result), convertExternalParamsToInternalParams(inParams))
}
