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
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/extractor"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logentry"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logtypes"
)

type defaultLogger struct {
	logger       logentry.Logger
	idsExtractor extractor.IDsFromRequest

	pathParamPerms   ParamPerms
	queryParamPerms  ParamPerms
	headerParamPerms ParamPerms
}

func (l *defaultLogger) Request(r Request) {
	log := poolRequestLogV2.Get()
	populateLogObject(log, r, l.idsExtractor, l.pathParamPerms, l.queryParamPerms, l.headerParamPerms)
	l.logger.LogObject(log)
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
	log := new(logtypes.RequestLogV2)
	populateLogObject(log, r, idsExtractor, pathParamPerms, queryParamPerms, headerParamPerms)
	return []wlog.Param{logentry.Param(log)}
}

// populateLogObject pre-populates a RequestLogV2 from the request and permissions: type, time, method, protocol,
// path, status, request/response sizes, duration, params/unsafeParams, and any IDs extractable from the request.
func populateLogObject(log *logtypes.RequestLogV2, r Request, idsExtractor extractor.IDsFromRequest, pathParamPerms, queryParamPerms, headerParamPerms ParamPerms) {
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
	log.RequestSize = r.Request.ContentLength
	log.ResponseSize = r.ResponseSize
	log.Duration = r.Duration.Nanoseconds() / 1000
	log.Method = r.Request.Method
	log.Uid = idsMap[wlog.UIDKey]
	log.Sid = idsMap[wlog.SIDKey]
	log.TokenId = idsMap[wlog.TokenIDKey]
	log.OrgId = idsMap[wlog.OrgIDKey]
	log.TraceId = idsMap[traceIDKey]
}

// populateRequestParams parses the path, header and query parameters. If any of the parameters are in a respective
// "forbidden" list, they are not logged at all. Otherwise, if a parameter is whitelisted it is added to safeParams and
// is added to unsafeParams otherwise. If a single key has multiple values, the value for that key in the log
// will be a slice that contains all of the values for the key.
func populateRequestParams(log *logtypes.RequestLogV2, r Request, pathParamPerms, queryParamPerms, headerParamPerms ParamPerms) {
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

func processKeyValPair(k, v string, safeDst, unsafeDst *logtypes.AnyMap, basePerms, reqPerms ParamPerms) {
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

func addAsMultiMap(k, v string, dst *logtypes.AnyMap) {
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
