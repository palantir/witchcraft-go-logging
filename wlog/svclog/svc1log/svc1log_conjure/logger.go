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

	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
)

var (
	DefaultOutput = os.Stdout
)

type Logger interface {
	Debug(msg string, params ...Param)
	Info(msg string, params ...Param)
	Warn(msg string, params ...Param)
	Error(msg string, params ...Param)
	SetLevel(level wlog.LogLevel)
}

type LoggerWithParams interface {
	Logger
	WithParams(params ...Param) Logger
}

func New(w io.Writer, level wlog.LogLevel, params ...Param) LoggerWithParams {
	return &defaultLogger{
		logger: wlog.NewDefaultLogger[logging.ServiceLogV1](w).WithParams(Type(), TimeNow()).WithParams(params...),
		level:  wlog.NewAtomicLogLevel(level),
	}
}

type defaultLogger struct {
	level  *wlog.AtomicLogLevel
	logger wlog.ConjureLogger[logging.ServiceLogV1]
}

func (l *defaultLogger) Debug(msg string, params ...Param) { l.log(wlog.DebugLevel, msg, params...) }

func (l *defaultLogger) Info(msg string, params ...Param) { l.log(wlog.InfoLevel, msg, params...) }

func (l *defaultLogger) Warn(msg string, params ...Param) { l.log(wlog.WarnLevel, msg, params...) }

func (l *defaultLogger) Error(msg string, params ...Param) { l.log(wlog.ErrorLevel, msg, params...) }

func (l *defaultLogger) log(level wlog.LogLevel, msg string, params ...Param) {
	if l.level.Enabled(level) {
		l.logger.Log(append([]Param{Level(level), Message(msg)}, params...)...)
	}
}

func (l *defaultLogger) SetLevel(level wlog.LogLevel) {
	l.level.SetLevel(level)
}

func (l *defaultLogger) WithParams(params ...Param) Logger {
	return &defaultLogger{
		level:  l.level,
		logger: l.logger.WithParams(params...),
	}
}
