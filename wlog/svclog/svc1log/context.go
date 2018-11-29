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
	"context"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.palantir.build/deployability/witchcraft-params-go"
)

type svc1LogContextKeyType string

const contextKey = svc1LogContextKeyType(TypeValue)

// WithLogger returns a copy of the provided context with the provided Logger included as a value. This operation will
// replace any logger that was previously set on the context (along with all parameters that may have been set on the
// logger).
func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKey, logger)
}

// WithLoggerParams returns a copy of the provided context whose logger is configured with the provided parameters. If
// the provided context does not have a logger set on it or if no parameters are provided, the original context is
// returned unmodified. If any of the provided parameters set safe or unsafe values, the returned context will also have
// those values set on it using the wparams.ContextWithSafeAndUnsafeParams function.
func WithLoggerParams(ctx context.Context, params ...Param) context.Context {
	logger := loggerFromContext(ctx)
	if logger == nil || len(params) == 0 {
		return ctx
	}
	// if the provided params set any safe or unsafe values, set those as wparams on the context
	if safeParams, unsafeParams := safeAndUnsafeParamsFromParams(params); len(safeParams) > 0 || len(unsafeParams) > 0 {
		ctx = wparams.ContextWithSafeAndUnsafeParams(ctx, safeParams, unsafeParams)
	}
	return WithLogger(ctx, WithParams(logger, params...))
}

// FromContext returns the Logger stored in the provided context or nil if no logger is set on the context. If a logger
// is returned, the returned logger has any safe or unsafe parameters stored on the context using wparams set on it.
func FromContext(ctx context.Context) Logger {
	logger := loggerFromContext(ctx)
	if logger == nil {
		return nil
	}
	if paramStorer := wparams.ParamStorerFromContext(ctx); paramStorer != nil && (len(paramStorer.SafeParams()) > 0 || len(paramStorer.UnsafeParams()) > 0) {
		logger = WithParams(logger, Params(paramStorer))
	}
	return logger
}

func safeAndUnsafeParamsFromParams(params []Param) (safe map[string]interface{}, unsafe map[string]interface{}) {
	logEntry := wlog.NewMapLogEntry()
	for _, currParam := range params {
		currParam.apply(logEntry)
	}
	return logEntry.AnyMapValues()[ParamsKey], logEntry.AnyMapValues()[wlog.UnsafeParamsKey]
}

// FromContext returns the Logger stored in the provided context or nil if no logger is set on the context.
func loggerFromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(contextKey).(Logger); ok {
		return logger
	}
	return nil
}
