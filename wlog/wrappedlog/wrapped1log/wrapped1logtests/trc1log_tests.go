// Copyright (c) 2021 Palantir Technologies. All rights reserved.
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

package wrapped1logtests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/palantir/pkg/objmatcher"
	"github.com/palantir/witchcraft-go-logging/wlog/logreader"
	"github.com/palantir/witchcraft-go-logging/wlog/trclog/trc1log"
	"github.com/palantir/witchcraft-go-tracing/wtracing"
	"github.com/palantir/witchcraft-go-tracing/wzipkin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Trc1LogJSONTestSuite(t *testing.T, entityName string, entityVersion string, loggerProvider func(w io.Writer) trc1log.Logger) {
	trc1LogJSONOutputTests(t, entityName, entityVersion, loggerProvider)
	durationFormatOutputTest(t, loggerProvider)
}

type Trc1TestCase struct {
	Name        string
	SpanOptions []wtracing.SpanOption
	JSONMatcher objmatcher.MapMatcher
}

func Trc1TestCases(entityName string, entityVersion string, clientSpan wtracing.Span) []Trc1TestCase {
	spanContext := clientSpan.Context()
	traceID := string(spanContext.TraceID)
	clientSpanID := string(spanContext.ID)
	return []Trc1TestCase{
		{
			Name: "trace.1 log entry",
			JSONMatcher: map[string]objmatcher.Matcher{
				"type":          objmatcher.NewEqualsMatcher("wrapped.1"),
				"entityName":    objmatcher.NewEqualsMatcher(entityName),
				"entityVersion": objmatcher.NewEqualsMatcher(entityVersion),
				"payload": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"type": objmatcher.NewEqualsMatcher("traceLogV1"),
					"traceLogV1": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"type": objmatcher.NewEqualsMatcher("trace.1"),
						"time": objmatcher.NewRegExpMatcher(".+"),
						"span": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"name":      objmatcher.NewEqualsMatcher("testOp"),
							"traceId":   objmatcher.NewEqualsMatcher(traceID),
							"id":        objmatcher.NewRegExpMatcher("[a-f0-9]+"),
							"parentId":  objmatcher.NewEqualsMatcher(clientSpanID),
							"timestamp": objmatcher.NewAnyMatcher(),
							"duration":  objmatcher.NewAnyMatcher(),
						}),
					}),
				}),
			},
		},
		{
			Name: "trace.1 log entry with server mode annotations",
			SpanOptions: []wtracing.SpanOption{
				wtracing.WithKind(wtracing.Server),
			},
			JSONMatcher: map[string]objmatcher.Matcher{
				"type":          objmatcher.NewEqualsMatcher("wrapped.1"),
				"entityName":    objmatcher.NewEqualsMatcher(entityName),
				"entityVersion": objmatcher.NewEqualsMatcher(entityVersion),
				"payload": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"type": objmatcher.NewEqualsMatcher("traceLogV1"),
					"traceLogV1": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"type": objmatcher.NewEqualsMatcher("trace.1"),
						"time": objmatcher.NewRegExpMatcher(".+"),
						"span": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"name":      objmatcher.NewEqualsMatcher("testOp"),
							"traceId":   objmatcher.NewEqualsMatcher(traceID),
							"id":        objmatcher.NewRegExpMatcher("[a-f0-9]+"),
							"parentId":  objmatcher.NewEqualsMatcher(clientSpanID),
							"timestamp": objmatcher.NewAnyMatcher(),
							"duration":  objmatcher.NewAnyMatcher(),
							"annotations": objmatcher.SliceMatcher([]objmatcher.Matcher{
								objmatcher.MapMatcher(map[string]objmatcher.Matcher{
									"value":     objmatcher.NewEqualsMatcher("sr"),
									"timestamp": objmatcher.NewAnyMatcher(),
									"endpoint": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
										"serviceName": objmatcher.NewEqualsMatcher("testService"),
										"ipv4":        objmatcher.NewEqualsMatcher("127.0.0.1"),
									}),
								}),
								objmatcher.MapMatcher(map[string]objmatcher.Matcher{
									"value":     objmatcher.NewEqualsMatcher("ss"),
									"timestamp": objmatcher.NewAnyMatcher(),
									"endpoint": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
										"serviceName": objmatcher.NewEqualsMatcher("testService"),
										"ipv4":        objmatcher.NewEqualsMatcher("127.0.0.1"),
									}),
								}),
							}),
						}),
					}),
				}),
			},
		},
		{
			Name: "trace.1 log entry with trace tags ",
			SpanOptions: []wtracing.SpanOption{
				wtracing.WithSpanTag("key0", "value0"),
			},
			JSONMatcher: map[string]objmatcher.Matcher{
				"type":          objmatcher.NewEqualsMatcher("wrapped.1"),
				"entityName":    objmatcher.NewEqualsMatcher(entityName),
				"entityVersion": objmatcher.NewEqualsMatcher(entityVersion),
				"payload": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"type": objmatcher.NewEqualsMatcher("traceLogV1"),
					"traceLogV1": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"type": objmatcher.NewEqualsMatcher("trace.1"),
						"time": objmatcher.NewRegExpMatcher(".+"),
						"span": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"name":      objmatcher.NewEqualsMatcher("testOp"),
							"traceId":   objmatcher.NewEqualsMatcher(traceID),
							"id":        objmatcher.NewRegExpMatcher("[a-f0-9]+"),
							"parentId":  objmatcher.NewEqualsMatcher(clientSpanID),
							"timestamp": objmatcher.NewAnyMatcher(),
							"duration":  objmatcher.NewAnyMatcher(),
							"tags": objmatcher.MapMatcher(
								map[string]objmatcher.Matcher{
									"key0": objmatcher.NewEqualsMatcher("value0"),
								},
							),
						}),
					}),
				}),
			},
		},
		{
			Name: "trace.1 log entry with trace tag overwrite",
			SpanOptions: []wtracing.SpanOption{
				wtracing.WithSpanTag("key0", "value0"),
				wtracing.WithSpanTag("key1", "value1"),
				wtracing.WithSpanTag("key0", "value2"),
			},
			JSONMatcher: map[string]objmatcher.Matcher{
				"type":          objmatcher.NewEqualsMatcher("wrapped.1"),
				"entityName":    objmatcher.NewEqualsMatcher(entityName),
				"entityVersion": objmatcher.NewEqualsMatcher(entityVersion),
				"payload": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"type": objmatcher.NewEqualsMatcher("traceLogV1"),
					"traceLogV1": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"type": objmatcher.NewEqualsMatcher("trace.1"),
						"time": objmatcher.NewRegExpMatcher(".+"),
						"span": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"name":      objmatcher.NewEqualsMatcher("testOp"),
							"traceId":   objmatcher.NewEqualsMatcher(traceID),
							"id":        objmatcher.NewRegExpMatcher("[a-f0-9]+"),
							"parentId":  objmatcher.NewEqualsMatcher(clientSpanID),
							"timestamp": objmatcher.NewAnyMatcher(),
							"duration":  objmatcher.NewAnyMatcher(),
							"tags": objmatcher.MapMatcher(
								map[string]objmatcher.Matcher{
									"key0": objmatcher.NewEqualsMatcher("value2"),
									"key1": objmatcher.NewEqualsMatcher("value1"),
								},
							),
						}),
					}),
				}),
			},
		},
	}
}

