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

func TestWrapped1Logs(t *testing.T) {
	RunLogTests(t, []LogTest{
		{
			name: "Test wrapped service.1 log",
			input: []string{
				`{"type":"wrapped.1","payload":{"type":"serviceLogV1","serviceLogV1":{"type":"service.1","time":"2017-04-12T17:41:07.744Z","level":"ERROR","message":"Error handling request {}","origin":"com.palantir.remoting2.servers.jersey.JsonExceptionMapper","thread":"qtp1360518503-16","params":{},"uid":null,"sid":null,"tokenId":null,"orgId":null,"traceId":"fa4f6a37ac662fbd","unsafeParams":{"0":"8df8ace6-a068-4094-a7ff-0273469302f5","throwableMessage":null}}},"entityName":"codex-hub","entityVersion":"v2.1.0"}`,
			},
			output: []string{
				`ERROR [2017-04-12T17:41:07.744Z] com.palantir.remoting2.servers.jersey.JsonExceptionMapper: Error handling request 8df8ace6-a068-4094-a7ff-0273469302f5 (0: 8df8ace6-a068-4094-a7ff-0273469302f5, throwableMessage: <nil>)`,
			},
		},
		{
			name: "Wrapped request.2 log",
			input: []string{
				`{"type":"wrapped.1","payload":{"type":"requestLogV2","requestLogV2":{"type":"request.2","time":"2017-05-08T21:34:03.571Z","method":"GET","protocol":"HTTP/2.0","path":"/long","params":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":0,"responseSize":108,"duration":68000,"uid":null,"sid":null,"tokenId":null,"orgId":null,"traceId":"d384560e01f9a657","unsafeParams":{"unsafeKey":"unsafeVal"}}},"entityName":"codex-hub","entityVersion":"v2.1.0"}`,
			},
			output: []string{
				`[2017-05-08T21:34:03.571Z] "GET /long HTTP/2.0" 200 108 68000`,
			},
		},
		{
			name: "Wrapped trace.1 log",
			input: []string{
				`{"type":"wrapped.1","payload":{"type":"traceLogV1","traceLogV1":{"type":"trace.1","time":"2017-05-08T21:34:03.571Z","span":{"traceId":"3b2ecfbb0eaf8640","id":"630591bf0eaf799c","name":"operation","parentId":null,"timestamp":1491518454199000,"duration":1079,"annotations":[{"timestamp":1491518454199000,"value":"lc","endpoint":{"serviceName":"serviceName","ipv4":"10.160.121.155"}}]}}},"entityName":"codex-hub","entityVersion":"v2.1.0"}`,
			},
			output: []string{
				`[2017-05-08T21:34:03.571Z] traceId: 3b2ecfbb0eaf8640 id: 630591bf0eaf799c name: operation duration: 1079 microseconds`,
			},
		},
	})
}
