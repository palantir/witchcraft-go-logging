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

package svc1log

import (
	"maps"
	"path"
	"runtime"
	"strconv"
	"time"

	"github.com/palantir/pkg/datetime"
	werror "github.com/palantir/witchcraft-go-error"
	"github.com/palantir/witchcraft-go-logging/internal/gopath"
	"github.com/palantir/witchcraft-go-logging/wapi/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
	wparams "github.com/palantir/witchcraft-go-params"
)

const (
	TypeValue = "service.1"
)

type Param = wloginternal.Param[logging.ServiceLogV1]

type paramFunc = wloginternal.ParamFunc[logging.ServiceLogV1]

func Type() Param {
	return paramFunc(func(l *logging.ServiceLogV1) {
		l.Type = TypeValue
	})
}

func Level(level wlog.LogLevel) Param {
	return withLevel(level.ToLoggingType())
}

func withLevel(level logging.LogLevel) Param {
	return paramFunc(func(l *logging.ServiceLogV1) {
		l.Level = level
	})
}

func Time(time time.Time) Param {
	return paramFunc(func(l *logging.ServiceLogV1) {
		l.Time = datetime.DateTime(time)
	})
}

func TimeNow() Param {
	// Defer execution of time.Now() until the log is actually written
	return paramFunc(func(l *logging.ServiceLogV1) {
		l.Time = datetime.DateTime(time.Now())
	})
}

func Origin(origin string) Param {
	return paramFunc(func(l *logging.ServiceLogV1) {
		l.Origin = &origin
	})
}

// OriginFromInitLine sets the "origin" field to be the filename and line of the location at which this function is
// called.
func OriginFromInitLine() Param {
	origin := ""
	if file, line, ok := initLineCaller(1); ok {
		origin = file + ":" + strconv.Itoa(line)
	}
	return Origin(origin)
}

// OriginFromInitPkg sets the "origin" field to be the package path of the location at which this function is called.
// The skipPkg parameter determines the level of the parent package that should be used. For example, if this function
// is called in a file with the package path "github.com/palantir/witchcraft-go-logging/wlog", then with skipPkg=0 the
// origin would be "github.com/palantir/witchcraft-go-logging/wlog", while with skipPkg=1 the origin would be
// "github.com/palantir/witchcraft-go-logging".
func OriginFromInitPkg(skipPkg int) Param {
	return Origin(CallerPkg(1, skipPkg))
}

// OriginFromCallLine sets the "origin" field to be the filename and line of the location at which the logger invocation
// is performed.
//
// Note that, when this parameter is used, every log invocation will perform a "runtime.Caller" call, which may not be
// suitable for performance-critical scenarios.
//
// Note that this parameter is tied to the implementation details of the logger implementations defined in the svc1log
// package (it hard-codes assumptions relating to the number of call stacks that must be skipped to reach the log site).
// Using this parameter with an svc1log.Logger implementation not defined in the svc1log package may result in incorrect
// output. If wrapping the default implementation of svc1log.Logger, OriginFromCallLineWithSkip allows for trimming
// additional stack frames.
func OriginFromCallLine() Param {
	return OriginFromCallLineWithSkip(0)
}

const defaultOriginFromCallLineStackSkip = 8

// OriginFromCallLineWithSkip is like OriginFromCallLine but allows for configuring additional skipped stack frames.
// This allows for libraries wrapping loggers to hide their implementation frames from the caller.
func OriginFromCallLineWithSkip(skipFrames int) Param {
	return paramFunc(func(l *logging.ServiceLogV1) {
		origin := ""
		if file, line, ok := initLineCaller(defaultOriginFromCallLineStackSkip + skipFrames); ok {
			origin = file + ":" + strconv.Itoa(line)
		}
		l.Origin = &origin
	})
}

