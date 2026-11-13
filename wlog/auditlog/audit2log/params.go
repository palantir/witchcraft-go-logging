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
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/internal/auditloginternal"
)

const (
	TypeValue = auditloginternal.Audit2TypeValue

	OtherUIDsKey     = auditloginternal.Audit2OtherUIDsKey
	OriginKey        = auditloginternal.Audit2OriginKey
	NameKey          = auditloginternal.Audit2NameKey
	ResultKey        = auditloginternal.Audit2ResultKey
	RequestParamsKey = auditloginternal.Audit2RequestParamsKey
	ResultParamsKey  = auditloginternal.Audit2ResultParamsKey
)

type Param interface {
	getParam() auditloginternal.AuditParam
}

var _ Param = (*paramStruct)(nil)

type paramStruct struct {
	param auditloginternal.AuditParam
}

func (p *paramStruct) getParam() auditloginternal.AuditParam {
	return p.param
}

func ApplyParam(p Param, entry wlog.LogEntry) {
	if p == nil {
		return
	}
	auditloginternal.ApplyAudit2Param(p.getParam(), entry)
}

func convertInternalParamToExportedParam(param auditloginternal.AuditParam) Param {
	return &paramStruct{
		param: param,
	}
}

// internalParam returns the AuditParam of p, or an empty one if p is nil.
func internalParam(p Param) auditloginternal.AuditParam {
	if p == nil {
		return auditloginternal.AuditParam{}
	}
	return p.getParam()
}

func UID(uid string) Param {
	return convertInternalParamToExportedParam(auditloginternal.UID(uid))
}

func SID(sid string) Param {
	return convertInternalParamToExportedParam(auditloginternal.SID(sid))
}

func TokenID(tokenID string) Param {
	return convertInternalParamToExportedParam(auditloginternal.Audit2TokenID(tokenID))
}

func OrgID(orgID string) Param {
	return convertInternalParamToExportedParam(auditloginternal.OrgID(orgID))
}

func TraceID(traceID string) Param {
	return convertInternalParamToExportedParam(auditloginternal.TraceID(traceID))
}

func OtherUIDs(otherUIDs ...string) Param {
	return convertInternalParamToExportedParam(auditloginternal.Audit2OtherUIDs(otherUIDs...))
}

func Origin(origin string) Param {
	return convertInternalParamToExportedParam(auditloginternal.Origin(origin))
}

func RequestParam(key string, value any) Param {
	return convertInternalParamToExportedParam(auditloginternal.Audit2RequestParam(key, value))
}

func RequestParams(requestParams map[string]any) Param {
	return convertInternalParamToExportedParam(auditloginternal.Audit2RequestParams(requestParams))
}

func ResultParam(key string, value any) Param {
	return convertInternalParamToExportedParam(auditloginternal.Audit2ResultParam(key, value))
}

func ResultParams(resultParams map[string]any) Param {
	return convertInternalParamToExportedParam(auditloginternal.Audit2ResultParams(resultParams))
}
