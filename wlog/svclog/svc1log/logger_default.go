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
	"github.com/palantir/witchcraft-go-logging/wapi/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
)

var objectPool = wloginternal.NewPool((*logging.ServiceLogV1).Reset)

type defaultLogger struct {
	logger wlog.Logger[logging.ServiceLogV1]
	level  *wlog.AtomicLogLevel
}

func (l *defaultLogger) Debug(msg string, params ...Param) {
	l.log(wlog.DebugLevel, msg, params...)
}

func (l *defaultLogger) Info(msg string, params ...Param) {
	l.log(wlog.InfoLevel, msg, params...)
}

func (l *defaultLogger) Warn(msg string, params ...Param) {
	l.log(wlog.WarnLevel, msg, params...)
}

func (l *defaultLogger) Error(msg string, params ...Param) {
	l.log(wlog.ErrorLevel, msg, params...)
}

func (l *defaultLogger) log(level wlog.LogLevel, msg string, params ...Param) {
	if l.Enabled(level) {
		wloginternal.LogObject(l.logger, objectPool, defaultParam(level, msg), params...)
	}
}

func (l *defaultLogger) SetLevel(level wlog.LogLevel) {
	l.level.SetLevel(level)
}

func (l *defaultLogger) Enabled(level wlog.LogLevel) bool {
	return l.level == nil || l.level.Enabled(level)
}
