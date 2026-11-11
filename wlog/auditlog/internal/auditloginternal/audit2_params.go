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
	"maps"

	"github.com/palantir/pkg/datetime"
	"github.com/palantir/witchcraft-go-logging/wlog"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logentry"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logtypes"
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

// Audit2ToParams builds parameters backed by a fresh audit.2 log that can outlive the log call.
func Audit2ToParams[P any](name string, result AuditResultType, params []P, get func(P) AuditParam) []wlog.Param {
	log := new(logtypes.AuditLogV2)
	populateAuditLogV2(log, &atLogTimeValues{}, name, result, params, get)
	return []wlog.Param{logentry.Param(log)}
}

// populateAuditLogV2 populates a caller-provided AuditLogV2 with the required fields, applies the provided params, and
// propagates a sibling v3 emission's logEntryId if one was set on values.
func populateAuditLogV2[P any](log *logtypes.AuditLogV2, values *atLogTimeValues, name string, result AuditResultType, params []P, get func(P) AuditParam) {
	log.Type = Audit2TypeValue
	log.Time = datetime.DateTime(values.getOrComputeTimeValue())
	log.Name = name
	log.Result = logtypes.AuditResult(result)
	for _, p := range params {
		if fn := get(p).Audit2ParamFn; fn != nil {
			fn(log)
		}
	}
	if id := values.logEntryID; id != "" {
		wloginternal.SetMapParam(&log.RequestParams, "_auditLogEntryId", any(id))
	}
}

// ApplyAudit2Param applies one parameter's audit.2 fields to entry.
func ApplyAudit2Param(param AuditParam, entry wlog.LogEntry) {
	if param.Audit2ParamFn == nil {
		return
	}
	log := new(logtypes.AuditLogV2)
	param.Audit2ParamFn(log)
	logentry.AuditLogV2Params(entry, log)
}

func Audit2TokenID(tokenID string) AuditParam {
	return AuditParam{
		Audit2ParamFn: func(log *logtypes.AuditLogV2) { log.TokenId = tokenID },
		Audit3ParamFn: func(log *logtypes.AuditLogV3) { log.TokenId = tokenID },
	}
}

func Audit2OtherUIDs(otherUIDs ...string) AuditParam {
	return AuditParam{
		Audit2ParamFn: func(log *logtypes.AuditLogV2) {
			log.OtherUids = otherUIDs
		},
	}
}

func Audit2RequestParam(key string, value any) AuditParam {
	return AuditParam{
		Audit2ParamFn: func(log *logtypes.AuditLogV2) {
			wloginternal.SetMapParam(&log.RequestParams, key, value)
		},
		Audit3ParamFn: func(log *logtypes.AuditLogV3) {
			if category, ok := value.(string); key == "category" && ok {
				log.Categories = append(log.Categories, category)
				return
			}
			wloginternal.SetMapParam(&log.RequestFields, key, value)
		},
	}
}

func Audit2RequestParams(requestParams map[string]any) AuditParam {
	return AuditParam{
		Audit2ParamFn: func(log *logtypes.AuditLogV2) {
			wloginternal.SetMapParams(&log.RequestParams, requestParams)
		},
		Audit3ParamFn: func(log *logtypes.AuditLogV3) {
			// In v3, "category" gets pulled out into Categories; the rest goes to RequestFields.
			params := requestParams
			if categoryVal, ok := requestParams["category"]; ok {
				if categoryValString, ok := categoryVal.(string); ok {
					params = maps.Clone(requestParams)
					delete(params, "category")
					log.Categories = append(log.Categories, categoryValString)
				}
			}
			wloginternal.SetMapParams(&log.RequestFields, params)
		},
	}
}

func Audit2ResultParam(key string, value any) AuditParam {
	return AuditParam{
		Audit2ParamFn: func(log *logtypes.AuditLogV2) {
			wloginternal.SetMapParam(&log.ResultParams, key, value)
		},
		Audit3ParamFn: func(log *logtypes.AuditLogV3) {
			wloginternal.SetMapParam(&log.ResultFields, key, value)
		},
	}
}

func Audit2ResultParams(resultParams map[string]any) AuditParam {
	return AuditParam{
		Audit2ParamFn: func(log *logtypes.AuditLogV2) {
			wloginternal.SetMapParams(&log.ResultParams, resultParams)
		},
		Audit3ParamFn: func(log *logtypes.AuditLogV3) {
			wloginternal.SetMapParams(&log.ResultFields, resultParams)
		},
	}
}
