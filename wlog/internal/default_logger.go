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

package wloginternal

import (
	"fmt"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"runtime"
	"strings"
	"time"
)

// WarnLoggerOutput returns the logger output for a default warning logger for a given logger type. The output includes
// the location at which the call was made, with the "skip" parameter determining how far back in the call stack to go
// for the location (for example, skip=0 specifies the line in this function, skip=1 specifies the line that called this
// function, etc.).
//
// This function is defined in an internal package because each logger type needs to define its own warning logger type
// but the format/content of the output should be consistent across them.
func WarnLoggerOutput(loggerType, output string, skip int) string {
	pc, fn, line, _ := runtime.Caller(skip)
	return fmt.Sprintf("[WARNING] %s[%s:%d]: usage of %s.Logger from FromContext that did not have that logger set: %s", runtime.FuncForPC(pc).Name(), fn, line, loggerType, strings.TrimSuffix(output, "\n"))
}

func ToServiceParams(msg string, level wlog.Param, inParams []svc1log.Param) []wlog.Param {
	outParams := make([]wlog.Param, len(defaultServiceParams)+2+len(inParams))
	copy(outParams, defaultServiceParams)
	outParams[len(defaultServiceParams)] = level
	outParams[len(defaultServiceParams)+1] = wlog.NewParam(func(entry wlog.LogEntry) {
		entry.StringValue(svc1log.MessageKey, msg)
	})
	for idx := range inParams {
		outParams[len(defaultServiceParams)+2+idx] = wlog.NewParam(func(entry wlog.LogEntry) {
			svc1log.ApplyParam(inParams[idx], entry)
		})
	}
	return outParams
}

var defaultServiceParams = []wlog.Param{
	wlog.NewParam(func(entry wlog.LogEntry) {
		entry.StringValue(wlog.TypeKey, svc1log.TypeValue)
		entry.StringValue(wlog.TimeKey, time.Now().Format(time.RFC3339Nano))
	}),
}
