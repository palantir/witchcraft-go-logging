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
	"github.com/palantir/witchcraft-go-logging/wapi/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/extractor"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
)

var objectPool = wloginternal.NewPool((*logging.RequestLogV2).Reset)

type defaultLogger struct {
	logger       wlog.Logger[logging.RequestLogV2]
	idsExtractor extractor.IDsFromRequest

	pathParamPerms   ParamPerms
	queryParamPerms  ParamPerms
	headerParamPerms ParamPerms
}

func (l *defaultLogger) Request(r Request) {
	wloginternal.LogObject(l.logger, objectPool, ToParams(r, l.idsExtractor, l.pathParamPerms, l.queryParamPerms, l.headerParamPerms))
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

func ToParams(r Request, idsExtractor extractor.IDsFromRequest, pathParamPerms, queryParamPerms, headerParamPerms ParamPerms) wloginternal.Param[logging.RequestLogV2] {
	safeParams, unsafeParams := parseRequestParams(r, pathParamPerms, queryParamPerms, headerParamPerms)

	reqPath := r.Request.URL.Path
	if r.RouteInfo.Template != "" {
		reqPath = r.RouteInfo.Template
	}

	// extract IDs from request
	idsMap := idsExtractor.ExtractIDs(r.Request)

	return wloginternal.ParamFunc[logging.RequestLogV2](func(l *logging.RequestLogV2) {
		l.Type = TypeValue

		l.Time = datetime.DateTime(time.Now())

		if method := r.Request.Method; method != "" {
			l.Method = &method
		}

		l.Protocol = r.Request.Proto

		l.Path = reqPath

		wloginternal.SetMapParams(&l.Params, safeParams)

		l.Status = r.ResponseStatus

		l.RequestSize = safelong.SafeLong(r.Request.ContentLength)

		l.ResponseSize = safelong.SafeLong(r.ResponseSize)

		l.Duration = safelong.SafeLong(r.Duration.Microseconds())

		if uid := idsMap[extractor.UIDKey]; uid != "" {
			l.Uid = (*logging.UserId)(&uid)
		}

		if sid := idsMap[extractor.SIDKey]; sid != "" {
			l.Sid = (*logging.SessionId)(&sid)
		}

		if tokenID := idsMap[extractor.TokenIDKey]; tokenID != "" {
			l.TokenId = (*logging.TokenId)(&tokenID)
		}

		if orgID := idsMap[extractor.OrgIDKey]; orgID != "" {
			l.OrgId = (*logging.OrgId)(&orgID)
		}

		if traceID := idsMap[extractor.TraceIDKey]; traceID != "" {
			l.TraceId = (*logging.TraceId)(&traceID)
		}

		wloginternal.SetMapParams(&l.UnsafeParams, unsafeParams)
	})
}

// parseRequestParams parses the path, header and query parameters. If any of the parameters are in a respective
// "forbidden" list, they are not logged at all. Otherwise, if a parameter is whitelisted it is added to safeParams and
// is added to unsafeParams otherwise. If a single key has multiple values, the value for that key in the returned field
// will be a slice that contains all of the values for the key.
func parseRequestParams(r Request, pathParamPerms, queryParamPerms, headerParamPerms ParamPerms) (safeParams map[string]any, unsafeParams map[string]any) {
	safeMap := make(map[string]interface{})
	unsafeMap := make(map[string]interface{})

	for pathParamKey, pathParamVal := range r.RouteInfo.PathParams {
		processKeyValPair(pathParamKey, pathParamVal, safeMap, unsafeMap, pathParamPerms, r.PathParamPerms)
	}
	for k, valSlice := range r.Request.URL.Query() {
		for _, v := range valSlice {
			processKeyValPair(k, v, safeMap, unsafeMap, queryParamPerms, r.QueryParamPerms)
		}
	}
	for k := range r.Request.Header {
		processKeyValPair(k, r.Request.Header.Get(k), safeMap, unsafeMap, headerParamPerms, r.HeaderParamPerms)
	}
	return safeMap, unsafeMap
}

func processKeyValPair(k, v string, safeDst, unsafeDst map[string]interface{}, basePerms, reqPerms ParamPerms) {
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

func addAsMultiMap(k, v string, m map[string]interface{}) {
	currVal, exists := m[k]
	if !exists {
		m[k] = v
		return
	}

	var newVal interface{}
	// value for key already exists in destination. If destination is a slice, append to it.
	// Otherwise, create a new slice, add existing value to it and append new value.
	switch t := currVal.(type) {
	default:
		newVal = append([]interface{}{t}, v)
	case []interface{}:
		newVal = append(t, v)
	}
	m[k] = newVal
}
