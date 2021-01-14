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

func TestRequest1Logs(t *testing.T) {
	RunLogTests(t, []LogTest{
		{
			name: "GET request",
			input: []string{
				`{"type":"request.1","time":"2017-05-08T21:34:03.571Z","method":"GET","protocol":"HTTP/2.0","path":"/long","pathParams":{},"queryParams":{},"headerParams":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":"0","responseSize":"0","duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"pathParams":{},"queryParams":{},"headerParams":{}}}`,
			},
			output: []string{
				`[2017-05-08T21:34:03.571Z] "GET /long HTTP/2.0" 200 0 68000`,
			},
		},
		{
			name: "Substitute path parameters from params",
			input: []string{
				`{"type":"request.1","time":"2017-05-08T21:34:03.571Z","method":"GET","protocol":"HTTP/2.0","path":"/long/{id}","pathParams":{"id":25},"queryParams":{},"headerParams":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":"0","responseSize":"0","duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"pathParams":{},"queryParams":{},"headerParams":{}}}`,
			},
			output: []string{
				`[2017-05-08T21:34:03.571Z] "GET /long/25 HTTP/2.0" 200 0 68000`,
			},
		},
		{
			name: "Substitute path parameters from unsafeParams",
			input: []string{
				`{"type":"request.1","time":"2017-05-08T21:34:03.571Z","method":"GET","protocol":"HTTP/2.0","path":"/long/{id}","pathParams":{},"queryParams":{},"headerParams":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":"0","responseSize":"0","duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"id":"unsafeId"}}`,
			},
			output: []string{
				`[2017-05-08T21:34:03.571Z] "GET /long/unsafeId HTTP/2.0" 200 0 68000`,
			},
		},
		{
			name: "Substitute path parameters with regexps",
			input: []string{
				`{"type":"request.1","time":"2017-05-08T21:34:03.571Z","method":"GET","protocol":"HTTP/2.0","path":"/long/{id:.+}/{path*}","pathParams":{"id":"myID","path":"the/path"},"queryParams":{},"headerParams":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":"0","responseSize":"0","duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{}}`,
			},
			output: []string{
				`[2017-05-08T21:34:03.571Z] "GET /long/myID/the/path HTTP/2.0" 200 0 68000`,
			},
		},
		{
			name: "Path parameter unsubstituted when no match exists",
			input: []string{
				`{"type":"request.1","time":"2017-05-08T21:34:03.571Z","method":"GET","protocol":"HTTP/2.0","path":"/long/{id:.+}/{path*}","pathParams":{"id":"myID"},"queryParams":{},"headerParams":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":"0","responseSize":"0","duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{}}`,
			},
			output: []string{
				`[2017-05-08T21:34:03.571Z] "GET /long/myID/{path*} HTTP/2.0" 200 0 68000`,
			},
		},
		{
			name: "Request without method",
			input: []string{
				`{"type":"request.1","time":"2017-05-08T21:34:03.571Z","protocol":"HTTP/2.0","path":"/long","pathParams":{},"queryParams":{},"headerParams":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":"0","responseSize":"0","duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"pathParams":{},"queryParams":{},"headerParams":{}}}`,
			},
			output: []string{
				`[2017-05-08T21:34:03.571Z] "/long HTTP/2.0" 200 0 68000`,
			},
		},
		{
			name: "Align output lines",
			input: []string{
				`{"type":"request.1","method":"GET","time":"2017-05-08T21:34:03.5Z","protocol":"HTTP/2.0","path":"/long","pathParams":{},"queryParams":{},"headerParams":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":"0","responseSize":"0","duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"pathParams":{},"queryParams":{},"headerParams":{}}}`,
				`{"type":"request.1","method":"GET","time":"2017-05-08T21:34:03.57Z","protocol":"HTTP/2.0","path":"/long","pathParams":{},"queryParams":{},"headerParams":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":"0","responseSize":"0","duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"pathParams":{},"queryParams":{},"headerParams":{}}}`,
				`{"type":"request.1","method":"GET","time":"2017-05-08T21:34:03.571Z","protocol":"HTTP/2.0","path":"/long","pathParams":{},"queryParams":{},"headerParams":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":"0","responseSize":"0","duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"pathParams":{},"queryParams":{},"headerParams":{}}}`,
				`{"type":"request.1","method":"GET","time":"2017-05-08T21:34:03.5714Z","protocol":"HTTP/2.0","path":"/long","pathParams":{},"queryParams":{},"headerParams":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":"0","responseSize":"0","duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"pathParams":{},"queryParams":{},"headerParams":{}}}`,
				`{"type":"request.1","method":"GET","time":"2017-05-08T21:34:03.57148Z","protocol":"HTTP/2.0","path":"/long","pathParams":{},"queryParams":{},"headerParams":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":"0","responseSize":"0","duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"pathParams":{},"queryParams":{},"headerParams":{}}}`,
				`{"type":"request.1","method":"GET","time":"2017-05-08T21:34:03.571485Z","protocol":"HTTP/2.0","path":"/long","pathParams":{},"queryParams":{},"headerParams":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":"0","responseSize":"0","duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"pathParams":{},"queryParams":{},"headerParams":{}}}`,
				`{"type":"request.1","method":"GET","time":"2017-05-08T21:34:03.5714850Z","protocol":"HTTP/2.0","path":"/long","pathParams":{},"queryParams":{},"headerParams":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":"0","responseSize":"0","duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"pathParams":{},"queryParams":{},"headerParams":{}}}`,
				`{"type":"request.1","method":"GET","time":"2017-05-08T21:34:03.57148509Z","protocol":"HTTP/2.0","path":"/long","pathParams":{},"queryParams":{},"headerParams":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":"0","responseSize":"0","duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"pathParams":{},"queryParams":{},"headerParams":{}}}`,
				`{"type":"request.1","method":"GET","time":"2017-05-08T21:34:03.571485095Z","protocol":"HTTP/2.0","path":"/long","pathParams":{},"queryParams":{},"headerParams":{"accept-encoding":"gzip","host":"localhost:8443","user-agent":"okhttp/3.5.0"},"status":200,"requestSize":"0","responseSize":"0","duration":68000,"uid":null,"sid":null,"tokenId":null,"traceId":"d384560e01f9a657","unsafeParams":{"pathParams":{},"queryParams":{},"headerParams":{}}}`,
			},
			output: []string{
				`[2017-05-08T21:34:03.5Z]   "GET /long HTTP/2.0" 200 0 68000`,
				`[2017-05-08T21:34:03.57Z]  "GET /long HTTP/2.0" 200 0 68000`,
				`[2017-05-08T21:34:03.571Z] "GET /long HTTP/2.0" 200 0 68000`,
				`[2017-05-08T21:34:03.5714Z]      "GET /long HTTP/2.0" 200 0 68000`,
				`[2017-05-08T21:34:03.57148Z]     "GET /long HTTP/2.0" 200 0 68000`,
				`[2017-05-08T21:34:03.571485Z]    "GET /long HTTP/2.0" 200 0 68000`,
				`[2017-05-08T21:34:03.571485Z]    "GET /long HTTP/2.0" 200 0 68000`,
				`[2017-05-08T21:34:03.57148509Z]  "GET /long HTTP/2.0" 200 0 68000`,
				`[2017-05-08T21:34:03.571485095Z] "GET /long HTTP/2.0" 200 0 68000`,
			},
		},
	})
}
