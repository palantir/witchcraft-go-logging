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

package req2log

import (
	"context"
	"os"
)

type req2LogContextKeyType string

const contextKey = req2LogContextKeyType(TypeValue)

// WithLogger returns a copy of the provided context with the provided Logger included as a value. This operation will
// replace any logger that was previously set on the context.
func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKey, logger)
}

// FromContext returns the Logger stored in the provided context. If no logger is set on the context, returns the logger
// created by calling DefaultLogger.
func FromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(contextKey).(Logger); ok {
		return logger
	}
	return New(os.Stdout)
}
