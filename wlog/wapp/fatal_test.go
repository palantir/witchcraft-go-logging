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

package wapp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/palantir/witchcraft-go-error"
	"github.com/palantir/witchcraft-go-logging/wlog"
	_ "github.com/palantir/witchcraft-go-logging/wlog-zap"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/stretchr/testify/assert"
)

func TestRunWithFatalLogging_Panic(t *testing.T) {
	buf := &bytes.Buffer{}
	ctx := getContextWithLogger(context.Background(), buf)
	err := RunWithFatalLogging(ctx, func(ctx context.Context) error {
		panic("foo")
	})
	assert.Error(t, err)
	st, safe := werror.ParamFromError(err, "stacktrace")
	assert.True(t, safe, "Expected stacktrace param to be safe")
	assert.NotNil(t, st, "Expected stacktrace param to not be nil")
	r, safe := werror.ParamFromError(err, "recovered")
	assert.False(t, safe, "Expected recovered param to be unsafe")
	assert.NotNil(t, r, "Expected recovered param value to not be nil")
}

func TestRunWithFatalLogging_Error(t *testing.T) {
	buf := &bytes.Buffer{}
	ctx := getContextWithLogger(context.Background(), buf)
	err := RunWithFatalLogging(ctx, func(ctx context.Context) error {
		return werror.Error("foo")
	})
	assert.NotNil(t, err)
	assert.Contains(t, buf.String(), "foo")
}

func TestRunWithFatalLogging_PanicWithError(t *testing.T) {
	buf := &bytes.Buffer{}
	ctx := getContextWithLogger(context.Background(), buf)
	var panicErr error
	err := RunWithFatalLogging(ctx, func(ctx context.Context) error {
		panicErr = werror.Error("foo", werror.SafeParam("verySafeParam", "blah"), werror.UnsafeParam("notSafeParam", "oogabooga"))
		panic(panicErr)
	})
	assert.Error(t, err)

	errorWithStacktrace := fmt.Sprintf("%+v", err)
	assert.Contains(t, errorWithStacktrace, panicErr.Error(), "Expected error returned through panic to be included in the stacktrace")

	// assert params from original error are present
	st, safe := werror.ParamFromError(err, "stacktrace")
	assert.True(t, safe, "Expected stacktrace param to be safe")
	assert.NotNil(t, st, "Expected stacktrace param to not be nil")
	vsp, safe := werror.ParamFromError(err, "verySafeParam")
	assert.True(t, safe, "Expected verySafeParam param to be safe")
	assert.Equal(t, "blah", vsp, "Expected verySafeParam param to match what was returned")
	nsp, safe := werror.ParamFromError(err, "notSafeParam")
	assert.False(t, safe, "Expected notSafeParam param to be unsafe")
	assert.Equal(t, "oogabooga", nsp, "Expected notSafeParam param to match what was returned")

	// assert original error is listed as the cause
	expectedPanicErr := werror.RootCause(err)
	assert.Equal(t, panicErr, expectedPanicErr, "Expected panic error to be root cause of recovered error")
	assert.EqualError(t, expectedPanicErr, "foo")
}

func getContextWithLogger(ctx context.Context, writer io.Writer) context.Context {
	logger := svc1log.New(writer, wlog.DebugLevel)
	ctx = svc1log.WithLogger(ctx, logger)
	return ctx
}
