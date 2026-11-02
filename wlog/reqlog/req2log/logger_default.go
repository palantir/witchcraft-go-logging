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

package req2log

import (
	"strings"
	"time"

	"github.com/palantir/pkg/datetime"
	"github.com/palantir/pkg/safelong"
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/extractor"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
)

// conjureLogger consumes borrowed log objects synchronously. Implementations must copy
// any object or nested data they retain after the method returns.
type conjureLogger interface {
	LogRequestLogV2(*logging.RequestLogV2)
}

type defaultLogger struct {
	logger       conjureLogger
	idsExtractor extractor.IDsFromRequest

	pathParamPerms   ParamPerms
	queryParamPerms  ParamPerms
	headerParamPerms ParamPerms
}

func (l *defaultLogger) Request(r Request) {
	log := poolRequestLogV2.Get()
	populateLogObject(log, r, l.idsExtractor, l.pathParamPerms, l.queryParamPerms, l.headerParamPerms)
	l.logger.LogRequestLogV2(log)
	poolRequestLogV2.Put(log)
}

func (l *defaultLogger) PathParamPerms() ParamPerms {
	return l.pathParamPerms
}

func (l *defaultLogger) QueryParamPerms() ParamPerms {
	return l.queryParamPerms
}

func (l *defaultLogger) HeaderParamPerms() ParamPerms {
	return l.headerParamPerms
}

// ToParams builds parameters backed by a fresh object that can outlive the log call.
func ToParams(r Request, idsExtractor extractor.IDsFromRequest, pathParamPerms, queryParamPerms, headerParamPerms ParamPerms) []wlog.Param {
	log := new(logging.RequestLogV2)
	populateLogObject(log, r, idsExtractor, pathParamPerms, queryParamPerms, headerParamPerms)
	return []wlog.Param{objectToParam(log)}
}

func ToObject(r Request, idsExtractor extractor.IDsFromRequest, pathParamPerms, queryParamPerms, headerParamPerms ParamPerms) logging.RequestLogV2 {
	var log logging.RequestLogV2
	populateLogObject(&log, r, idsExtractor, pathParamPerms, queryParamPerms, headerParamPerms)
	return log
}

// populateLogObject pre-populates a RequestLogV2 from the request and permissions: type, time, method, protocol,
// path, status, request/response sizes, duration, params/unsafeParams, and any IDs extractable from the request.
func populateLogObject(log *logging.RequestLogV2, r Request, idsExtractor extractor.IDsFromRequest, pathParamPerms, queryParamPerms, headerParamPerms ParamPerms) {
	populateRequestParams(log, r, pathParamPerms, queryParamPerms, headerParamPerms)

	reqPath := r.Request.URL.Path
	if r.RouteInfo.Template != "" {
		reqPath = r.RouteInfo.Template
	}

	idsMap := idsExtractor.ExtractIDs(r.Request)

	log.Type = TypeValue
	log.Time = datetime.DateTime(time.Now())
	log.Protocol = r.Request.Proto
	log.Path = reqPath
	log.Status = r.ResponseStatus
	log.RequestSize = safelong.SafeLong(r.Request.ContentLength)
	log.ResponseSize = safelong.SafeLong(r.ResponseSize)
	log.Duration = safelong.SafeLong(r.Duration.Nanoseconds() / 1000)

	if r.Request.Method != "" {
		log.Method = &r.Request.Method
	}
	setOptionalID(&log.Uid, idsMap, wlog.UIDKey)
	setOptionalID(&log.Sid, idsMap, wlog.SIDKey)
	setOptionalID(&log.TokenId, idsMap, wlog.TokenIDKey)
	setOptionalID(&log.OrgId, idsMap, wlog.OrgIDKey)
	setOptionalID(&log.TraceId, idsMap, traceIDKey)
}

