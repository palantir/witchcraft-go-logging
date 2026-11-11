// Copyright (c) 2025 Palantir Technologies. All rights reserved.
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

package audit3log

import (
	"context"

	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
)

type audit3LogContextKeyType string

const contextKey = audit3LogContextKeyType(TypeValue)

// WithLogger returns a copy of the provided context with the provided Logger included as a value. This operation will
// replace any logger that was previously set on the context (along with all parameters that may have been set on the
// logger).
func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKey, logger)
}

// WithLoggerParams configures the context's logger with params, using the default logger if none is set.
// With no parameters, it returns the original context.
func WithLoggerParams(ctx context.Context, params ...Param) context.Context {
	if len(params) == 0 {
		return ctx
	}
	return WithLogger(ctx, WithParams(loggerFromContext(ctx), params...))
}

// FromContext returns the context's logger, or one created by the configured default logger creator.
// Context IDs and the wtracing trace ID are applied to the returned logger.
func FromContext(ctx context.Context) Logger {
	logger := loggerFromContext(ctx)
	var params []Param
	if uid := wloginternal.IDFromContext(ctx, wloginternal.UIDKey); uid != nil {
		params = append(params, UID(*uid))
	}
	if sid := wloginternal.IDFromContext(ctx, wloginternal.SIDKey); sid != nil {
		params = append(params, SID(*sid))
	}
	if tokenID := wloginternal.IDFromContext(ctx, wloginternal.TokenIDKey); tokenID != nil {
		params = append(params, TokenID(*tokenID))
	}
	if orgID := wloginternal.IDFromContext(ctx, wloginternal.OrgIDKey); orgID != nil {
		params = append(params, OrgID(*orgID))
	}
	if traceID := wtracing.TraceIDFromContext(ctx); traceID != "" {
		params = append(params, TraceID(string(traceID)))
	}
	return WithParams(logger, params...)
}

// IsSetOnContext returns true if an audit.3 logger is set on the provided context, false otherwise.
// Note that, even if this function returns "false", that does not mean that "FromContext" cannot called on this
// context, as "FromContext" will use the defaultLoggerCreator to create a logger if none is set on the context.
func IsSetOnContext(ctx context.Context) bool {
	_, ok := ctx.Value(contextKey).(Logger)
	return ok
}

// loggerFromContext returns the logger stored in the provided context. If no logger is set on the context, returns the
// logger created by the configured default logger creator.
func loggerFromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(contextKey).(Logger); ok {
		return logger
	}
	return defaultLoggerCreator()
}
