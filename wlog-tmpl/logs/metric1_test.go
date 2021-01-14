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

package logs_test

import (
	"testing"
)

func TestMetric1Logs(t *testing.T) {
	RunLogTests(t, []LogTest{
		{
			name: "Metric log entry",
			input: []string{
				`{"type":"metric.1","time":"2017-05-08T21:34:03.571Z","metricName":"com.palantir.lime.searchResult","metricType":"gauge","values":{"resultIndex":2},"tags":{"deployment":"yellow", "env":"prod"},"uid":null,"sid":null,"tokenId":null,"unsafeParams":{"searchTerm":"unsafeVal"}}`,
			},
			output: []string{
				`[2017-05-08T21:34:03.571Z] METRIC com.palantir.lime.searchResult gauge (resultIndex: 2) (deployment: yellow, env: prod) (searchTerm: unsafeVal)`,
			},
		},
		{
			name: "Align output lines",
			input: []string{
				`{"type":"metric.1","time":"2017-05-08T21:34:03.5Z","metricName":"com.palantir.lime.searchResult","metricType":"gauge","values":{"resultIndex":2},"tags":{"deployment":"yellow", "env":"prod"},"uid":null,"sid":null,"tokenId":null,"unsafeParams":{"searchTerm":"unsafeVal"}}`,
				`{"type":"metric.1","time":"2017-05-08T21:34:03.57Z","metricName":"com.palantir.lime.searchResult","metricType":"gauge","values":{"resultIndex":2},"tags":{"deployment":"yellow", "env":"prod"},"uid":null,"sid":null,"tokenId":null,"unsafeParams":{"searchTerm":"unsafeVal"}}`,
				`{"type":"metric.1","time":"2017-05-08T21:34:03.571Z","metricName":"com.palantir.lime.searchResult","metricType":"gauge","values":{"resultIndex":2},"tags":{"deployment":"yellow", "env":"prod"},"uid":null,"sid":null,"tokenId":null,"unsafeParams":{"searchTerm":"unsafeVal"}}`,
			},
			output: []string{
				`[2017-05-08T21:34:03.5Z]   METRIC com.palantir.lime.searchResult gauge (resultIndex: 2) (deployment: yellow, env: prod) (searchTerm: unsafeVal)`,
				`[2017-05-08T21:34:03.57Z]  METRIC com.palantir.lime.searchResult gauge (resultIndex: 2) (deployment: yellow, env: prod) (searchTerm: unsafeVal)`,
				`[2017-05-08T21:34:03.571Z] METRIC com.palantir.lime.searchResult gauge (resultIndex: 2) (deployment: yellow, env: prod) (searchTerm: unsafeVal)`,
			},
		},
	})
}