func durationFormatOutputTest(t *testing.T, loggerProvider func(w io.Writer) trc1log.Logger) {
	buf := &bytes.Buffer{}
	logger := loggerProvider(buf)
	tracer, err := wzipkin.NewTracer(logger)
	require.NoError(t, err)
	span := tracer.StartSpan("testOp")
	time.Sleep(100 * time.Millisecond)
	// Finish() triggers logging
	span.Finish()

	entries, err := logreader.EntriesFromContent(buf.Bytes())
	require.NoError(t, err)
	require.Equal(t, 1, len(entries), "trace log should have exactly 1 entry")
	// Ensure the duration matches the sleep amount
	intValue := getDurationValue(t, entries[0])
	assert.True(t, intValue*time.Microsecond < 200*time.Millisecond, fmt.Sprintf("duration must be less than 200 milliseconds and is %d", intValue))
	assert.True(t, intValue*time.Microsecond > 100*time.Millisecond, fmt.Sprintf("duration must be more than 100 milliseconds and is %d", intValue))
}

func trc1LogJSONOutputTests(t *testing.T, entityName string, entityVersion string, loggerProvider func(w io.Writer) trc1log.Logger) {
	tracer, err := wzipkin.NewTracer(wtracing.NewNoopReporter())
	require.NoError(t, err)
	clientSpan := tracer.StartSpan("testOp", wtracing.WithKind(wtracing.Client))
	defer clientSpan.Finish()

	for i, tc := range Trc1TestCases(entityName, entityVersion, clientSpan) {
		t.Run(tc.Name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := loggerProvider(buf)

			tracer, err := wzipkin.NewTracer(
				logger,
				wtracing.WithLocalEndpoint(&wtracing.Endpoint{
					ServiceName: "testService",
					IPv4:        net.IPv4(127, 0, 0, 1),
					Port:        1234,
				}),
			)
			require.NoError(t, err)
			span := tracer.StartSpan("testOp", append([]wtracing.SpanOption{wtracing.WithParent(clientSpan)}, tc.SpanOptions...)...)
			// Finish() triggers logging
			span.Finish()

			entries, err := logreader.EntriesFromContent(buf.Bytes())
			require.NoError(t, err)
			require.Equal(t, 1, len(entries), "trace log should have exactly 1 entry")
			assert.NoError(t, tc.JSONMatcher.Matches(map[string]interface{}(entries[0])), "Case %d: %s\n%v", i, tc.Name, err)
		})
	}
}

func getDurationValue(t *testing.T, entry logreader.Entry) time.Duration {
	payload, ok := entry["payload"]
	require.True(t, ok)
	payloadAsMap, ok := payload.(map[string]interface{})
	require.True(t, ok)
	traceLog, ok := payloadAsMap["traceLogV1"]
	require.True(t, ok)
	traceLogAsMap, ok := traceLog.(map[string]interface{})
	require.True(t, ok)
	value, ok := traceLogAsMap["span"]
	require.True(t, ok)
	valueAsMap, ok := value.(map[string]interface{})
	require.True(t, ok)
	durationValue, ok := valueAsMap["duration"]
	require.True(t, ok)
	durationAsJSONNumber, ok := durationValue.(json.Number)
	require.True(t, ok)
	intValue, err := durationAsJSONNumber.Int64()
	require.NoError(t, err)
	return time.Duration(intValue)
}
