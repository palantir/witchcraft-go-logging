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
	"io"
	"os"

	"github.com/palantir/witchcraft-go-logging/wapi/logging"
	wlog "github.com/palantir/witchcraft-go-logging/wlog2"
)

var (
	DefaultOutput = os.Stdout

	levelDebug = logging.New_LogLevel(logging.LogLevel_DEBUG)
	levelInfo  = logging.New_LogLevel(logging.LogLevel_INFO)
	levelWarn  = logging.New_LogLevel(logging.LogLevel_WARN)
	levelError = logging.New_LogLevel(logging.LogLevel_ERROR)
)

type Logger interface {
	Debug(msg string, params ...Param)
	Info(msg string, params ...Param)
	Warn(msg string, params ...Param)
	Error(msg string, params ...Param)
	SetLevel(level wlog.LogLevel)
	WithParams(params ...Param) Logger
}

func New(w io.Writer, level wlog.LogLevel, params ...Param) Logger {
	return &defaultLogger{
		logger: wlog.NewDefaultLogger(w, Type(), TimeNow()).WithParams(params...),
		level:  wlog.NewAtomicLogLevel(level),
	}
}

type defaultLogger struct {
	level  *wlog.AtomicLogLevel
	logger wlog.ConjureLogger[logging.ServiceLogV1]
}

func (l *defaultLogger) WithParams(params ...Param) Logger {
	return &defaultLogger{
		level:  l.level,
		logger: l.logger.WithParams(params...),
	}
}

func (l *defaultLogger) Debug(msg string, params ...Param) {
	l.doLog(wlog.DebugLevel, levelDebug, msg, params...)
}

func (l *defaultLogger) Info(msg string, params ...Param) {
	l.doLog(wlog.InfoLevel, levelInfo, msg, params...)
}

func (l *defaultLogger) Warn(msg string, params ...Param) {
	l.doLog(wlog.WarnLevel, levelWarn, msg, params...)
}

func (l *defaultLogger) Error(msg string, params ...Param) {
	l.doLog(wlog.ErrorLevel, levelError, msg, params...)
}

func (l *defaultLogger) doLog(wLevel wlog.LogLevel, cLevel logging.LogLevel, msg string, params ...Param) {
	if l.level.Enabled(wLevel) {
		l.logger.Log(append([]Param{withLevel(cLevel), Message(msg)}, params...)...)
	}
}

func (l *defaultLogger) SetLevel(level wlog.LogLevel) {
	l.level.SetLevel(level)
}
