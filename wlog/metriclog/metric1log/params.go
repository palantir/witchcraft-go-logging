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
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
)

const (
	TypeValue = "metric.1"

	MetricNameKey = "metricName"
	MetricTypeKey = "metricType"
	ValuesKey     = "values"
	TagsKey       = "tags"
)

type Param = wloginternal.Param[logging.MetricLogV1]

type paramFunc = wloginternal.ParamFunc[logging.MetricLogV1]

func ApplyParam(p Param, entry wlog.LogEntry) {
	if p == nil {
		return
	}
	log := new(logging.MetricLogV1)
	wloginternal.ApplyParam(log, p)
	if _, ok := p.(timeParam); ok {
		entry.StringValue(wlog.TimeKey, log.Time.String())
	}
	applyObjectParams(log, entry)
}

// populateLogObject pre-populates the type, time, metric name, and metric type fields and applies user params in-place.
func populateLogObject(log *logging.MetricLogV1, metricName, metricType string, params []Param) {
	log.Type = TypeValue
	log.Time = datetime.DateTime(time.Now())
	log.MetricName = metricName
	log.MetricType = metricType
	wloginternal.ApplyParams(log, params)
	wloginternal.OmitEmpty(&log.Uid)
	wloginternal.OmitEmpty(&log.Sid)
	wloginternal.OmitEmpty(&log.TokenId)
	wloginternal.OmitEmpty(&log.OrgId)
}

func objectToParam(log *logging.MetricLogV1) wlog.Param {
	return wlog.NewParam(func(entry wlog.LogEntry) {
		entry.StringValue(wlog.TypeKey, log.Type)
		entry.StringValue(wlog.TimeKey, log.Time.String())
		entry.StringValue(MetricNameKey, log.MetricName)
		entry.StringValue(MetricTypeKey, log.MetricType)
		applyObjectParams(log, entry)
	})
}

func applyObjectParams(log *logging.MetricLogV1, entry wlog.LogEntry) {
	if len(log.Values) > 0 {
		entry.AnyMapValue(ValuesKey, log.Values)
	}
	if len(log.Tags) > 0 {
		entry.StringMapValue(TagsKey, log.Tags)
	}
	if log.Uid != nil {
		entry.OptionalStringValue(wlog.UIDKey, string(*log.Uid))
	}
	if log.Sid != nil {
		entry.OptionalStringValue(wlog.SIDKey, string(*log.Sid))
	}
	if log.TokenId != nil {
		entry.OptionalStringValue(wlog.TokenIDKey, string(*log.TokenId))
	}
	if log.OrgId != nil {
		entry.OptionalStringValue(wlog.OrgIDKey, string(*log.OrgId))
	}
	if len(log.UnsafeParams) > 0 {
		entry.AnyMapValue(wlog.UnsafeParamsKey, log.UnsafeParams)
	}
}

type timeParam struct {
	paramFunc
}

// Time overrides the time at which the metric was recorded. By default, metrics use the time at which they are logged.
func Time(value time.Time) Param {
	return timeParam{paramFunc(func(log *logging.MetricLogV1) {
		log.Time = datetime.DateTime(value)
	})}
}

func Value(key string, value any) Param {
	return paramFunc(func(log *logging.MetricLogV1) {
		wloginternal.SetMapParam(&log.Values, key, value)
	})
}

func Values(values map[string]any) Param {
	return paramFunc(func(log *logging.MetricLogV1) {
		wloginternal.SetMapParams(&log.Values, values)
	})
}

func Tag(key, value string) Param {
	return paramFunc(func(log *logging.MetricLogV1) {
		wloginternal.SetMapParam(&log.Tags, key, value)
	})
}

func Tags(values map[string]string) Param {
	return paramFunc(func(log *logging.MetricLogV1) {
		wloginternal.SetMapParams(&log.Tags, values)
	})
}

func UID(uid string) Param {
	return paramFunc(func(log *logging.MetricLogV1) {
		log.Uid = (*logging.UserId)(&uid)
	})
}

func SID(sid string) Param {
	return paramFunc(func(log *logging.MetricLogV1) {
		log.Sid = (*logging.SessionId)(&sid)
	})
}

func TokenID(tokenID string) Param {
	return paramFunc(func(log *logging.MetricLogV1) {
		log.TokenId = (*logging.TokenId)(&tokenID)
	})
}

func OrgID(orgID string) Param {
	return paramFunc(func(log *logging.MetricLogV1) {
		log.OrgId = (*logging.OrganizationId)(&orgID)
	})
}

func UnsafeParam(key string, value any) Param {
	return paramFunc(func(log *logging.MetricLogV1) {
		wloginternal.SetMapParam(&log.UnsafeParams, key, value)
	})
}

func UnsafeParams(unsafe map[string]any) Param {
	return paramFunc(func(log *logging.MetricLogV1) {
		wloginternal.SetMapParams(&log.UnsafeParams, unsafe)
	})
}
