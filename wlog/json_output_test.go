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
	"io"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/diaglog/diag1log"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/palantir/witchcraft-go-logging/wlog/wrappedlog/wrapped1log"
	"github.com/stretchr/testify/require"
)

type blockingJSONValue struct {
	entered chan<- struct{}
	release <-chan struct{}
}

func (v blockingJSONValue) MarshalJSON() ([]byte, error) {
	v.entered <- struct{}{}
	<-v.release
	return []byte(`"value"`), nil
}

func TestJSONMarshalingCanOverlap(t *testing.T) {
	for name, provider := range jsonProviders() {
		t.Run(name, func(t *testing.T) {
			entered := make(chan struct{}, 2)
			release := make(chan struct{})
			logger := svc1log.NewFromCreator(io.Discard, wlog.InfoLevel, provider.NewLeveledLogger)
			var wg sync.WaitGroup
			t.Cleanup(func() {
				close(release)
				wg.Wait()
			})
			for range 2 {
				wg.Go(func() {
					logger.Info("message", svc1log.SafeParam("key", blockingJSONValue{entered: entered, release: release}))
				})
			}
			timer := time.NewTimer(5 * time.Second)
			defer timer.Stop()
			for range 2 {
				select {
				case <-entered:
				case <-timer.C:
					t.Fatal("concurrent log calls did not reach marshaling together")
				}
			}
		})
	}
}

func TestJSONConcurrentRecords(t *testing.T) {
	for name, provider := range jsonProviders() {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			backend := provider.NewLeveledLogger(&buf, wlog.InfoLevel)
			logger := svc1log.NewFromCreator(&buf, wlog.InfoLevel, func(io.Writer, wlog.LogLevel) wlog.LeveledLogger {
				return backend
			})
			var wg sync.WaitGroup
			for range 8 {
				wg.Go(func() {
					for range 25 {
						logger.Info("typed")
						backend.(wlog.Logger).Log(wlog.StringParam("message", "entry"))
					}
				})
			}
			wg.Wait()
			lines := bytes.Split(bytes.TrimSpace(buf.Bytes()), []byte("\n"))
			require.Len(t, lines, 400)
			counts := map[string]int{}
			for _, line := range lines {
				var record map[string]any
				require.NoError(t, json.Unmarshal(line, &record))
				counts[record["message"].(string)]++
			}
			require.Equal(t, map[string]int{"typed": 200, "entry": 200}, counts)
		})
	}
}

func TestJSONRecoversAfterMarshalError(t *testing.T) {
	cycle := map[string]any{}
	cycle["self"] = cycle
	for name, provider := range jsonProviders() {
		for _, wrapped := range []bool{false, true} {
			for badName, bad := range map[string]any{"function": func() {}, "cycle": cycle} {
				t.Run(name+"/"+badName+"/wrapped="+strconv.FormatBool(wrapped), func(t *testing.T) {
					var buf bytes.Buffer
					var logger svc1log.Logger
					if wrapped {
						logger = wrapped1log.NewFromProvider(&buf, wlog.InfoLevel, provider, "entity", "version").Service()
					} else {
						logger = svc1log.NewFromCreator(&buf, wlog.InfoLevel, provider.NewLeveledLogger)
					}
					logger.Info("before")
					logger.Info(strings.Repeat("bad", 10000), svc1log.SafeParam("unsupported", bad))
					logger.Info("after")
					lines := bytes.Split(bytes.TrimSpace(buf.Bytes()), []byte("\n"))
					require.Len(t, lines, 2)
					for i, message := range []string{"before", "after"} {
						var record map[string]any
						require.NoError(t, json.Unmarshal(lines[i], &record))
						if wrapped {
							record = record["payload"].(map[string]any)["serviceLogV1"].(map[string]any)
						}
						require.Equal(t, message, record["message"])
					}
				})
			}
		}
	}
}

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

type failOnceWriter struct {
	bytes.Buffer
	failed bool
}

func (w *failOnceWriter) Write(p []byte) (int, error) {
	if !w.failed {
		w.failed = true
		return 0, errors.New("write failed")
	}
	return w.Buffer.Write(p)
}

func TestJSONRecoversAfterWriterError(t *testing.T) {
	for name, provider := range jsonProviders() {
		t.Run(name, func(t *testing.T) {
			var writer failOnceWriter
			logger := svc1log.NewFromCreator(&writer, wlog.InfoLevel, provider.NewLeveledLogger)
			logger.Info("failed")
			logger.Info("after")
			var record map[string]any
			require.NoError(t, json.Unmarshal(writer.Bytes(), &record))
			require.Equal(t, "after", record["message"])
		})
	}
}

func TestJSONEntryTimestamp(t *testing.T) {
	for name, provider := range jsonProviders() {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := provider.NewLogger(&buf)
			logger.Log(wlog.StringParam("key", "value"))
			var record map[string]any
			require.NoError(t, json.Unmarshal(buf.Bytes(), &record))
			_, err := time.Parse(time.RFC3339Nano, record["time"].(string))
			require.NoError(t, err)
			buf.Reset()
			logger.Log(wlog.StringParam("time", "supplied"))
			require.NoError(t, json.Unmarshal(buf.Bytes(), &record))
			require.Equal(t, "supplied", record["time"])
		})
	}
}
