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
	"time"

	"github.com/palantir/witchcraft-go-logging/wlog"
)

const (
	Audit2TypeValue = "audit.2"

	Audit2OtherUIDsKey     = "otherUids"
	Audit2OriginKey        = "origin"
	Audit2NameKey          = "name"
	Audit2ResultKey        = "result"
	Audit2RequestParamsKey = "requestParams"
	Audit2ResultParamsKey  = "resultParams"
)

// Audit2Param is effectively a marker type that wraps AuditParam.
// The AuditParam interface requires specifying outputs for both audit.v2 and audit.v3, but the Audit2Param type should
// only be used for parameters that are "native" to audit.v2.
type Audit2Param struct {
	Param AuditParam
}

func Audit2ToParams(name string, result AuditResultType, inParams []Audit2Param) []wlog.Param {
	outParams := make([]wlog.Param, 1+len(inParams))
	outParams[0] = wlog.NewParam(func(entry wlog.LogEntry) {
		entry.StringValue(wlog.TypeKey, Audit2TypeValue)
		entry.StringValue(wlog.TimeKey, time.Now().Format(time.RFC3339Nano))
		auditNameResultParam(name, result).Audit2ParamFn(entry)
	})
	for idx := range inParams {
		outParams[1+idx] = wlog.NewParam(inParams[idx].Param.Audit2ParamFn)
	}
	return outParams
}

func Audit2UID(uid string) Audit2Param {
	return Audit2Param{
		Param: uidParam(uid),
	}
}

func Audit2SID(sid string) Audit2Param {
	return Audit2Param{
		Param: sidParam(sid),
	}
}

func Audit2TokenID(tokenID string) Audit2Param {
	return Audit2Param{
		Param: tokenIDParam(tokenID),
	}
}

func Audit2OrgID(orgID string) Audit2Param {
	return Audit2Param{
		Param: orgIDParam(orgID),
	}
}

func Audit2TraceID(traceID string) Audit2Param {
	return Audit2Param{
		Param: traceIDParam(traceID),
	}
}

func Audit2OtherUIDs(otherUIDs ...string) Audit2Param {
	return Audit2Param{
		Param: AuditParam{
			Audit2ParamFn: func(entry wlog.LogEntry) {
				entry.StringListValue(Audit2OtherUIDsKey, otherUIDs)
			},
		},
	}
}

func Audit2Origin(origin string) Audit2Param {
	return Audit2Param{
		Param: originParam(origin),
	}
}

func Audit2RequestParam(key string, value interface{}) Audit2Param {
	return Audit2RequestParams(map[string]interface{}{
		key: value,
	})
}

func Audit2RequestParams(requestParams map[string]interface{}) Audit2Param {
	return Audit2Param{
		Param: AuditParam{
			Audit2ParamFn: func(entry wlog.LogEntry) {
				entry.AnyMapValue(Audit2RequestParamsKey, requestParams)
			},
		},
	}
}

func Audit2ResultParam(key string, value interface{}) Audit2Param {
	return Audit2ResultParams(map[string]interface{}{
		key: value,
	})
}

func Audit2ResultParams(resultParams map[string]interface{}) Audit2Param {
	return Audit2Param{
		Param: AuditParam{
			Audit2ParamFn: func(entry wlog.LogEntry) {
				entry.AnyMapValue(Audit2ResultParamsKey, resultParams)
			},
		},
	}
}
