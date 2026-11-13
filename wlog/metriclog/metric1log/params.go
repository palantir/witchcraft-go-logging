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
	"github.com/palantir/witchcraft-go-logging/wlog"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logentry"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logtypes"
)

const (
	TypeValue = "metric.1"

	MetricNameKey = "metricName"
	MetricTypeKey = "metricType"
	ValuesKey     = "values"
	TagsKey       = "tags"
)

type Param = wloginternal.Param[logtypes.MetricLogV1]

type paramFunc = wloginternal.ParamFunc[logtypes.MetricLogV1]

func ApplyParam(p Param, entry wlog.LogEntry) {
	if p == nil {
		return
	}
	log := new(logtypes.MetricLogV1)
	wloginternal.ApplyParam(log, p)
	if _, ok := p.(timeParam); ok {
		entry.StringValue(wlog.TimeKey, log.Time.String())
	}
	logentry.MetricLogV1Params(entry, log)
}

// populateLogObject pre-populates the type, time, metric name, and metric type fields and applies user params in-place.
func populateLogObject(log *logtypes.MetricLogV1, metricName, metricType string, params []Param) {
	log.Type = TypeValue
	log.Time = datetime.DateTime(time.Now())
	log.MetricName = metricName
	log.MetricType = metricType
	wloginternal.ApplyParams(log, params)
}

type timeParam struct {
	paramFunc
}

// Time overrides the time at which the metric was recorded. By default, metrics use the time at which they are logged.
func Time(value time.Time) Param {
	return timeParam{paramFunc(func(log *logtypes.MetricLogV1) {
		log.Time = datetime.DateTime(value)
	})}
}

func Value(key string, value any) Param {
	return paramFunc(func(log *logtypes.MetricLogV1) {
		wloginternal.SetMapParam(&log.Values, key, value)
	})
}

func Values(values map[string]any) Param {
	return paramFunc(func(log *logtypes.MetricLogV1) {
		wloginternal.SetMapParams(&log.Values, values)
	})
}

func Tag(key, value string) Param {
	return paramFunc(func(log *logtypes.MetricLogV1) {
		wloginternal.SetMapParam(&log.Tags, key, value)
	})
}

func Tags(values map[string]string) Param {
	return paramFunc(func(log *logtypes.MetricLogV1) {
		wloginternal.SetMapParams(&log.Tags, values)
	})
}

func UID(uid string) Param {
	return paramFunc(func(log *logtypes.MetricLogV1) {
		log.Uid = uid
	})
}

func SID(sid string) Param {
	return paramFunc(func(log *logtypes.MetricLogV1) {
		log.Sid = sid
	})
}

func TokenID(tokenID string) Param {
	return paramFunc(func(log *logtypes.MetricLogV1) {
		log.TokenId = tokenID
	})
}

func OrgID(orgID string) Param {
	return paramFunc(func(log *logtypes.MetricLogV1) {
		log.OrgId = orgID
	})
}

func UnsafeParam(key string, value any) Param {
	return paramFunc(func(log *logtypes.MetricLogV1) {
		wloginternal.SetMapParam(&log.UnsafeParams, key, value)
	})
}

func UnsafeParams(unsafe map[string]any) Param {
	return paramFunc(func(log *logtypes.MetricLogV1) {
		wloginternal.SetMapParams(&log.UnsafeParams, unsafe)
	})
}
