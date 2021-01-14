// Copyright (c) 2021 Palantir Technologies. All rights reserved.
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

package svc1logr

import (
	"fmt"
	"path/filepath"

	"github.com/go-logr/logr"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

type logger struct {
	origin    string
	svcLogger svc1log.Logger
}

// New returns a go-logr interface implementation that uses svc1log internally.
func New(svcLogger svc1log.Logger, origin string) logr.Logger {
	svcLogger = svc1log.WithParams(svcLogger, svc1log.Origin(origin))
	return &logger{
		origin:    origin,
		svcLogger: svcLogger,
	}
}

func (l *logger) Info(msg string, keysAndValues ...interface{}) {
	l.svcLogger.Info(msg, toSafeParams(l.svcLogger, keysAndValues))
}

func (l *logger) Enabled() bool {
	return true
}

func (l *logger) Error(err error, msg string, keysAndValues ...interface{}) {
	l.svcLogger.Error(msg, svc1log.Stacktrace(err), toSafeParams(l.svcLogger, keysAndValues))
}

func (l *logger) V(level int) logr.InfoLogger {
	return New(l.svcLogger, l.origin)
}

func (l *logger) WithValues(keysAndValues ...interface{}) logr.Logger {
	svcLogger := svc1log.WithParams(l.svcLogger, toSafeParams(l.svcLogger, keysAndValues))
	return New(svcLogger, l.origin)
}

func (l *logger) WithName(name string) logr.Logger {
	return New(l.svcLogger, filepath.Join(l.origin, name))
}

func toSafeParams(svcLogger svc1log.Logger, keysAndValues []interface{}) svc1log.Param {
	if len(keysAndValues)%2 == 1 {
		svcLogger.Error("KeysAndValues pair slice has an odd number of arguments; ignoring all",
			svc1log.SafeParam("keysAndValuesLen", len(keysAndValues)))
		return svc1log.SafeParams(map[string]interface{}{})
	}

	params := make(map[string]interface{}, len(keysAndValues)/2)
	for i := 0; i < len(keysAndValues); i = i + 2 {
		key, ok := keysAndValues[i].(string)
		if !ok {
			svcLogger.Error("Key type is not string",
				svc1log.SafeParam("actualType", fmt.Sprintf("%T", keysAndValues[i])),
				svc1log.SafeParam("key", key))
			continue
		}
		params[key] = keysAndValues[i+1]
	}
	return svc1log.SafeParams(params)
}
