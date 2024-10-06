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

package evt2log

import (
	"time"

	"github.com/palantir/pkg/datetime"
	"github.com/palantir/witchcraft-go-logging/wapi/logging"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
)

const (
	TypeValue = "event.2"
)

type Param = wloginternal.Param[logging.EventLogV2]

type paramFunc = wloginternal.ParamFunc[logging.EventLogV2]

func defaultParam(name string) Param {
	return paramFunc(func(l *logging.EventLogV2) {
		l.Type = TypeValue
		l.Time = datetime.DateTime(time.Now())
		l.EventName = name
	})
}

func Value(key string, value interface{}) Param {
	return paramFunc(func(l *logging.EventLogV2) {
		wloginternal.SetMapParam(&l.Values, key, value)
	})
}

func Values(values map[string]interface{}) Param {
	return paramFunc(func(l *logging.EventLogV2) {
		wloginternal.SetMapParams(&l.Values, values)
	})
}

func Tags(values map[string]string) Param {
	return paramFunc(func(l *logging.EventLogV2) {
		wloginternal.SetMapParams(&l.Tags, values)
	})
}

func Tag(key, value string) Param {
	return paramFunc(func(l *logging.EventLogV2) {
		wloginternal.SetMapParam(&l.Tags, key, value)
	})
}

func UID(uid string) Param {
	return paramFunc(func(l *logging.EventLogV2) {
		l.Uid = (*logging.UserId)(&uid)
	})
}

func SID(sid string) Param {
	return paramFunc(func(l *logging.EventLogV2) {
		l.Sid = (*logging.SessionId)(&sid)
	})
}

func TokenID(tokenID string) Param {
	return paramFunc(func(l *logging.EventLogV2) {
		l.TokenId = (*logging.TokenId)(&tokenID)
	})
}

func OrgID(orgID string) Param {
	return paramFunc(func(l *logging.EventLogV2) {
		l.OrgId = (*logging.OrgId)(&orgID)
	})
}

func TraceID(traceID string) Param {
	return paramFunc(func(l *logging.EventLogV2) {
		l.TraceId = (*logging.TraceId)(&traceID)
	})
}

func UnsafeParams(unsafe map[string]interface{}) Param {
	return paramFunc(func(l *logging.EventLogV2) {
		wloginternal.SetMapParams(&l.UnsafeParams, unsafe)
	})
}

func UnsafeParam(key string, value interface{}) Param {
	return paramFunc(func(l *logging.EventLogV2) {
		wloginternal.SetMapParam(&l.UnsafeParams, key, value)
	})
}
