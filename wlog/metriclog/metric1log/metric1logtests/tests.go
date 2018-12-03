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

package metric1logtests

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/palantir/pkg/objmatcher"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/witchcraft-go-logging/wlog/metriclog/metric1log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	Name         string
	MetricName   string
	MetricType   string
	Tags         map[string]string
	Values       map[string]interface{}
	UID          string
	SID          string
	TokenID      string
	UnsafeParams map[string]interface{}
	JSONMatcher  objmatcher.MapMatcher
}

func (tc TestCase) Params() []metric1log.Param {
	return []metric1log.Param{
		metric1log.Values(tc.Values),
		metric1log.UID(tc.UID),
		metric1log.SID(tc.SID),
		metric1log.TokenID(tc.TokenID),
		metric1log.Tags(tc.Tags),
		metric1log.UnsafeParams(tc.UnsafeParams),
	}
}

func TestCases() []TestCase {
	return []TestCase{
		{
			Name:       "basic metric log entry",
			MetricName: "com.palantir.deployability.logtrough.iteratorage.millis",
			MetricType: "histogram",
			UID:        "user-1",
			SID:        "session-1",
			Values: map[string]interface{}{
				"max":    1400,
				"mean":   70,
				"stddev": 20,
				"p95":    100,
				"p99":    1200,
				"p999":   1350,
				"count":  100,
			},
			Tags: map[string]string{
				"shardId": "shard-1234",
			},
			TokenID: "X-Y-Z",
			UnsafeParams: map[string]interface{}{
				"Password": "HelloWorld!",
			},
			JSONMatcher: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"metricName": objmatcher.NewEqualsMatcher("com.palantir.deployability.logtrough.iteratorage.millis"),
				"time":       objmatcher.NewRegExpMatcher(".+"),
				"type":       objmatcher.NewEqualsMatcher("metric.1"),
				"metricType": objmatcher.NewEqualsMatcher("histogram"),
				"uid":        objmatcher.NewEqualsMatcher("user-1"),
				"sid":        objmatcher.NewEqualsMatcher("session-1"),
				"values": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"max":    objmatcher.NewEqualsMatcher(json.Number("1400")),
					"mean":   objmatcher.NewEqualsMatcher(json.Number("70")),
					"stddev": objmatcher.NewEqualsMatcher(json.Number("20")),
					"p95":    objmatcher.NewEqualsMatcher(json.Number("100")),
					"p99":    objmatcher.NewEqualsMatcher(json.Number("1200")),
					"p999":   objmatcher.NewEqualsMatcher(json.Number("1350")),
					"count":  objmatcher.NewEqualsMatcher(json.Number("100")),
				}),
				"tags": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"shardId": objmatcher.NewEqualsMatcher("shard-1234"),
				}),
				"tokenId": objmatcher.NewEqualsMatcher("X-Y-Z"),
				"unsafeParams": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"Password": objmatcher.NewEqualsMatcher("HelloWorld!"),
				}),
			}),
		},
	}
}

func JSONTestSuite(t *testing.T, loggerProvider func(w io.Writer) metric1log.Logger) {
	jsonOutputTests(t, loggerProvider)
	valueIsntOverwrittenByValues(t, loggerProvider)
	extraValuesDoNotAppear(t, loggerProvider)
	tagIsntOverwrittenByTags(t, loggerProvider)
	extraTagsDoNotAppear(t, loggerProvider)
}

func jsonOutputTests(t *testing.T, loggerProvider func(w io.Writer) metric1log.Logger) {
	for i, tc := range TestCases() {
		t.Run(tc.Name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := loggerProvider(buf)

			logger.Metric(tc.MetricName, tc.MetricType, tc.Params()...)

			gotMetricLog := map[string]interface{}{}
			logEntry := buf.Bytes()
			err := safejson.Unmarshal(logEntry, &gotMetricLog)
			require.NoError(t, err, "Case %d: %s\nMetric log line is not a valid map: %v", i, tc.Name, string(logEntry))

			assert.NoError(t, tc.JSONMatcher.Matches(gotMetricLog), "Case %d: %s", i, tc.Name)
		})
	}
}

// Verifies that if different parameters are specified using Value and Values params, all of the values are present in
// the final output (that is, these parameters should be additive).
func valueIsntOverwrittenByValues(t *testing.T, loggerProvider func(w io.Writer) metric1log.Logger) {
	t.Run("Value and Values params are additive", func(t *testing.T) {
		var buf bytes.Buffer
		logger := loggerProvider(&buf)

		logger.Metric("metric", "metric-type", metric1log.Value("key", "value"), metric1log.Values(map[string]interface{}{"keys": "values"}))

		gotMetricLog := map[string]interface{}{}
		logEntry := buf.Bytes()
		err := safejson.Unmarshal(logEntry, &gotMetricLog)
		require.NoError(t, err, "Metric log line is not a valid map: %v", string(logEntry))
		assert.NoError(t, objmatcher.MapMatcher(map[string]objmatcher.Matcher{
			"metricName": objmatcher.NewEqualsMatcher("metric"),
			"time":       objmatcher.NewRegExpMatcher(".+"),
			"type":       objmatcher.NewEqualsMatcher("metric.1"),
			"metricType": objmatcher.NewEqualsMatcher("metric-type"),
			"values": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"key":  objmatcher.NewEqualsMatcher("value"),
				"keys": objmatcher.NewEqualsMatcher("values"),
			}),
		}).Matches(gotMetricLog))
	})
}

