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
	"path"
	"runtime"
	"strconv"

	werror "github.com/palantir/witchcraft-go-error"
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/internal/gopath"
	wparams "github.com/palantir/witchcraft-go-params"
)

type Param interface {
	apply(log *logging.ServiceLogV1)
}

type paramFunc func(log *logging.ServiceLogV1)

func (f paramFunc) apply(log *logging.ServiceLogV1) {
	f(log)
}

// Origin sets the "origin" field to be the provided value if it is non-empty.
func Origin(origin string) Param {
	return paramFunc(func(log *logging.ServiceLogV1) {
		log.Origin = &origin
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
	return Origin(CallerPkg(1, 0))
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
	return paramFunc(func(log *logging.ServiceLogV1) {
		if file, line, ok := initLineCaller(defaultOriginFromCallLineStackSkip + skipFrames); ok {
			origin := file + ":" + strconv.Itoa(line)
			log.Origin = &origin
		}
	})
}

func initLineCaller(skip int) (string, int, bool) {
	// the 1 skips the current "initLineCaller" function
	_, file, line, ok := runtime.Caller(1 + skip)
	if ok {
		file = gopath.TrimPrefix(file)
	}
	return file, line, ok
}

func SafeParam(key string, value interface{}) Param {
	return paramFunc(func(log *logging.ServiceLogV1) {
		if log.Params == nil {
			log.Params = map[string]interface{}{key: value}
		} else {
			log.Params[key] = value
		}
	})
}

func SafeParams(safe map[string]interface{}) Param {
	return paramFunc(func(log *logging.ServiceLogV1) {
		if log.Params == nil {
			log.Params = make(map[string]interface{}, len(safe))
		}
		for k, v := range safe {
			log.Params[k] = v
		}
	})
}

func UID(uid string) Param {
	return paramFunc(func(log *logging.ServiceLogV1) {
		log.Uid = (*logging.UserId)(&uid)
	})
}

func SID(sid string) Param {
	return paramFunc(func(log *logging.ServiceLogV1) {
		log.Sid = (*logging.SessionId)(&sid)
	})
}

func TokenID(tokenID string) Param {
	return paramFunc(func(log *logging.ServiceLogV1) {
		log.TokenId = (*logging.TokenId)(&tokenID)
	})
}

func TraceID(traceID string) Param {
	return paramFunc(func(log *logging.ServiceLogV1) {
		log.TraceId = (*logging.TraceId)(&traceID)
	})
}

func Stacktrace(err error) Param {
	return paramFunc(func(log *logging.ServiceLogV1) {
		if err == nil {
			return
		}

		stacktrace := werror.GenerateErrorString(err, false)
		log.Stacktrace = &stacktrace

		// add all safe and unsafe parameters stored in error
		safeParams, unsafeParams := werror.ParamsFromError(err)
		SafeParams(safeParams).apply(log)
		UnsafeParams(unsafeParams).apply(log)
	})
}

func UnsafeParam(key string, value interface{}) Param {
	return paramFunc(func(log *logging.ServiceLogV1) {
		if log.UnsafeParams == nil {
			log.UnsafeParams = map[string]interface{}{key: value}
		} else {
			log.UnsafeParams[key] = value
		}
	})
}

func UnsafeParams(unsafe map[string]interface{}) Param {
	return paramFunc(func(log *logging.ServiceLogV1) {
		if log.UnsafeParams == nil {
			log.UnsafeParams = make(map[string]interface{}, len(unsafe))
		}
		for k, v := range unsafe {
			log.UnsafeParams[k] = v
		}
	})
}

func Params(object wparams.ParamStorer) Param {
	return paramFunc(func(log *logging.ServiceLogV1) {
		if object != nil {
			SafeParams(object.SafeParams()).apply(log)
			UnsafeParams(object.UnsafeParams()).apply(log)
		}
	})
}

func Tag(key, value string) Param {
	return paramFunc(func(log *logging.ServiceLogV1) {
		if log.Tags == nil {
			log.Tags = map[string]string{key: value}
		} else {
			log.Tags[key] = value
		}
	})
}

func Tags(tags map[string]string) Param {
	return paramFunc(func(log *logging.ServiceLogV1) {
		if log.Tags == nil {
			log.Tags = make(map[string]string, len(tags))
		}
		for k, v := range tags {
			log.Tags[k] = v
		}
	})
}
