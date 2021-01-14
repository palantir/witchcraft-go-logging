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

func TestTrace1Logs(t *testing.T) {
	RunLogTests(t, []LogTest{
		{
			name: "Trace log",
			input: []string{
				`{"type":"trace.1","time":"2017-05-08T21:34:03.571Z","span":{"traceId":"3b2ecfbb0eaf8640","id":"630591bf0eaf799c","name":"operation","parentId":null,"timestamp":1491518454199000,"duration":1079,"annotations":[{"timestamp":1491518454199000,"value":"lc","endpoint":{"serviceName":"serviceName","ipv4":"10.160.121.155"}}]}}`,
			},
			output: []string{
				`[2017-05-08T21:34:03.571Z] traceId: 3b2ecfbb0eaf8640 id: 630591bf0eaf799c name: operation duration: 1079 microseconds`,
			},
		},
		{
			name: "Align output lines",
			input: []string{
				`{"type":"trace.1","time":"2017-05-08T21:34:03.5Z","span":{"traceId":"3b2ecfbb0eaf8640","id":"630591bf0eaf799c","name":"operation","parentId":null,"timestamp":1491518454199000,"duration":1079,"annotations":[{"timestamp":1491518454199000,"value":"lc","endpoint":{"serviceName":"serviceName","ipv4":"10.160.121.155"}}]}}`,
				`{"type":"trace.1","time":"2017-05-08T21:34:03.57Z","span":{"traceId":"3b2ecfbb0eaf8640","id":"630591bf0eaf799c","name":"operation","parentId":null,"timestamp":1491518454199000,"duration":1079,"annotations":[{"timestamp":1491518454199000,"value":"lc","endpoint":{"serviceName":"serviceName","ipv4":"10.160.121.155"}}]}}`,
				`{"type":"trace.1","time":"2017-05-08T21:34:03.571Z","span":{"traceId":"3b2ecfbb0eaf8640","id":"630591bf0eaf799c","name":"operation","parentId":null,"timestamp":1491518454199000,"duration":1079,"annotations":[{"timestamp":1491518454199000,"value":"lc","endpoint":{"serviceName":"serviceName","ipv4":"10.160.121.155"}}]}}`,
			},
			output: []string{
				`[2017-05-08T21:34:03.5Z]   traceId: 3b2ecfbb0eaf8640 id: 630591bf0eaf799c name: operation duration: 1079 microseconds`,
				`[2017-05-08T21:34:03.57Z]  traceId: 3b2ecfbb0eaf8640 id: 630591bf0eaf799c name: operation duration: 1079 microseconds`,
				`[2017-05-08T21:34:03.571Z] traceId: 3b2ecfbb0eaf8640 id: 630591bf0eaf799c name: operation duration: 1079 microseconds`,
			},
		},
	})
}