// CallerPkg returns a package path based on the location at which this function is called and the parameters given to
// the function. This can be used in conjunction with the "Origin" param to set the origin field programmatically.
//
// The parentCaller parameter specifies the number of "parents" to go back in the call stack, while the parentPkg
// parameter determines the level of the parent package that should be used. For example, if this function is called in
// a file with the package path "github.com/palantir/witchcraft-go-logging/wlog" and that function is called from a file
// with the package path "github.com/palantir/project/helper", then with parentCaller=0 and parentPkg=0 the returned
// value would be "github.com/palantir/witchcraft-go-logging/wlog", while with parentCaller=1 and parentPkg=1 the value
// would be "github.com/palantir/project" (parentCaller=1 sets the package to "github.com/palantir/project/helper" and
// parentPkg=1 causes the package to become "github.com/palantir/project").
func CallerPkg(parentCaller, parentPkg int) string {
	origin := ""
	if file, _, ok := initLineCaller(1 + parentCaller); ok {
		origin = path.Dir(file)
		for i := 0; i < parentPkg; i++ {
			origin = path.Dir(origin)
		}
	}
	return origin
}

func initLineCaller(skip int) (string, int, bool) {
	// the 1 skips the current "initLineCaller" function
	_, file, line, ok := runtime.Caller(1 + skip)
	if ok {
		file = gopath.TrimPrefix(file)
	}
	return file, line, ok
}

func Thread(thread string) Param {
	return paramFunc(func(l *logging.ServiceLogV1) {
		l.Thread = &thread
	})
}

func Message(message string) Param {
	return paramFunc(func(l *logging.ServiceLogV1) {
		l.Message = message
	})
}

func Params(params wparams.ParamStorer) Param {
	return paramFunc(func(l *logging.ServiceLogV1) {
		if params != nil {
			wloginternal.ApplyParams(l, SafeParams(params.SafeParams()), UnsafeParams(params.UnsafeParams()))
		}
	})
}

func SafeParams(params map[string]any) Param {
	return paramFunc(func(l *logging.ServiceLogV1) {
		if l.Params == nil {
			l.Params = maps.Clone(params)
		} else {
			for k, v := range params {
				l.Params[k] = v
			}
		}
	})
}

func SafeParam(key string, value any) Param {
	return paramFunc(func(l *logging.ServiceLogV1) {
		if l.Params == nil {
			l.Params = map[string]any{key: value}
		} else {
			l.Params[key] = value
		}
	})
}

func UID(uid string) Param {
	return paramFunc(func(l *logging.ServiceLogV1) {
		l.Uid = (*logging.UserId)(&uid)
	})
}

func SID(sid string) Param {
	return paramFunc(func(l *logging.ServiceLogV1) {
		l.Sid = (*logging.SessionId)(&sid)
	})
}

func TokenID(tokenId string) Param {
	return paramFunc(func(l *logging.ServiceLogV1) {
		l.TokenId = (*logging.TokenId)(&tokenId)
	})
}

func OrgID(orgId string) Param {
	return paramFunc(func(l *logging.ServiceLogV1) {
		l.OrgId = (*logging.OrgId)(&orgId)
	})
}

func TraceID(traceId string) Param {
	return paramFunc(func(l *logging.ServiceLogV1) {
		l.TraceId = (*logging.TraceId)(&traceId)
	})
}

func Stacktrace(err error) Param {
	return paramFunc(func(l *logging.ServiceLogV1) {
		if err != nil {
			stacktrace := werror.GenerateErrorString(err, false)
			l.Stacktrace = &stacktrace

			// add all safe and unsafe parameters stored in error
			safeParams, unsafeParams := werror.ParamsFromError(err)
			wloginternal.ApplyParams(l, SafeParams(safeParams), UnsafeParams(unsafeParams))
		}
	})
}

func UnsafeParams(unsafeParams map[string]any) Param {
	return paramFunc(func(l *logging.ServiceLogV1) {
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
	return paramFunc(func(l *logging.ServiceLogV1) {
		if l.UnsafeParams == nil {
			l.UnsafeParams = map[string]any{key: value}
		} else {
			l.UnsafeParams[key] = value
		}
	})
}

func Tags(tags map[string]string) Param {
	return paramFunc(func(l *logging.ServiceLogV1) {
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
	return paramFunc(func(l *logging.ServiceLogV1) {
		if l.Tags == nil {
			l.Tags = map[string]string{key: value}
		} else {
			l.Tags[key] = value
		}
	})
}
