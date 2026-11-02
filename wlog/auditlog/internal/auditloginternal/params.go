// Copyright (c) 2025 Palantir Technologies. All rights reserved.
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

package auditloginternal

import (
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
)

// AuditParam holds the per-format mutators for an audit parameter. A single Param can mutate either or both of the
// audit.v2 and audit.v3 log structs depending on which fields it represents.
type AuditParam struct {
	Audit2ParamFn func(*logging.AuditLogV2)
	Audit3ParamFn func(*logging.AuditLogV3)
	LogEntryID    *string
}

func uidParam(uid string) AuditParam {
	return optStringParam(uid,
		func(log *logging.AuditLogV2, v *logging.UserId) { log.Uid = v },
		func(log *logging.AuditLogV3, v *logging.UserId) { log.Uid = v })
}

func sidParam(sid string) AuditParam {
	return optStringParam(sid,
		func(log *logging.AuditLogV2, v *logging.SessionId) { log.Sid = v },
		func(log *logging.AuditLogV3, v *logging.SessionId) { log.Sid = v })
}

func tokenIDParam(tokenID string) AuditParam {
	return optStringParam(tokenID,
		func(log *logging.AuditLogV2, v *logging.TokenId) { log.TokenId = v },
		func(log *logging.AuditLogV3, v *logging.TokenId) { log.TokenId = v })
}

func orgIDParam(orgID string) AuditParam {
	return optStringParam(orgID,
		func(log *logging.AuditLogV2, v *logging.OrganizationId) { log.OrgId = v },
		func(log *logging.AuditLogV3, v *logging.OrganizationId) { log.OrgId = v })
}

func traceIDParam(traceID string) AuditParam {
	return optStringParam(traceID,
		func(log *logging.AuditLogV2, v *logging.TraceId) { log.TraceId = v },
		func(log *logging.AuditLogV3, v *logging.TraceId) { log.TraceId = v })
}

func originParam(origin string) AuditParam {
	return optStringParam(origin,
		func(log *logging.AuditLogV2, v *string) { log.Origin = v },
		func(log *logging.AuditLogV3, v *string) { log.Origin = v })
}

// optStringParam builds an AuditParam that sets an optional string field (of newtype T) on both the audit.v2 and
// audit.v3 structs. Empty values clear previously supplied values.
func optStringParam[T ~string](value string, set2 func(*logging.AuditLogV2, *T), set3 func(*logging.AuditLogV3, *T)) AuditParam {
	v := T(value)
	return AuditParam{
		Audit2ParamFn: func(log *logging.AuditLogV2) { set2(log, &v) },
		Audit3ParamFn: func(log *logging.AuditLogV3) { set3(log, &v) },
	}
}
