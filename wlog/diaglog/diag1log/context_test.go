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
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/diaglog/diag1log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithLoggerFromContextRoundTrip(t *testing.T) {
	expected := diag1log.New(ioutil.Discard)
	ctx := diag1log.WithLogger(context.Background(), expected)
	got := diag1log.FromContext(ctx)
	assert.Equal(t, expected, got)
}

func TestFromContextWithoutLogger(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "stdout")
	require.NoError(t, err)
	defer func() { _ = f.Close() }()
	previousStdout, previousProvider := os.Stdout, wlog.DefaultLoggerProvider()
	os.Stdout = f
	wlog.SetDefaultLoggerProvider(wlog.NewJSONMarshalLoggerProvider())
	defer func() {
		os.Stdout = previousStdout
		wlog.SetDefaultLoggerProvider(previousProvider)
	}()
	diag1log.FromContext(context.Background()).Diagnostic(logging.NewDiagnosticFromGeneric(logging.GenericDiagnostic{
		DiagnosticType: "test",
		Value:          "value",
	}))
	data, err := os.ReadFile(f.Name())
	require.NoError(t, err)
	var record logging.DiagnosticLogV1
	require.NoError(t, json.Unmarshal(data, &record))
	require.Equal(t, diag1log.TypeValue, record.Type)
}
