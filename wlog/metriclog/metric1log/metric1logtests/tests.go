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
	//jsonLoggerUpdateTest(t, loggerProvider)
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

//func jsonLoggerUpdateTest(t *testing.T, loggerProvider func(params wlog.LoggerParams, origin string) svc1log.Logger) {
//	t.Run("update JSON logger", func(t *testing.T) {
//		currCase := TestCases()[0]
//
//		buf := bytes.Buffer{}
//		logger := loggerProvider(wlog.LoggerParams{
//			Level:  wlog.ErrorLevel,
//			Output: &buf,
//		}, currCase.Origin)
//
//		// log at info level
//		logger.Info(currCase.Message, currCase.LogParams...)
//
//		// output should be empty
//		assert.Equal(t, "", buf.String())
//
//		// update configuration to log at info level
//		updatable, ok := logger.(wlog.UpdatableLogger)
//		require.True(t, ok, "logger does not support updating")
//
//		updated := updatable.UpdateLogger(wlog.LoggerParams{
//			Level:  wlog.InfoLevel,
//			Output: &buf,
//		})
//		assert.True(t, updated)
//
//		// log at info level
//		logger.Info(currCase.Message, currCase.LogParams...)
//
//		// output should exist and match
//		gotServiceLog := map[string]interface{}{}
//		logEntry := buf.Bytes()
//		err := safejson.Unmarshal(logEntry, &gotServiceLog)
//		require.NoError(t, err, "Service log line is not a valid map: %v", string(logEntry))
//
//		assert.NoError(t, currCase.JSONMatcher.Matches(gotServiceLog), "No match")
//	})
//}
