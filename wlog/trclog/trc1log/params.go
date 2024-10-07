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

package trc1log

import (
	"time"

	"github.com/palantir/pkg/datetime"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
	"github.com/palantir/witchcraft-go-logging/wtypes"
)

const (
	TypeValue = "trace.1"
)

type Param = wloginternal.Param[wtypes.TraceLogV1]

type paramFunc = wloginternal.ParamFunc[wtypes.TraceLogV1]

func Type() Param {
	return paramFunc(func(l *wtypes.TraceLogV1) {
		l.Type = TypeValue
	})
}

func Time(time time.Time) Param {
	return paramFunc(func(l *wtypes.TraceLogV1) {
		l.Time = datetime.DateTime(time)
	})
}

func TimeNow() Param {
	// Defer execution of time.Now() until the log is actually written
	return paramFunc(func(l *wtypes.TraceLogV1) {
		l.Time = datetime.DateTime(time.Now())
	})
}

func UID(uid string) Param {
	return paramFunc(func(l *wtypes.TraceLogV1) {
		l.Uid = (*wtypes.UserID)(&uid)
	})
}

func SID(sid string) Param {
	return paramFunc(func(l *wtypes.TraceLogV1) {
		l.Sid = (*wtypes.SessionID)(&sid)
	})
}

func TokenID(tokenID string) Param {
	return paramFunc(func(l *wtypes.TraceLogV1) {
		l.TokenId = (*wtypes.TokenID)(&tokenID)
	})
}

func OrgID(orgID string) Param {
	return paramFunc(func(l *wtypes.TraceLogV1) {
		l.OrgId = (*wtypes.OrgID)(&orgID)
	})
}

func UnsafeParam(key string, value interface{}) Param {
	return paramFunc(func(l *wtypes.TraceLogV1) {
		wloginternal.SetMapParam(&l.UnsafeParams, key, value)
	})
}

func UnsafeParams(unsafe map[string]interface{}) Param {
	return paramFunc(func(l *wtypes.TraceLogV1) {
		wloginternal.SetMapParams(&l.UnsafeParams, unsafe)
	})
}
