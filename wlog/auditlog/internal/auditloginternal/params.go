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
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logtypes"
)

// AuditParam holds the per-format mutators for an audit parameter. A single Param can mutate either or both of the
// audit.v2 and audit.v3 log structs depending on which fields it represents.
type AuditParam struct {
	Audit2ParamFn func(*logtypes.AuditLogV2)
	Audit3ParamFn func(*logtypes.AuditLogV3)
}

func UID(uid string) AuditParam {
	return AuditParam{
		Audit2ParamFn: func(log *logtypes.AuditLogV2) { log.Uid = uid },
		Audit3ParamFn: func(log *logtypes.AuditLogV3) { log.Uid = uid },
	}
}

func SID(sid string) AuditParam {
	return AuditParam{
		Audit2ParamFn: func(log *logtypes.AuditLogV2) { log.Sid = sid },
		Audit3ParamFn: func(log *logtypes.AuditLogV3) { log.Sid = sid },
	}
}

func OrgID(orgID string) AuditParam {
	return AuditParam{
		Audit2ParamFn: func(log *logtypes.AuditLogV2) { log.OrgId = orgID },
		Audit3ParamFn: func(log *logtypes.AuditLogV3) { log.OrgId = orgID },
	}
}

func TraceID(traceID string) AuditParam {
	return AuditParam{
		Audit2ParamFn: func(log *logtypes.AuditLogV2) { log.TraceId = traceID },
		Audit3ParamFn: func(log *logtypes.AuditLogV3) { log.TraceId = traceID },
	}
}

func Origin(origin string) AuditParam {
	return AuditParam{
		Audit2ParamFn: func(log *logtypes.AuditLogV2) { log.Origin = origin },
		Audit3ParamFn: func(log *logtypes.AuditLogV3) { log.Origin = origin },
	}
}
