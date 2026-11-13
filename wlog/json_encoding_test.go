// Copyright (c) 2026 Palantir Technologies. All rights reserved.
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

package wlog_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/diaglog/diag1log"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/palantir/witchcraft-go-logging/wlog/wrappedlog/wrapped1log"
	"github.com/stretchr/testify/require"
)

type jsonError struct{}

func (jsonError) Error() string { return "json error" }

func (jsonError) MarshalJSON() ([]byte, error) { return []byte(`{"custom":true}`), nil }

type textAppendError struct{}

func (textAppendError) Error() string { return "error message" }

func (textAppendError) AppendText(dst []byte) ([]byte, error) {
	return append(dst, "custom text"...), nil
}

type unexportedFields struct{ x int }

func TestJSONEncodesV1CompatibleValues(t *testing.T) {
	for name, provider := range jsonProviders() {
		for _, wrapped := range []bool{false, true} {
			t.Run(name+"/wrapped="+strconv.FormatBool(wrapped), func(t *testing.T) {
				var buf bytes.Buffer
				var logger svc1log.Logger
				if wrapped {
					logger = wrapped1log.NewFromProvider(&buf, wlog.InfoLevel, provider, "entity", "version").Service()
				} else {
					logger = svc1log.NewFromCreator(&buf, wlog.InfoLevel, provider.NewLeveledLogger)
				}
				logger.Info("message\xff", svc1log.SafeParams(map[string]any{
					"duration":   time.Second,
					"error":      errors.New("failed"),
					"jsonError":  jsonError{},
					"textError":  textAppendError{},
					"unexported": unexportedFields{x: 1},
				}))
				var record map[string]any
				require.NoError(t, json.Unmarshal(buf.Bytes(), &record), buf.String())
				if wrapped {
					record = record["payload"].(map[string]any)["serviceLogV1"].(map[string]any)
				}
				require.Equal(t, "message�", record["message"])
				require.Equal(t, map[string]any{
					"duration":   float64(time.Second),
					"error":      "failed",
					"jsonError":  map[string]any{"custom": true},
					"textError":  "custom text",
					"unexported": map[string]any{},
				}, record["params"])
			})
		}
	}
}

func TestJSONSortsMapKeys(t *testing.T) {
	for name, provider := range jsonProviders() {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := svc1log.NewFromCreator(&buf, wlog.InfoLevel, provider.NewLeveledLogger)
			logger.Info("message", svc1log.SafeParams(map[string]any{"e": 5, "c": 3, "a": 1, "d": 4, "b": 2}))
			require.Contains(t, buf.String(), `"params":{"a":1,"b":2,"c":3,"d":4,"e":5}`)
		})
	}
}

func TestJSONDiagnosticErrorValues(t *testing.T) {
	values := map[string]any{"error": errors.New("failed"), "custom": jsonError{}, "duration": time.Second}
	want := map[string]any{"error": "failed", "custom": map[string]any{"custom": true}, "duration": float64(time.Second)}
	for name, provider := range jsonProviders() {
		for _, wrapped := range []bool{false, true} {
			t.Run(name+"/wrapped="+strconv.FormatBool(wrapped), func(t *testing.T) {
				var buf bytes.Buffer
				logger := diag1log.NewFromCreator(&buf, provider.NewLogger)
				if wrapped {
					logger = wrapped1log.NewFromProvider(&buf, wlog.InfoLevel, provider, "entity", "version").Diagnostic()
				}
				logger.Diagnostic(logging.NewDiagnosticFromGeneric(logging.GenericDiagnostic{
					DiagnosticType: "example",
					Value:          values,
				}), diag1log.UnsafeParams(values))
				var record map[string]any
				require.NoError(t, json.Unmarshal(buf.Bytes(), &record), buf.String())
				if wrapped {
					record = record["payload"].(map[string]any)["diagnosticLogV1"].(map[string]any)
				}
				generic := record["diagnostic"].(map[string]any)["generic"].(map[string]any)
				require.Equal(t, want, generic["value"])
				require.Equal(t, want, record["unsafeParams"])
			})
		}
	}
}
