// Copyright (c) 2020 Palantir Technologies. All rights reserved.
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

package diag1log

import (
	"context"

	"github.com/palantir/witchcraft-go-logging/wtypes"
)

type diag1LogContextKeyType string

const contextKey = diag1LogContextKeyType(TypeValue)

// Diagnostic uses the logger stored in the context to emit a diagnostic.1 log.
func Diagnostic(ctx context.Context, diagnostic wtypes.Diagnostic, params ...Param) {
	FromContext(ctx).Diagnostic(diagnostic, params...)
}

func FromContext(ctx context.Context) Logger {
	untypedLogger := ctx.Value(contextKey)
	if untypedLogger == nil {
		return &defaultLogger{}
	}
	logger, ok := untypedLogger.(Logger)
	if !ok {
		return &defaultLogger{}
	}
	return logger
}

func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKey, logger)
}
