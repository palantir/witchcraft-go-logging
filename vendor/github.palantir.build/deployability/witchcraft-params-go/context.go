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

package wparams

import (
	"context"
)

type witchcraftParamsContextKeyType string

const contextKey = witchcraftParamsContextKeyType("witchcraftParams")

// ContextWithParams returns a copy of the provided context that contains all of the safe and unsafe parameters provided
// by the provided ParamStorers. If the provided context already has safe/unsafe params, the newly returned context will
// contain the result of merging the previous parameters with the provided parameters.
func ContextWithParamStorers(ctx context.Context, params ...ParamStorer) context.Context {
	return context.WithValue(ctx, contextKey, NewParamStorer(append([]ParamStorer{ParamStorerFromContext(ctx)}, params...)...))
}

func ContextWithSafeParam(ctx context.Context, key string, value interface{}) context.Context {
	return ContextWithSafeParams(ctx, map[string]interface{}{
		key: value,
	})
}

func ContextWithSafeParams(ctx context.Context, safeParams map[string]interface{}) context.Context {
	return ContextWithParamStorers(ctx, NewSafeParamStorer(safeParams))
}

func ContextWithUnsafeParam(ctx context.Context, key string, value interface{}) context.Context {
	return ContextWithUnsafeParams(ctx, map[string]interface{}{
		key: value,
	})
}

func ContextWithUnsafeParams(ctx context.Context, unsafeParams map[string]interface{}) context.Context {
	return ContextWithParamStorers(ctx, NewUnsafeParamStorer(unsafeParams))
}

func ContextWithSafeAndUnsafeParams(ctx context.Context, safeParams, unsafeParams map[string]interface{}) context.Context {
	return ContextWithParamStorers(ctx, NewSafeAndUnsafeParamStorer(safeParams, unsafeParams))
}

// ParamStorerFromContext returns the ParamStorer stored in the provided context. Returns nil if the provided context
// does not contain a ParamStorer.
func ParamStorerFromContext(ctx context.Context) ParamStorer {
	val := ctx.Value(contextKey)
	if paramStorer, ok := val.(ParamStorer); ok {
		return paramStorer
	}
	return nil
}

// SafeAndUnsafeParamsFromContext returns the safe and unsafe parameters stored in the ParamStorer returned by
// ParamStorerFromContext for the provided context. Returns nil maps if the provided context does not have a
// ParamStorer.
func SafeAndUnsafeParamsFromContext(ctx context.Context) (safeParams map[string]interface{}, unsafeParams map[string]interface{}) {
	storer := ParamStorerFromContext(ctx)
	if storer == nil {
		return nil, nil
	}
	return storer.SafeParams(), storer.UnsafeParams()
}
