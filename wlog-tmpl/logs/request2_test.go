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

func TestRequest2Logs(t *testing.T) {
	RunLogTests(t, []LogTest{
		{
			name: "GET request",
			input: []string{
				`{"type":"request.2","time":"2017-05-08T21:34:03.571Z","method":"GET","protocol":"HTTP/2.0","path":"/long","params":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":0,"responseSize":108,"duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"unsafeKey":"unsafeVal"}}`,
			},
			output: []string{
				`[2017-05-08T21:34:03.571Z] "GET /long HTTP/2.0" 200 108 68000`,
			},
		},
		{
			name: "Substitute path parameters from params",
			input: []string{
				`{"type":"request.2","time":"2017-05-08T21:34:03.571Z","method":"GET","protocol":"HTTP/2.0","path":"/long/{id}","params":{"id":25,"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":0,"responseSize":108,"duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"unsafeKey":"unsafeVal"}}`,
			},
			output: []string{
				`[2017-05-08T21:34:03.571Z] "GET /long/25 HTTP/2.0" 200 108 68000`,
			},
		},
		{
			name: "Substitute path parameters from unsafeParams",
			input: []string{
				`{"type":"request.2","time":"2017-05-08T21:34:03.571Z","method":"GET","protocol":"HTTP/2.0","path":"/long/{id}","params":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":0,"responseSize":108,"duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"unsafeKey":"unsafeVal","id":"unsafeId"}}`,
			},
			output: []string{
				`[2017-05-08T21:34:03.571Z] "GET /long/unsafeId HTTP/2.0" 200 108 68000`,
			},
		},
		{
			name: "Substitute path parameters with regexps",
			input: []string{
				`{"type":"request.2","time":"2017-05-08T21:34:03.571Z","method":"GET","protocol":"HTTP/2.0","path":"/long/{id:.+}/{path*}","params":{"id":"myID","path":"the/path","accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":0,"responseSize":108,"duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"unsafeKey":"unsafeVal"}}`,
			},
			output: []string{
				`[2017-05-08T21:34:03.571Z] "GET /long/myID/the/path HTTP/2.0" 200 108 68000`,
			},
		},
		{
			name: "Path parameter unsubstituted when no match exists",
			input: []string{
				`{"type":"request.2","time":"2017-05-08T21:34:03.571Z","method":"GET","protocol":"HTTP/2.0","path":"/long/{id:.+}/{path*}","params":{"id":"myID","accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":0,"responseSize":108,"duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"unsafeKey":"unsafeVal"}}`,
			},
			output: []string{
				`[2017-05-08T21:34:03.571Z] "GET /long/myID/{path*} HTTP/2.0" 200 108 68000`,
			},
		},
		{
			name: "Request without method",
			input: []string{
				`{"type":"request.2","time":"2017-05-08T21:34:03.571Z","protocol":"HTTP/2.0","path":"/long","params":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":0,"responseSize":108,"duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"unsafeKey":"unsafeVal"}}`,
			},
			output: []string{
				`[2017-05-08T21:34:03.571Z] "/long HTTP/2.0" 200 108 68000`,
			},
		},
		{
			name: "Align output lines",
			input: []string{
				`{"type":"request.2","time":"2017-05-08T21:34:03.5Z","method":"GET","protocol":"HTTP/2.0","path":"/long","params":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":0,"responseSize":108,"duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"unsafeKey":"unsafeVal"}}`,
				`{"type":"request.2","time":"2017-05-08T21:34:03.57Z","method":"GET","protocol":"HTTP/2.0","path":"/long","params":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":0,"responseSize":108,"duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"unsafeKey":"unsafeVal"}}`,
				`{"type":"request.2","time":"2017-05-08T21:34:03.571Z","method":"GET","protocol":"HTTP/2.0","path":"/long","params":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":0,"responseSize":108,"duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"unsafeKey":"unsafeVal"}}`,
				`{"type":"request.2","time":"2017-05-08T21:34:03.5714Z","method":"GET","protocol":"HTTP/2.0","path":"/long","params":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":0,"responseSize":108,"duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"unsafeKey":"unsafeVal"}}`,
				`{"type":"request.2","time":"2017-05-08T21:34:03.57148Z","method":"GET","protocol":"HTTP/2.0","path":"/long","params":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":0,"responseSize":108,"duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"unsafeKey":"unsafeVal"}}`,
				`{"type":"request.2","time":"2017-05-08T21:34:03.571485Z","method":"GET","protocol":"HTTP/2.0","path":"/long","params":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":0,"responseSize":108,"duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"unsafeKey":"unsafeVal"}}`,
				`{"type":"request.2","time":"2017-05-08T21:34:03.5714850Z","method":"GET","protocol":"HTTP/2.0","path":"/long","params":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":0,"responseSize":108,"duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"unsafeKey":"unsafeVal"}}`,
				`{"type":"request.2","time":"2017-05-08T21:34:03.57148509Z","method":"GET","protocol":"HTTP/2.0","path":"/long","params":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":0,"responseSize":108,"duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"unsafeKey":"unsafeVal"}}`,
				`{"type":"request.2","time":"2017-05-08T21:34:03.571485095Z","method":"GET","protocol":"HTTP/2.0","path":"/long","params":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":0,"responseSize":108,"duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"unsafeKey":"unsafeVal"}}`,
			},
			output: []string{
				`[2017-05-08T21:34:03.5Z]   "GET /long HTTP/2.0" 200 108 68000`,
				`[2017-05-08T21:34:03.57Z]  "GET /long HTTP/2.0" 200 108 68000`,
				`[2017-05-08T21:34:03.571Z] "GET /long HTTP/2.0" 200 108 68000`,
				`[2017-05-08T21:34:03.5714Z]      "GET /long HTTP/2.0" 200 108 68000`,
				`[2017-05-08T21:34:03.57148Z]     "GET /long HTTP/2.0" 200 108 68000`,
				`[2017-05-08T21:34:03.571485Z]    "GET /long HTTP/2.0" 200 108 68000`,
				`[2017-05-08T21:34:03.571485Z]    "GET /long HTTP/2.0" 200 108 68000`,
				`[2017-05-08T21:34:03.57148509Z]  "GET /long HTTP/2.0" 200 108 68000`,
				`[2017-05-08T21:34:03.571485095Z] "GET /long HTTP/2.0" 200 108 68000`,
			},
		},
	})
}
