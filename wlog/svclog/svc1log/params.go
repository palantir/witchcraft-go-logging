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
	"slices"
	"strconv"
	"strings"
	"sync"

	werror "github.com/palantir/witchcraft-go-error"
	"github.com/palantir/witchcraft-go-logging/internal/gopath"
	"github.com/palantir/witchcraft-go-logging/wlog"
	wparams "github.com/palantir/witchcraft-go-params"
)

const (
	TypeValue = "service.1"

	LevelKey        = "level"
	LevelDebugValue = "DEBUG"
	LevelInfoValue  = "INFO"
	LevelWarnValue  = "WARN"
	LevelErrorValue = "ERROR"

	OriginKey     = "origin"
	ThreadKey     = "thread"
	MessageKey    = "message"
	ParamsKey     = "params"
	StacktraceKey = "stacktrace"
	TagsKey       = "tags"
)

type Param interface {
	apply(entry wlog.LogEntry)
}

func ApplyParam(p Param, entry wlog.LogEntry) {
	if p == nil {
		return
	}
	p.apply(entry)
}

type paramFunc func(entry wlog.LogEntry)

func (f paramFunc) apply(entry wlog.LogEntry) {
	f(entry)
}

// Origin sets the "origin" field to be the provided value if it is non-empty.
func Origin(origin string) Param {
	return paramFunc(func(logger wlog.LogEntry) {
		logger.OptionalStringValue(OriginKey, origin)
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
		for range parentPkg {
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
// is performed: the innermost stack frame outside the wlog packages.
//
// Note that, when this parameter is used, every log invocation walks the stack, which may not be suitable for
// performance-critical scenarios.
func OriginFromCallLine() Param {
	return OriginFromCallLineSkippingPackages()
}

// OriginFromCallLineSkippingPackages is like OriginFromCallLine but also skips frames in the provided packages and their
// subpackages. This allows libraries that wrap loggers to hide their frames from the origin.
func OriginFromCallLineSkippingPackages(pkgs ...string) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.OptionalStringValue(OriginKey, originOutside(pkgs))
	})
}

const wlogPackage = "github.com/palantir/witchcraft-go-logging/wlog"

// originOutside returns the origin of the innermost caller outside the wlog packages, the provider modules such as
// wlog-zap that apply params while they build an entry, and pkgs.
func originOutside(pkgs []string) string {
	// runtime.Callers unwinds as many frames as the buffer holds, so start small and grow only if needed.
	var buf [16]uintptr
	pcs := buf[:]
	for {
		n := runtime.Callers(3, pcs)
		for _, pc := range pcs[:n] {
			for _, frame := range framesAt(pc) {
				if !frame.inWlog && !slices.ContainsFunc(pkgs, func(p string) bool { return inPackage(frame.pkg, p) }) {
					return frame.origin
				}
			}
		}
		if n < len(pcs) {
			return ""
		}
		pcs = make([]uintptr, 2*len(pcs))
	}
}

type pcFrame struct {
	pkg    string
	origin string
	inWlog bool
}

// pcFrames caches framesAt, so that a log call resolves only call sites it has not seen before.
var pcFrames sync.Map

// framesAt returns the frames at a return PC from runtime.Callers, innermost first. A PC has more than one frame when
// calls are inlined into it.
func framesAt(pc uintptr) []pcFrame {
	if frames, ok := pcFrames.Load(pc); ok {
		return frames.([]pcFrame)
	}
	var out []pcFrame
	frames := runtime.CallersFrames([]uintptr{pc})
	for {
		frame, more := frames.Next()
		pkg := funcPackage(frame.Function)
		// Test packages are callers even when they live in a wlog directory.
		inWlog := (inPackage(pkg, wlogPackage) || strings.HasPrefix(pkg, wlogPackage+"-")) && !strings.HasSuffix(pkg, "_test")
		out = append(out, pcFrame{pkg: pkg, origin: gopath.TrimPrefix(frame.File) + ":" + strconv.Itoa(frame.Line), inWlog: inWlog})
		if !more {
			break
		}
	}
	stored, _ := pcFrames.LoadOrStore(pc, out)
	return stored.([]pcFrame)
}

// funcPackage returns the package path of a function name reported by the runtime, such as "a/b.(*T).M[...]".
func funcPackage(fn string) string {
	fn, _, _ = strings.Cut(fn, "[")
	lastSlash := strings.LastIndexByte(fn, '/')
	if dot := strings.IndexByte(fn[lastSlash+1:], '.'); dot >= 0 {
		return fn[:lastSlash+1+dot]
	}
	return fn
}

func inPackage(pkg, parent string) bool {
	rest, ok := strings.CutPrefix(pkg, parent)
	return ok && (rest == "" || rest[0] == '/')
}

const defaultOriginFromCallLineStackSkip = 8

// OriginFromCallLineWithSkip sets the "origin" field to the location that is skipFrames frames above the svc1log
// logger call. Unlike OriginFromCallLine, it counts frames, so skipFrames must include every frame that wraps the
// logger, including frames added by svc1log itself such as the wrapper that FromContext adds for context params.
// Prefer OriginFromCallLineSkippingPackages, which does not depend on the number of wrapping frames.
func OriginFromCallLineWithSkip(skipFrames int) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.OptionalStringValue(OriginKey, originAbove(defaultOriginFromCallLineStackSkip+skipFrames))
	})
}

