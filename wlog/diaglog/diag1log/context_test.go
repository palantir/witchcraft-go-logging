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

package diag1log_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog/diaglog/diag1log"
	"github.com/stretchr/testify/assert"
)

func TestWithLoggerFromContextRoundTrip(t *testing.T) {
	expected := diag1log.New(ioutil.Discard)
	ctx := diag1log.WithLogger(context.Background(), expected)
	got := diag1log.FromContext(ctx)
	assert.Equal(t, expected, got)
}

type testDiagLogger struct{}

func (testDiagLogger) Diagnostic(diagnostic logging.Diagnostic, params ...diag1log.Param) {}
