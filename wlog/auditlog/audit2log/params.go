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
	"time"

	"github.com/palantir/pkg/datetime"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
	"github.com/palantir/witchcraft-go-logging/wtypes"
)

const (
	TypeValue = "audit.2"
)

type Param = wloginternal.Param[wtypes.AuditLogV2]

type paramFunc = wloginternal.ParamFunc[wtypes.AuditLogV2]

func defaultParam(name string, result AuditResultType) Param {
	return paramFunc(func(l *wtypes.AuditLogV2) {
		l.Type = TypeValue
		l.Time = datetime.DateTime(time.Now())
		l.Name = name
		l.Result = result
	})
}

func UID(uid string) Param {
	return paramFunc(func(l *wtypes.AuditLogV2) {
		l.Uid = (*wtypes.UserId)(&uid)
	})
}

func SID(sid string) Param {
	return paramFunc(func(l *wtypes.AuditLogV2) {
		l.Sid = (*wtypes.SessionId)(&sid)
	})
}

func TokenID(tokenID string) Param {
	return paramFunc(func(l *wtypes.AuditLogV2) {
		l.TokenId = (*wtypes.TokenId)(&tokenID)
	})
}

func OrgID(orgID string) Param {
	return paramFunc(func(l *wtypes.AuditLogV2) {
		l.OrgId = (*wtypes.OrgId)(&orgID)
	})
}

func TraceID(traceID string) Param {
	return paramFunc(func(l *wtypes.AuditLogV2) {
		l.TraceId = (*wtypes.TraceId)(&traceID)
	})
}

func OtherUIDs(otherUIDs ...string) Param {
	return paramFunc(func(l *wtypes.AuditLogV2) {
		for _, uid := range otherUIDs {
			l.OtherUids = append(l.OtherUids, wtypes.UserId(uid))
		}
	})
}

func Origin(origin string) Param {
	return paramFunc(func(l *wtypes.AuditLogV2) {
		l.Origin = &origin
	})
}

func RequestParam(key string, value interface{}) Param {
	return paramFunc(func(l *wtypes.AuditLogV2) {
		wloginternal.SetMapParam(&l.RequestParams, key, value)
	})
}

func RequestParams(requestParams map[string]interface{}) Param {
	return paramFunc(func(l *wtypes.AuditLogV2) {
		wloginternal.SetMapParams(&l.RequestParams, requestParams)
	})
}

func ResultParam(key string, value interface{}) Param {
	return paramFunc(func(l *wtypes.AuditLogV2) {
		wloginternal.SetMapParam(&l.ResultParams, key, value)
	})
}

func ResultParams(resultParams map[string]interface{}) Param {
	return paramFunc(func(l *wtypes.AuditLogV2) {
		wloginternal.SetMapParams(&l.ResultParams, resultParams)
	})
}
