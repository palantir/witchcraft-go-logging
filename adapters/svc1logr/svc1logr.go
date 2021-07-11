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
	"path"

	"github.com/go-logr/logr"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
)

type svc1logr struct {
	origin  string
	logger  svc1log.Logger
	level   int
	enabled bool
}

// New provides an implementation of [logr.Logger](https://pkg.go.dev/github.com/go-logr/logr#Logger) which delegates to
// a witchcraft svc1log logger. Use this package for software which expects a logr.Logger but where the calling code
// has access to an svc1log.Logger or wants the output to be Witchcraft-compatible.  Info() is logged as INFO and
// Error() as ERROR. All values are converted to safe parameters.
func New(logger svc1log.Logger, origin string, level int) logr.Logger {
	logger = svc1log.WithParams(logger, svc1log.Origin(origin))
	return &svc1logr{
		origin:  origin,
		logger:  logger,
		level:   level,
		enabled: true,
	}
}

func (s *svc1logr) Info(msg string, keysAndValues ...interface{}) {
	if s.Enabled() {
		s.logger.Info(msg, toSafeParams(s.logger, keysAndValues))
	}
}

func (s *svc1logr) Enabled() bool {
	return s.enabled
}

func (s *svc1logr) Error(err error, msg string, keysAndValues ...interface{}) {
	s.logger.Error(msg, svc1log.Stacktrace(err), toSafeParams(s.logger, keysAndValues))
}

func (s *svc1logr) V(level int) logr.Logger {
	enabled := true
	if level > s.level {
		enabled = false
	}
	return &svc1logr{
		origin:  s.origin,
		logger:  s.logger,
		level:   s.level,
		enabled: enabled,
	}
}

func (s *svc1logr) WithValues(keysAndValues ...interface{}) logr.Logger {
	logger := svc1log.WithParams(s.logger, toSafeParams(s.logger, keysAndValues))
	return New(logger, s.origin, s.level)
}

func (s *svc1logr) WithName(name string) logr.Logger {
	return New(s.logger, path.Join(s.origin, name), s.level)
}

func toSafeParams(logger svc1log.Logger, keysAndValues []interface{}) svc1log.Param {
	if len(keysAndValues)%2 == 1 {
		logger.Error("Logging error: keysAndValues pair slice has an odd number of arguments; ignoring all",
			svc1log.SafeParam("keysAndValuesLen", len(keysAndValues)))
		return svc1log.SafeParams(map[string]interface{}{})
	}

	params := make(map[string]interface{}, len(keysAndValues)/2)
	for i := 0; i < len(keysAndValues); i = i + 2 {
		key, ok := keysAndValues[i].(string)
		if !ok {
			logger.Error("Logging error: Key type is not string",
				svc1log.SafeParam("actualType", fmt.Sprintf("%T", keysAndValues[i])),
				svc1log.SafeParam("key", key))
			continue
		}
		params[key] = keysAndValues[i+1]
	}
	return svc1log.SafeParams(params)
}
