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

package app

import (
	"bytes"
	"context"
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
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		assert.Contains(t, buf.String(), "panicking")
	}()
	ctx := getContextWithLogger(context.Background(), buf)
	_ = RunWithFatalLogging(ctx, func(ctx context.Context) error {
		panic("foo")
	})
	assert.Fail(t, "unexpected")
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

func getContextWithLogger(ctx context.Context, writer io.Writer) context.Context {
	logger := svc1log.New(writer, wlog.DebugLevel)
	ctx = svc1log.WithLogger(ctx, logger)
	return ctx
}
