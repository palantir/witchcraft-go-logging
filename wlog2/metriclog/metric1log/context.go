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

package metric1log

import (
	"context"

	wloginternal "github.com/palantir/witchcraft-go-logging/wlog2/internal"
)

type contextKey struct{}

// WithLogger returns a copy of the provided context with the provided Logger included as a value. This operation will
// replace any logger that was previously set on the context (along with all parameters that may have been set on the
// logger).
func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}

// WithLoggerParams returns a copy of the provided context whose logger is configured with the provided parameters. If
// no parameters are provided, the original context is returned unmodified. If the provided context did not have a
// logger set on it, the returned context will contain the default logger configured with the provided parameters.
func WithLoggerParams(ctx context.Context, params ...Param) context.Context {
	if len(params) == 0 {
		return ctx
	}
	return WithLogger(ctx, loggerFromContext(ctx).WithParams(params...))
}

// FromContext returns the Logger stored in the provided context. If no logger is set on the context, returns the logger
// created by calling DefaultLogger. If the context contains a TraceID set using wtracing, the returned logger has that
// TraceID set on it as a parameter. Any safe or unsafe parameters stored on the context using wparams are also set as
// parameters on the returned logger.
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
	return logger.WithParams(params...)
}

// loggerFromContext returns the logger stored in the provided context. If no logger is set on the context, returns the
// logger created by calling DefaultLogger.
func loggerFromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(contextKey{}).(Logger); ok {
		return logger
	}
	return New(DefaultOutput)
}