// originAbove returns the origin of the frame that is depth frames above its caller, counting inlined frames.
func originAbove(depth int) string {
	var buf [32]uintptr
	pcs := buf[:min(depth+1, len(buf))]
	if depth >= len(buf) {
		pcs = make([]uintptr, depth+1)
	}
	for _, pc := range pcs[:runtime.Callers(2, pcs)] {
		frames := framesAt(pc)
		if depth < len(frames) {
			return frames[depth].origin
		}
		depth -= len(frames)
	}
	return ""
}

func initLineCaller(skip int) (string, int, bool) {
	// the 1 skips the current "initLineCaller" function
	_, file, line, ok := runtime.Caller(1 + skip)
	if ok {
		file = gopath.TrimPrefix(file)
	}
	return file, line, ok
}

func SafeParam(key string, value any) Param {
	return SafeParams(map[string]any{
		key: value,
	})
}

func SafeParams(safe map[string]any) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.AnyMapValue(ParamsKey, safe)
	})
}

func UID(uid string) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.OptionalStringValue(wlog.UIDKey, uid)
	})
}

func SID(sid string) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.OptionalStringValue(wlog.SIDKey, sid)
	})
}

func TokenID(tokenID string) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.OptionalStringValue(wlog.TokenIDKey, tokenID)
	})
}

func OrgID(orgID string) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.OptionalStringValue(wlog.OrgIDKey, orgID)
	})
}

func TraceID(traceID string) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.OptionalStringValue(wlog.TraceIDKey, traceID)
	})
}

func Stacktrace(err error) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		if err == nil {
			return
		}
		entry.StringValue(StacktraceKey, werror.GenerateErrorString(err, false))

		// add all safe and unsafe parameters stored in error
		safeParams, unsafeParams := werror.ParamsFromError(err)
		SafeParams(safeParams).apply(entry)
		UnsafeParams(unsafeParams).apply(entry)
	})
}

func UnsafeParam(key string, value any) Param {
	return UnsafeParams(map[string]any{
		key: value,
	})
}

func UnsafeParams(unsafe map[string]any) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.AnyMapValue(wlog.UnsafeParamsKey, unsafe)
	})
}

func Params(object wparams.ParamStorer) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		SafeParams(object.SafeParams()).apply(entry)
		UnsafeParams(object.UnsafeParams()).apply(entry)
	})
}

func Tag(key, value string) Param {
	return Tags(map[string]string{key: value})
}

func Tags(tags map[string]string) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.StringMapValue(TagsKey, tags)
	})
}
