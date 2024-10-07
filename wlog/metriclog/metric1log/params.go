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

package metric1log

import (
	"time"

	"github.com/palantir/pkg/datetime"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
	"github.com/palantir/witchcraft-go-logging/wtypes"
)

const (
	TypeValue = "metric.1"
)

type Param = wloginternal.Param[wtypes.MetricLogV1]

type paramFunc = wloginternal.ParamFunc[wtypes.MetricLogV1]

func defaultParam(name string, typ string) Param {
	return paramFunc(func(l *wtypes.MetricLogV1) {
		l.Type = TypeValue
		l.Time = datetime.DateTime(time.Now())
		l.MetricName = name
		l.MetricType = typ
	})
}

func Value(key string, value interface{}) Param {
	return paramFunc(func(l *wtypes.MetricLogV1) {
		wloginternal.SetMapParam(&l.Values, key, value)
	})
}

func Values(values map[string]interface{}) Param {
	return paramFunc(func(l *wtypes.MetricLogV1) {
		wloginternal.SetMapParams(&l.Values, values)
	})
}

func Tag(key, value string) Param {
	return paramFunc(func(l *wtypes.MetricLogV1) {
		wloginternal.SetMapParam(&l.Tags, key, value)
	})
}

func Tags(values map[string]string) Param {
	return paramFunc(func(l *wtypes.MetricLogV1) {
		wloginternal.SetMapParams(&l.Tags, values)
	})
}

func UID(uid string) Param {
	return paramFunc(func(l *wtypes.MetricLogV1) {
		l.Uid = (*wtypes.UserId)(&uid)
	})
}

func SID(sid string) Param {
	return paramFunc(func(l *wtypes.MetricLogV1) {
		l.Sid = (*wtypes.SessionId)(&sid)
	})
}

func TokenID(tokenID string) Param {
	return paramFunc(func(l *wtypes.MetricLogV1) {
		l.TokenId = (*wtypes.TokenId)(&tokenID)
	})
}

func OrgID(orgID string) Param {
	return paramFunc(func(l *wtypes.MetricLogV1) {
		l.OrgId = (*wtypes.OrgId)(&orgID)
	})
}

func UnsafeParam(key string, value interface{}) Param {
	return paramFunc(func(l *wtypes.MetricLogV1) {
		wloginternal.SetMapParam(&l.UnsafeParams, key, value)
	})
}

func UnsafeParams(unsafe map[string]interface{}) Param {
	return paramFunc(func(l *wtypes.MetricLogV1) {
		wloginternal.SetMapParams(&l.UnsafeParams, unsafe)
	})
}
