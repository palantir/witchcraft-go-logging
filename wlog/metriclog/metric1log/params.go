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
	"maps"
	"time"

	"github.com/palantir/pkg/datetime"
	"github.com/palantir/witchcraft-go-logging/wapi/logging"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
)

const (
	TypeValue = "metric.1"
)

type Param = wloginternal.Param[logging.MetricLogV1]

type paramFunc = wloginternal.ParamFunc[logging.MetricLogV1]

func Type() Param {
	return paramFunc(func(l *logging.MetricLogV1) {
		l.Type = "metric.1"
	})
}

func Time(time time.Time) Param {
	return paramFunc(func(l *logging.MetricLogV1) {
		l.Time = datetime.DateTime(time)
	})
}

func TimeNow() Param {
	// Defer execution of time.Now() until the log is actually written
	return paramFunc(func(l *logging.MetricLogV1) {
		l.Time = datetime.DateTime(time.Now())
	})
}

func MetricName(name string) Param {
	return paramFunc(func(l *logging.MetricLogV1) {
		l.MetricName = name
	})
}

func MetricType(metric string) Param {
	return paramFunc(func(l *logging.MetricLogV1) {
		l.MetricType = metric
	})
}

func Values(params map[string]any) Param {
	return paramFunc(func(l *logging.MetricLogV1) {
		if l.Values == nil {
			l.Values = maps.Clone(params)
		} else {
			for k, v := range params {
				l.Values[k] = v
			}
		}
	})
}

func Value(key string, value any) Param {
	return paramFunc(func(l *logging.MetricLogV1) {
		if l.Values == nil {
			l.Values = map[string]any{key: value}
		} else {
			l.Values[key] = value
		}
	})
}

func Tags(tags map[string]string) Param {
	return paramFunc(func(l *logging.MetricLogV1) {
		if l.Tags == nil {
			l.Tags = maps.Clone(tags)
		} else {
			for k, v := range tags {
				l.Tags[k] = v
			}
		}
	})
}

func Tag(key, value string) Param {
	return paramFunc(func(l *logging.MetricLogV1) {
		if l.Tags == nil {
			l.Tags = map[string]string{key: value}
		} else {
			l.Tags[key] = value
		}
	})
}

func UID(uid string) Param {
	return paramFunc(func(l *logging.MetricLogV1) {
		l.Uid = (*logging.UserId)(&uid)
	})
}

func SID(sid string) Param {
	return paramFunc(func(l *logging.MetricLogV1) {
		l.Sid = (*logging.SessionId)(&sid)
	})
}

func TokenID(tokenId string) Param {
	return paramFunc(func(l *logging.MetricLogV1) {
		l.TokenId = (*logging.TokenId)(&tokenId)
	})
}

func OrgID(orgId string) Param {
	return paramFunc(func(l *logging.MetricLogV1) {
		l.OrgId = (*logging.OrgId)(&orgId)
	})
}

func UnsafeParams(unsafeParams map[string]any) Param {
	return paramFunc(func(l *logging.MetricLogV1) {
		if l.UnsafeParams == nil {
			l.UnsafeParams = maps.Clone(unsafeParams)
		} else {
			for k, v := range unsafeParams {
				l.UnsafeParams[k] = v
			}
		}
	})
}

func UnsafeParam(key string, value any) Param {
	return paramFunc(func(l *logging.MetricLogV1) {
		if l.UnsafeParams == nil {
			l.UnsafeParams = map[string]any{key: value}
		} else {
			l.UnsafeParams[key] = value
		}
	})
}
