// Copyright (c) 2022 Palantir Technologies. All rights reserved.
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

package svc1zap

import (
	"strings"

	"github.com/palantir/witchcraft-go-logging/internal/gopath"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type svc1zapCore struct {
	log                svc1log.Logger
	mutator            func(entry zapcore.Entry) (zapcore.Entry, bool)
	originFromCallLine bool
	newParamFunc       func(key string, value interface{}) svc1log.Param
}

// New returns a zap logger that delegates to the provided svc1log logger.
// The enabled/disabled log level configuration on the returned zap logger is ignored in favor of the svc1log configuration.
func New(logger svc1log.Logger, opts ...Option) *zap.Logger {
	core := NewCore(logger, opts...)
	z := zap.New(core)
	if core.(*svc1zapCore).originFromCallLine {
		z = z.WithOptions(zap.AddCaller())
	}
	return z
}

func NewCore(logger svc1log.Logger, opts ...Option) zapcore.Core {
	core := &svc1zapCore{log: logger}
	for _, opt := range opts {
		opt(core)
	}
	return core
}

type Option func(*svc1zapCore)

// WithOriginFromZapCaller enables zap.AddCaller() and uses the caller file and line to construct the origin value.
// Similar to svc1log.OriginFromCallLine().
func WithOriginFromZapCaller() Option {
	return func(core *svc1zapCore) { core.originFromCallLine = true }
}

// WithNewParamFunc provides a function for constructing svc1log.Param values from zap fields.
// Use this option to control parameter safety. By default, all fields are converted to unsafe params.
// If newParam returns nil, the field is skipped.
func WithNewParamFunc(newParam func(key string, value interface{}) svc1log.Param) Option {
	return func(core *svc1zapCore) { core.newParamFunc = newParam }
}

// WithEntryMutatorFunc provides a function for modifying or skipping entries dynamically.
// If mutator is set, ok must return true for the message to be logged.
func WithEntryMutatorFunc(mutator func(entry zapcore.Entry) (out zapcore.Entry, ok bool)) Option {
	return func(core *svc1zapCore) { core.mutator = mutator }
}

func (c svc1zapCore) Enabled(level zapcore.Level) bool {
	if checker, ok := c.log.(wlog.LevelChecker); ok {
		switch level {
		case zapcore.DebugLevel:
			return checker.Enabled(wlog.DebugLevel)
		case zapcore.InfoLevel:
			return checker.Enabled(wlog.InfoLevel)
		case zapcore.WarnLevel:
			return checker.Enabled(wlog.WarnLevel)
		default:
			return checker.Enabled(wlog.ErrorLevel)
		}
	}
	return true
}

func (c svc1zapCore) With(fields []zapcore.Field) zapcore.Core {
	clone := c
	clone.log = svc1log.WithParams(c.log, c.fieldsToWlogParams(fields)...)
	return &clone
}

func (c svc1zapCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if !c.Enabled(entry.Level) {
		return ce
	}

	if c.mutator != nil {
		var ok bool
		entry, ok = c.mutator(entry)
		if !ok {
			return ce
		}
	}
	return ce.AddCore(entry, c)
}

func (c svc1zapCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	message := formatMessage(entry)
	params := c.fieldsToWlogParams(fields)
	if c.originFromCallLine && entry.Caller.Defined {
		params = append(params, svc1log.Origin(gopath.TrimPrefix(entry.Caller.FullPath())))
	}
	switch entry.Level {
	case zapcore.DebugLevel:
		c.log.Debug(message, params...)
	case zapcore.InfoLevel:
		c.log.Info(message, params...)
	case zapcore.WarnLevel:
		c.log.Warn(message, params...)
	default:
		c.log.Error(message, params...)
	}
	return nil
}

func (c svc1zapCore) Sync() error { return nil }

func (c svc1zapCore) fieldsToWlogParams(fields []zapcore.Field) []svc1log.Param {
	var params []svc1log.Param
	for key, value := range fieldsToMap(fields) {
		if c.newParamFunc != nil {
			if p := c.newParamFunc(key, value); p != nil {
				params = append(params, p)
			}
		} else {
			params = append(params, svc1log.UnsafeParam(key, value))
		}
	}
	return params
}

func formatMessage(entry zapcore.Entry) string {
	if entry.LoggerName == "" {
		return entry.Message
	}
	sb := strings.Builder{}
	sb.Grow(len(entry.LoggerName) + 2 + len(entry.Message))
	sb.WriteString(entry.LoggerName)
	sb.WriteString(": ")
	sb.WriteString(entry.Message)
	return sb.String()
}

func fieldsToMap(fields []zapcore.Field) map[string]interface{} {
	params := zapcore.NewMapObjectEncoder()
	for _, field := range fields {
		if field.Key == "token" {
			continue // who logs a token...?
		}
		field.AddTo(params)
	}
	return params.Fields
}