// objectToParam emits all set fields of a RequestLogV2 to a wlog.LogEntry.
func objectToParam(log *logging.RequestLogV2) wlog.Param {
	return wlog.NewParam(func(entry wlog.LogEntry) {
		entry.StringValue(wlog.TypeKey, log.Type)
		entry.StringValue(wlog.TimeKey, log.Time.String())
		if log.Method != nil {
			entry.OptionalStringValue(methodKey, *log.Method)
		}
		entry.StringValue(protocolKey, log.Protocol)
		entry.StringValue(pathKey, log.Path)
		if len(log.Params) > 0 {
			entry.AnyMapValue(paramsKey, log.Params)
		}
		entry.IntValue(statusKey, int32(log.Status))
		entry.SafeLongValue(requestSizeKey, int64(log.RequestSize))
		entry.SafeLongValue(responseSizeKey, int64(log.ResponseSize))
		entry.SafeLongValue(durationKey, int64(log.Duration))
		if log.Uid != nil {
			entry.StringValue(wlog.UIDKey, string(*log.Uid))
		}
		if log.Sid != nil {
			entry.StringValue(wlog.SIDKey, string(*log.Sid))
		}
		if log.TokenId != nil {
			entry.StringValue(wlog.TokenIDKey, string(*log.TokenId))
		}
		if log.OrgId != nil {
			entry.StringValue(wlog.OrgIDKey, string(*log.OrgId))
		}
		if log.TraceId != nil {
			entry.StringValue(traceIDKey, string(*log.TraceId))
		}
		if len(log.UnsafeParams) > 0 {
			entry.AnyMapValue(wlog.UnsafeParamsKey, log.UnsafeParams)
		}
	})
}

type entryLogger struct {
	logger wlog.Logger
}

func (l entryLogger) LogRequestLogV2(log *logging.RequestLogV2) {
	l.logger.Log(objectToParam(log))
}

func setOptionalID[T ~string](dst **T, idsMap map[string]string, key string) {
	if v := idsMap[key]; v != "" {
		*dst = new(T(v))
	}
}

// populateRequestParams parses the path, header and query parameters. If any of the parameters are in a respective
// "forbidden" list, they are not logged at all. Otherwise, if a parameter is whitelisted it is added to safeParams and
// is added to unsafeParams otherwise. If a single key has multiple values, the value for that key in the log
// will be a slice that contains all of the values for the key.
func populateRequestParams(log *logging.RequestLogV2, r Request, pathParamPerms, queryParamPerms, headerParamPerms ParamPerms) {
	for pathParamKey, pathParamVal := range r.RouteInfo.PathParams {
		processKeyValPair(pathParamKey, pathParamVal, &log.Params, &log.UnsafeParams, pathParamPerms, r.PathParamPerms)
	}
	for k, valSlice := range r.Request.URL.Query() {
		for _, v := range valSlice {
			processKeyValPair(k, v, &log.Params, &log.UnsafeParams, queryParamPerms, r.QueryParamPerms)
		}
	}
	for k := range r.Request.Header {
		processKeyValPair(k, r.Request.Header.Get(k), &log.Params, &log.UnsafeParams, headerParamPerms, r.HeaderParamPerms)
	}
}

func processKeyValPair(k, v string, safeDst, unsafeDst *map[string]any, basePerms, reqPerms ParamPerms) {
	// lowercase keys are used for lookups. Convert once here to avoid multiple unnecessary allocations.
	// Note that, if a key is added to an output map, the original unconverted key should be added.
	lowerK := strings.ToLower(k)

	// NOTE: logic below is duplicated twice for basePerms and reqPerms instead of making perms a var-arg param and
	// iterating over it in the interest of performance (avoids slice allocation on every call).
	if basePerms != nil && basePerms.Forbidden(lowerK) {
		// key is forbidden/blacklisted: do not record at all
		return
	}
	if reqPerms != nil && reqPerms.Forbidden(lowerK) {
		// key is forbidden/blacklisted: do not record at all
		return
	}

	// iterate over whitelist in separate loop after all of the forbidden keys are processed because forbidden
	// takes precedence over whitelist
	if basePerms != nil && basePerms.Safe(lowerK) {
		// key is whitelisted and not forbidden: add to safe
		addAsMultiMap(k, v, safeDst)
		return
	}
	if reqPerms != nil && reqPerms.Safe(lowerK) {
		// key is whitelisted and not forbidden: add to safe
		addAsMultiMap(k, v, safeDst)
		return
	}

	// not in any forbidden list or whitelist: add to unsafe
	addAsMultiMap(k, v, unsafeDst)
}

func addAsMultiMap(k, v string, dst *map[string]any) {
	currVal, exists := (*dst)[k]
	if !exists {
		wloginternal.SetMapParam(dst, k, any(v))
		return
	}

	var newVal any
	// value for key already exists in destination. If destination is a slice, append to it.
	// Otherwise, create a new slice, add existing value to it and append new value.
	switch t := currVal.(type) {
	default:
		newVal = append([]any{t}, v)
	case []any:
		newVal = append(t, v)
	}
	(*dst)[k] = newVal
}
