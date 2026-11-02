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
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
)

var (
	LevelDebug = logging.New_LogLevel(logging.LogLevel_DEBUG)
	LevelInfo  = logging.New_LogLevel(logging.LogLevel_INFO)
	LevelWarn  = logging.New_LogLevel(logging.LogLevel_WARN)
	LevelError = logging.New_LogLevel(logging.LogLevel_ERROR)
)

// conjureLogger consumes borrowed log objects synchronously. Implementations must copy
// any object or nested data they retain after the method returns.
type conjureLogger interface {
	LogServiceLogV1(*logging.ServiceLogV1)
}

type defaultLogger struct {
	logger   conjureLogger
	delegate wlog.LeveledLogger
	level    wlog.LevelChecker
}

func (l *defaultLogger) Debug(msg string, params ...Param) {
	if l.Enabled(wlog.DebugLevel) {
		l.log(LevelDebug, msg, params)
	}
}

func (l *defaultLogger) Info(msg string, params ...Param) {
	if l.Enabled(wlog.InfoLevel) {
		l.log(LevelInfo, msg, params)
	}
}

func (l *defaultLogger) Warn(msg string, params ...Param) {
	if l.Enabled(wlog.WarnLevel) {
		l.log(LevelWarn, msg, params)
	}
}

func (l *defaultLogger) Error(msg string, params ...Param) {
	if l.Enabled(wlog.ErrorLevel) {
		l.log(LevelError, msg, params)
	}
}

func (l *defaultLogger) log(level logging.LogLevel, msg string, params []Param) {
	log := poolServiceLogV1.Get()
	populateLogObject(log, level, msg, params)
	l.logger.LogServiceLogV1(log)
	poolServiceLogV1.Put(log)
}

func (l *defaultLogger) SetLevel(level wlog.LogLevel) {
	l.delegate.SetLevel(level)
}

func (l *defaultLogger) Enabled(level wlog.LogLevel) bool {
	return l.level == nil || l.level.Enabled(level)
}

type entryLogger struct {
	logger wlog.LeveledLogger
}

func (l entryLogger) LogServiceLogV1(log *logging.ServiceLogV1) {
	param := ObjectToParam(log)
	switch log.Level.Value() {
	case logging.LogLevel_DEBUG:
		l.logger.Debug(log.Message, param)
	case logging.LogLevel_INFO:
		l.logger.Info(log.Message, param)
	case logging.LogLevel_WARN:
		l.logger.Warn(log.Message, param)
	case logging.LogLevel_ERROR:
		l.logger.Error(log.Message, param)
	}
}

// populateLogObject pre-populates the type, time, level, and message fields and applies user params in-place.
func populateLogObject(log *logging.ServiceLogV1, level logging.LogLevel, message string, params []Param) {
	log.Type = TypeValue
	log.Time = datetime.DateTime(time.Now())
	log.Level = level
	log.Message = message
	wloginternal.ApplyParams(log, params)
	wloginternal.OmitEmpty(&log.Origin)
	wloginternal.OmitEmpty(&log.Uid)
	wloginternal.OmitEmpty(&log.Sid)
	wloginternal.OmitEmpty(&log.TokenId)
	wloginternal.OmitEmpty(&log.OrgId)
	wloginternal.OmitEmpty(&log.TraceId)
}

func ToParams(level wlog.Param, inParams []Param) []wlog.Param {
	out := make([]wlog.Param, 0, 2+len(inParams))
	out = append(out, wlog.NewParam(func(entry wlog.LogEntry) {
		entry.StringValue(wlog.TypeKey, TypeValue)
		entry.StringValue(wlog.TimeKey, time.Now().Format(time.RFC3339Nano))
	}), level)
	for _, p := range inParams {
		out = append(out, wlog.NewParam(func(entry wlog.LogEntry) {
			ApplyParam(p, entry)
		}))
	}
	return out
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

func ToObject(level logging.LogLevel, message string, inParams []Param) logging.ServiceLogV1 {
	var log logging.ServiceLogV1
	populateLogObject(&log, level, message, inParams)
	return log
}
