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
	"time"

	"github.com/palantir/pkg/datetime"
	"github.com/palantir/witchcraft-go-logging/wlog"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logentry"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logtypes"
)

type defaultLogger struct {
	logger   logentry.Logger
	delegate wlog.LeveledLogger
	level    wlog.LevelChecker
}

func (l *defaultLogger) Debug(msg string, params ...Param) {
	if l.Enabled(wlog.DebugLevel) {
		l.log(logtypes.LogLevelDebug, msg, params)
	}
}

func (l *defaultLogger) Info(msg string, params ...Param) {
	if l.Enabled(wlog.InfoLevel) {
		l.log(logtypes.LogLevelInfo, msg, params)
	}
}

func (l *defaultLogger) Warn(msg string, params ...Param) {
	if l.Enabled(wlog.WarnLevel) {
		l.log(logtypes.LogLevelWarn, msg, params)
	}
}

func (l *defaultLogger) Error(msg string, params ...Param) {
	if l.Enabled(wlog.ErrorLevel) {
		l.log(logtypes.LogLevelError, msg, params)
	}
}

func (l *defaultLogger) log(level logtypes.LogLevel, msg string, params []Param) {
	log := poolServiceLogV1.Get()
	populateLogObject(log, level, msg, params)
	l.logger.LogObject(log)
	poolServiceLogV1.Put(log)
}

func (l *defaultLogger) SetLevel(level wlog.LogLevel) {
	l.delegate.SetLevel(level)
}

func (l *defaultLogger) Enabled(level wlog.LogLevel) bool {
	return l.level == nil || l.level.Enabled(level)
}

// populateLogObject pre-populates the type, time, level, and message fields and applies user params in-place.
func populateLogObject(log *logtypes.ServiceLogV1, level logtypes.LogLevel, message string, params []Param) {
	log.Type = TypeValue
	log.Time = datetime.DateTime(time.Now())
	log.Level = level
	log.Message = message
	wloginternal.ApplyParams(log, params)
}

func ToParams(level wlog.Param, inParams []Param) []wlog.Param {
	return []wlog.Param{
		wlog.NewParam(func(entry wlog.LogEntry) {
			entry.StringValue(wlog.TypeKey, TypeValue)
			entry.StringValue(wlog.TimeKey, time.Now().Format(time.RFC3339Nano))
		}),
		level,
		// Apply all parameters to one object so that a later empty value clears an earlier one.
		wlog.NewParam(func(entry wlog.LogEntry) {
			log := new(logtypes.ServiceLogV1)
			wloginternal.ApplyParams(log, inParams)
			logentry.ServiceLogV1Params(entry, log)
		}),
	}
}

var (
	debugLevelParam = wlog.StringParam(LevelKey, LevelDebugValue)
	infoLevelParam  = wlog.StringParam(LevelKey, LevelInfoValue)
	warnLevelParam  = wlog.StringParam(LevelKey, LevelWarnValue)
	errorLevelParam = wlog.StringParam(LevelKey, LevelErrorValue)
)

func DebugLevelParam() wlog.Param { return debugLevelParam }
func InfoLevelParam() wlog.Param  { return infoLevelParam }
func WarnLevelParam() wlog.Param  { return warnLevelParam }
func ErrorLevelParam() wlog.Param { return errorLevelParam }