// Verifies that parameters remain separate between different logger calls (ensures there is not a bug where parameters
// are modified by making a logger call).
func extraValuesDoNotAppear(t *testing.T, loggerProvider func(w io.Writer) metric1log.Logger) {
	t.Run("Value and Values params stay separate across logger calls", func(t *testing.T) {
		var buf bytes.Buffer
		logger := loggerProvider(&buf)

		reusedParams := metric1log.Values(map[string]interface{}{"keys": "values"})
		logger.Metric("metric", "metric-type", reusedParams, metric1log.Value("key", "value"))

		gotMetricLog := map[string]interface{}{}
		logEntry := buf.Bytes()
		err := safejson.Unmarshal(logEntry, &gotMetricLog)
		require.NoError(t, err, "Metric log line is not a valid map: %v", string(logEntry))
		assert.NoError(t, objmatcher.MapMatcher(map[string]objmatcher.Matcher{
			"metricName": objmatcher.NewEqualsMatcher("metric"),
			"time":       objmatcher.NewRegExpMatcher(".+"),
			"type":       objmatcher.NewEqualsMatcher("metric.1"),
			"metricType": objmatcher.NewEqualsMatcher("metric-type"),
			"values": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"key":  objmatcher.NewEqualsMatcher("value"),
				"keys": objmatcher.NewEqualsMatcher("values"),
			}),
		}).Matches(gotMetricLog))

		buf.Reset()
		logger.Metric("metric", "metric-type", reusedParams)

		gotMetricLog = map[string]interface{}{}
		logEntry = buf.Bytes()
		err = safejson.Unmarshal(logEntry, &gotMetricLog)
		require.NoError(t, err, "Metric log line is not a valid map: %v", string(logEntry))
		assert.NoError(t, objmatcher.MapMatcher(map[string]objmatcher.Matcher{
			"metricName": objmatcher.NewEqualsMatcher("metric"),
			"time":       objmatcher.NewRegExpMatcher(".+"),
			"type":       objmatcher.NewEqualsMatcher("metric.1"),
			"metricType": objmatcher.NewEqualsMatcher("metric-type"),
			"values": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"keys": objmatcher.NewEqualsMatcher("values"),
			}),
		}).Matches(gotMetricLog))
	})
}

// Verifies that if different parameters are specified using Tag and Tags params, all of the tags are present in the
// final output (that is, these parameters should be additive).
func tagIsntOverwrittenByTags(t *testing.T, loggerProvider func(w io.Writer) metric1log.Logger) {
	t.Run("Tag and Tags params are additive", func(t *testing.T) {
		var buf bytes.Buffer
		logger := loggerProvider(&buf)

		logger.Metric("metric", "metric-type", metric1log.Tag("key", "value"), metric1log.Tags(map[string]string{"keys": "values"}))

		gotMetricLog := map[string]interface{}{}
		logEntry := buf.Bytes()
		err := safejson.Unmarshal(logEntry, &gotMetricLog)
		require.NoError(t, err, "Metric log line is not a valid map: %v", string(logEntry))
		assert.NoError(t, objmatcher.MapMatcher(map[string]objmatcher.Matcher{
			"metricName": objmatcher.NewEqualsMatcher("metric"),
			"time":       objmatcher.NewRegExpMatcher(".+"),
			"type":       objmatcher.NewEqualsMatcher("metric.1"),
			"metricType": objmatcher.NewEqualsMatcher("metric-type"),
			"tags": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"key":  objmatcher.NewEqualsMatcher("value"),
				"keys": objmatcher.NewEqualsMatcher("values"),
			}),
		}).Matches(gotMetricLog))
	})
}

// Verifies that parameters remain separate between different logger calls (ensures there is not a bug where parameters
// are modified by making a logger call).
func extraTagsDoNotAppear(t *testing.T, loggerProvider func(w io.Writer) metric1log.Logger) {
	t.Run("Tag and Tags params stay separate across logger calls", func(t *testing.T) {
		var buf bytes.Buffer
		logger := loggerProvider(&buf)

		reusedParams := metric1log.Tags(map[string]string{"keys": "values"})
		logger.Metric("metric", "metric-type", reusedParams, metric1log.Tag("key", "value"))

		gotMetricLog := map[string]interface{}{}
		logEntry := buf.Bytes()
		err := safejson.Unmarshal(logEntry, &gotMetricLog)
		require.NoError(t, err, "Metric log line is not a valid map: %v", string(logEntry))
		assert.NoError(t, objmatcher.MapMatcher(map[string]objmatcher.Matcher{
			"metricName": objmatcher.NewEqualsMatcher("metric"),
			"time":       objmatcher.NewRegExpMatcher(".+"),
			"type":       objmatcher.NewEqualsMatcher("metric.1"),
			"metricType": objmatcher.NewEqualsMatcher("metric-type"),
			"tags": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"key":  objmatcher.NewEqualsMatcher("value"),
				"keys": objmatcher.NewEqualsMatcher("values"),
			}),
		}).Matches(gotMetricLog))

		buf.Reset()
		logger.Metric("metric", "metric-type", reusedParams)

		gotMetricLog = map[string]interface{}{}
		logEntry = buf.Bytes()
		err = safejson.Unmarshal(logEntry, &gotMetricLog)
		require.NoError(t, err, "Metric log line is not a valid map: %v", string(logEntry))
		assert.NoError(t, objmatcher.MapMatcher(map[string]objmatcher.Matcher{
			"metricName": objmatcher.NewEqualsMatcher("metric"),
			"time":       objmatcher.NewRegExpMatcher(".+"),
			"type":       objmatcher.NewEqualsMatcher("metric.1"),
			"metricType": objmatcher.NewEqualsMatcher("metric-type"),
			"tags": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"keys": objmatcher.NewEqualsMatcher("values"),
			}),
		}).Matches(gotMetricLog))
	})
}
