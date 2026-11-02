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
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
)

const (
	TypeValue = "event.2"

	EventNameKey = "eventName"
	ValuesKey    = "values"
	TagsKey      = "tags"
)

type Param = wloginternal.Param[logging.EventLogV2]

type paramFunc = wloginternal.ParamFunc[logging.EventLogV2]

func ApplyParam(p Param, entry wlog.LogEntry) {
	if p == nil {
		return
	}
	log := new(logging.EventLogV2)
	wloginternal.ApplyParam(log, p)
	applyObjectParams(log, entry)
}

// populateLogObject pre-populates the type, time, and event-name fields and applies user params in-place.
// Plain function (not a Param-returning closure) to avoid an extra heap allocation per log call.
func populateLogObject(log *logging.EventLogV2, name string, params []Param) {
	log.Type = TypeValue
	log.Time = datetime.DateTime(time.Now())
	log.EventName = name
	wloginternal.ApplyParams(log, params)
	wloginternal.OmitEmpty(&log.Uid)
	wloginternal.OmitEmpty(&log.Sid)
	wloginternal.OmitEmpty(&log.TokenId)
	wloginternal.OmitEmpty(&log.OrgId)
	wloginternal.OmitEmpty(&log.TraceId)
}

func objectToParam(log *logging.EventLogV2) wlog.Param {
	return wlog.NewParam(func(entry wlog.LogEntry) {
		entry.StringValue(wlog.TypeKey, log.Type)
		entry.StringValue(wlog.TimeKey, log.Time.String())
		entry.OptionalStringValue(EventNameKey, log.EventName)
		applyObjectParams(log, entry)
	})
}

func applyObjectParams(log *logging.EventLogV2, entry wlog.LogEntry) {
	if len(log.Values) > 0 {
		entry.AnyMapValue(ValuesKey, log.Values)
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
	if log.TraceId != nil {
		entry.OptionalStringValue(wlog.TraceIDKey, string(*log.TraceId))
	}
	if len(log.UnsafeParams) > 0 {
		entry.AnyMapValue(wlog.UnsafeParamsKey, log.UnsafeParams)
	}
	if len(log.Tags) > 0 {
		entry.StringMapValue(TagsKey, log.Tags)
	}
}
func Value(key string, value any) Param {
	return paramFunc(func(log *logging.EventLogV2) {
		wloginternal.SetMapParam(&log.Values, key, value)
	})
}

func Values(values map[string]any) Param {
	return paramFunc(func(log *logging.EventLogV2) {
		wloginternal.SetMapParams(&log.Values, values)
	})
}

func Tag(key, value string) Param {
	return paramFunc(func(log *logging.EventLogV2) {
		wloginternal.SetMapParam(&log.Tags, key, value)
	})
}

func Tags(values map[string]string) Param {
	return paramFunc(func(log *logging.EventLogV2) {
		wloginternal.SetMapParams(&log.Tags, values)
	})
}

func UID(uid string) Param {
	return paramFunc(func(log *logging.EventLogV2) {
		log.Uid = (*logging.UserId)(&uid)
	})
}

func SID(sid string) Param {
	return paramFunc(func(log *logging.EventLogV2) {
		log.Sid = (*logging.SessionId)(&sid)
	})
}

func TokenID(tokenID string) Param {
	return paramFunc(func(log *logging.EventLogV2) {
		log.TokenId = (*logging.TokenId)(&tokenID)
	})
}

func OrgID(orgID string) Param {
	return paramFunc(func(log *logging.EventLogV2) {
		log.OrgId = (*logging.OrganizationId)(&orgID)
	})
}

func UnsafeParam(key string, value any) Param {
	return paramFunc(func(log *logging.EventLogV2) {
		wloginternal.SetMapParam(&log.UnsafeParams, key, value)
	})
}

func UnsafeParams(unsafe map[string]any) Param {
	return paramFunc(func(log *logging.EventLogV2) {
		wloginternal.SetMapParams(&log.UnsafeParams, unsafe)
	})
}

func TraceID(traceID string) Param {
	return paramFunc(func(log *logging.EventLogV2) {
		log.TraceId = (*logging.TraceId)(&traceID)
	})
}
