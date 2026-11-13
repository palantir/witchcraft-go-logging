// Copyright (c) 2025 Palantir Technologies. All rights reserved.
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

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit2log"
	"github.com/palantir/witchcraft-go-logging/wlog/auditlog/audit2log/audit2logtests"
	"github.com/palantir/witchcraft-go-logging/wlog/diaglog/diag1log"
	"github.com/palantir/witchcraft-go-logging/wlog/diaglog/diag1log/diag1logtests"
	"github.com/palantir/witchcraft-go-logging/wlog/evtlog/evt2log"
	"github.com/palantir/witchcraft-go-logging/wlog/evtlog/evt2log/evt2logtests"
	"github.com/palantir/witchcraft-go-logging/wlog/metriclog/metric1log"
	"github.com/palantir/witchcraft-go-logging/wlog/metriclog/metric1log/metric1logtests"
	"github.com/palantir/witchcraft-go-logging/wlog/reqlog/req2log"
	"github.com/palantir/witchcraft-go-logging/wlog/reqlog/req2log/req2logtests"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log/svc1logtests"
	"github.com/palantir/witchcraft-go-logging/wlog/trclog/trc1log"
	"github.com/palantir/witchcraft-go-logging/wlog/trclog/trc1log/trc1logtests"
	"github.com/palantir/witchcraft-go-logging/wlog/wlogtests"
	"github.com/palantir/witchcraft-go-logging/wlog/wrappedlog/wrapped1log"
	"github.com/palantir/witchcraft-go-logging/wlog/wrappedlog/wrapped1log/wrapped1logtests"
	"github.com/stretchr/testify/require"
)

func jsonProviders() map[string]wlog.LoggerProvider {
	return map[string]wlog.LoggerProvider{
		"object": wlog.NewJSONMarshalLoggerProvider(),
	}
}

func TestJSONProviders(t *testing.T) {
	for name, provider := range jsonProviders() {
		t.Run(name, func(t *testing.T) {
			t.Run("wlog", func(t *testing.T) {
				wlogtests.JSONTestSuite(t, provider.NewLogger)
			})
			t.Run("service", func(t *testing.T) {
				svc1logtests.JSONTestSuite(t, func(w io.Writer, level wlog.LogLevel, origin string) svc1log.Logger {
					return svc1log.NewFromCreator(w, level, provider.NewLeveledLogger, svc1log.Origin(origin))
				})
			})
			t.Run("request", func(t *testing.T) {
				req2logtests.JSONTestSuite(t, func(w io.Writer, params ...req2log.LoggerCreatorParam) req2log.Logger {
					return req2log.NewFromCreator(w, provider.NewLogger, params...)
				})
			})
			t.Run("event", func(t *testing.T) {
				evt2logtests.JSONTestSuite(t, func(w io.Writer) evt2log.Logger {
					return evt2log.NewFromCreator(w, provider.NewLogger)
				})
			})
			t.Run("trace", func(t *testing.T) {
				trc1logtests.JSONTestSuite(t, func(w io.Writer) trc1log.Logger {
					return trc1log.NewFromCreator(w, provider.NewLogger)
				})
			})
			t.Run("metric", func(t *testing.T) {
				metric1logtests.JSONTestSuite(t, func(w io.Writer) metric1log.Logger {
					return metric1log.NewFromCreator(w, provider.NewLogger)
				})
			})
			t.Run("audit2", func(t *testing.T) {
				audit2logtests.JSONTestSuite(t, func(w io.Writer) audit2log.Logger {
					return audit2log.NewFromCreator(w, provider.NewLogger)
				})
			})
			t.Run("diagnostic", func(t *testing.T) {
				diag1logtests.JSONTestSuite(t, func(w io.Writer) diag1log.Logger {
					return diag1log.NewFromCreator(w, provider.NewLogger)
				})
			})
			const entityName, entityVersion = "entity", "version"
			t.Run("wrapped/audit2", func(t *testing.T) {
				wrapped1logtests.Audit2LogJSONTestSuite(t, entityName, entityVersion, func(w io.Writer) audit2log.Logger {
					return wrapped1log.NewFromProvider(w, wlog.InfoLevel, provider, entityName, entityVersion).Audit()
				})
			})
			t.Run("wrapped/diagnostic", func(t *testing.T) {
				wrapped1logtests.Diag1LogJSONTestSuite(t, entityName, entityVersion, func(w io.Writer) diag1log.Logger {
					return wrapped1log.NewFromProvider(w, wlog.InfoLevel, provider, entityName, entityVersion).Diagnostic()
				})
			})
			t.Run("wrapped/event", func(t *testing.T) {
				wrapped1logtests.Evt2LogJSONTestSuite(t, entityName, entityVersion, func(w io.Writer) evt2log.Logger {
					return wrapped1log.NewFromProvider(w, wlog.InfoLevel, provider, entityName, entityVersion).Event()
				})
			})
			t.Run("wrapped/metric", func(t *testing.T) {
				wrapped1logtests.Metric1LogJSONTestSuite(t, entityName, entityVersion, func(w io.Writer) metric1log.Logger {
					return wrapped1log.NewFromProvider(w, wlog.InfoLevel, provider, entityName, entityVersion).Metric()
				})
			})
			t.Run("wrapped/request", func(t *testing.T) {
				wrapped1logtests.Req2LogJSONTestSuite(t, entityName, entityVersion, func(w io.Writer, params ...req2log.LoggerCreatorParam) req2log.Logger {
					return wrapped1log.NewFromProvider(w, wlog.InfoLevel, provider, entityName, entityVersion).Request(params...)
				})
			})
			t.Run("wrapped/service", func(t *testing.T) {
				wrapped1logtests.Svc1LogJSONTestSuite(t, entityName, entityVersion, func(w io.Writer, level wlog.LogLevel, origin string) svc1log.Logger {
					return wrapped1log.NewFromProvider(w, level, provider, entityName, entityVersion).Service(svc1log.Origin(origin))
				})
			})
			t.Run("wrapped/trace", func(t *testing.T) {
				wrapped1logtests.Trc1LogJSONTestSuite(t, entityName, entityVersion, func(w io.Writer) trc1log.Logger {
					return wrapped1log.NewFromProvider(w, wlog.InfoLevel, provider, entityName, entityVersion).Trace()
				})
			})
		})
	}
}

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

// Each log type sets its own time; wrapped.1 has none on its envelope.
func TestJSONEntryWritesNoTimestamp(t *testing.T) {
	for name, provider := range jsonProviders() {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			provider.NewLogger(&buf).Log(wlog.StringParam("key", "value"))
			require.JSONEq(t, `{"key":"value"}`, buf.String())
		})
	}
}
