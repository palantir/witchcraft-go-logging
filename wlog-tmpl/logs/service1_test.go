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

func TestService1Logs(t *testing.T) {
	RunLogTests(t, []LogTest{
		{
			name: "Error level log",
			input: []string{
				`{"type":"service.1","time":"2017-04-12T17:41:07.744Z","level":"ERROR","message":"Error handling request {}","origin":"com.palantir.remoting2.servers.jersey.JsonExceptionMapper","thread":"qtp1360518503-16","params":{},"uid":null,"sid":null,"tokenId":null,"orgId":null,"traceId":"fa4f6a37ac662fbd","unsafeParams":{"0":"8df8ace6-a068-4094-a7ff-0273469302f5","throwableMessage":null}}`,
			},
			output: []string{
				`ERROR [2017-04-12T17:41:07.744Z] com.palantir.remoting2.servers.jersey.JsonExceptionMapper: Error handling request 8df8ace6-a068-4094-a7ff-0273469302f5 (0: 8df8ace6-a068-4094-a7ff-0273469302f5, throwableMessage: <nil>)`,
			},
		},
		{
			name: "Ordering of printed maps is deterministic",
			input: []string{
				`{"type":"service.1","level":"ERROR","unsafeParams":{"b":"b","a":"a","c":"c","e":"e","d":"d"}}`,
			},
			output: []string{
				`ERROR [0001-01-01T00:00:00Z]      (a: a, b: b, c: c, d: d, e: e)`,
			},
		},
		{
			name: "No param substitution - named safe param",
			input: []string{
				`{"type":"service.1","time":"2017-05-25T18:49:10.652Z","level":"INFO","message":"Special node for '{node}' already exists","origin":"com.palantir.example.NodeCreator","params":{"node":"my-special-node"},"unsafeParams":{}}`,
			},
			output: []string{
				`INFO  [2017-05-25T18:49:10.652Z] com.palantir.example.NodeCreator: Special node for '{node}' already exists (node: my-special-node)`,
			},
		},
		{
			name: "No param substitution - named unsafe param",
			input: []string{
				`{"type":"service.1","time":"2017-05-25T18:49:10.652Z","level":"INFO","message":"Special node for '{node}' already exists","origin":"com.palantir.example.NodeCreator","params":{},"unsafeParams":{"node":"my-special-node"}}`,
			},
			output: []string{
				`INFO  [2017-05-25T18:49:10.652Z] com.palantir.example.NodeCreator: Special node for '{node}' already exists (node: my-special-node)`,
			},
		},
		{
			name: "No param substitution - unlabeled safe param",
			input: []string{
				`{"type":"service.1","time":"2017-05-25T18:49:10.652Z","level":"INFO","message":"Special node for '{}' already exists","origin":"com.palantir.example.NodeCreator","params":{"node":"my-special-node"},"unsafeParams":{}}`,
			},
			output: []string{
				`INFO  [2017-05-25T18:49:10.652Z] com.palantir.example.NodeCreator: Special node for '{}' already exists (node: my-special-node)`,
			},
		},
		{
			name: "No param substitution - unlabeled unsafe param",
			input: []string{
				`{"type":"service.1","time":"2017-05-25T18:49:10.652Z","level":"INFO","message":"Special node for '{}' already exists","origin":"com.palantir.example.NodeCreator","params":{},"unsafeParams":{"node":"my-special-node"}}`,
			},
			output: []string{
				`INFO  [2017-05-25T18:49:10.652Z] com.palantir.example.NodeCreator: Special node for '{}' already exists (node: my-special-node)`,
			},
		},
		{
			name: "No param substitution - positional safe param",
			input: []string{
				`{"type":"service.1","time":"2017-05-25T18:49:10.652Z","level":"INFO","message":"Special node for '{}' already exists","origin":"com.palantir.example.NodeCreator","params":{"node":"my-special-node"},"unsafeParams":{}}`,
			},
			output: []string{
				`INFO  [2017-05-25T18:49:10.652Z] com.palantir.example.NodeCreator: Special node for '{}' already exists (node: my-special-node)`,
			},
		},
		{
			name: "Message param substitution - positional unsafe param",
			input: []string{
				`{"type":"service.1","time":"2017-05-25T18:49:10.652Z","level":"INFO","message":"Special node for '{}' already exists","origin":"com.palantir.example.NodeCreator","params":{},"unsafeParams":{"0":"my-special-node"}}`,
			},
			output: []string{
				`INFO  [2017-05-25T18:49:10.652Z] com.palantir.example.NodeCreator: Special node for 'my-special-node' already exists (0: my-special-node)`,
			},
		},
		{
			name: "Message param substitution - nested positional param",
			input: []string{
				`{"type":"service.1","time":"2017-05-25T18:49:10.652Z","level":"INFO","message":"Special node for '{{{}}' already exists","origin":"com.palantir.example.NodeCreator","params":{},"unsafeParams":{"0":"my-special-node"}}`,
			},
			output: []string{
				`INFO  [2017-05-25T18:49:10.652Z] com.palantir.example.NodeCreator: Special node for '{{my-special-node}' already exists (0: my-special-node)`,
			},
		},
		{
			name: "Stacktrace and message param substitution - stack trace substitutes from safe param",
			input: []string{
				`{"type":"service.1","time":"2017-04-12T17:41:07.744Z","level":"ERROR","message":"Error handling request {}, safe: {}","origin":"com.palantir.remoting2.servers.jersey.JsonExceptionMapper","thread":"qtp1360518503-16","params":{"request": "/foo","throwableMessage":"Message"},"uid":null,"sid":null,"tokenId":null,"orgId":null,"traceId":"fa4f6a37ac662fbd","stacktrace":"java.lang.NullPointerException: {throwableMessage}\n\tcom.palantir.edu.profiles.resource.ProfileResource.getUserId(ProfileResource.java:36)\n","unsafeParams":{"0":"8df8ace6-a068-4094-a7ff-0273469302f5"}}`,
			},
			output: []string{
				`ERROR [2017-04-12T17:41:07.744Z] com.palantir.remoting2.servers.jersey.JsonExceptionMapper: Error handling request 8df8ace6-a068-4094-a7ff-0273469302f5, safe: {} (request: /foo, throwableMessage: Message) (0: 8df8ace6-a068-4094-a7ff-0273469302f5)
java.lang.NullPointerException: Message
	com.palantir.edu.profiles.resource.ProfileResource.getUserId(ProfileResource.java:36)
`,
			},
		},
		{
			name: "Stacktrace and message param substitution - stack trace substitutes from unsafe param",
			input: []string{
				`{"type":"service.1","time":"2017-04-12T17:41:07.744Z","level":"ERROR","message":"Error handling request {}, safe: {}","origin":"com.palantir.remoting2.servers.jersey.JsonExceptionMapper","thread":"qtp1360518503-16","params":{"request": "/foo"},"uid":null,"sid":null,"tokenId":null,"orgId":null,"traceId":"fa4f6a37ac662fbd","stacktrace":"java.lang.NullPointerException: {throwableMessage}\n\tcom.palantir.edu.profiles.resource.ProfileResource.getUserId(ProfileResource.java:36)\n","unsafeParams":{"0":"8df8ace6-a068-4094-a7ff-0273469302f5","throwableMessage":"Message"}}`,
			},
			output: []string{
				`ERROR [2017-04-12T17:41:07.744Z] com.palantir.remoting2.servers.jersey.JsonExceptionMapper: Error handling request 8df8ace6-a068-4094-a7ff-0273469302f5, safe: {} (request: /foo) (0: 8df8ace6-a068-4094-a7ff-0273469302f5, throwableMessage: Message)
java.lang.NullPointerException: Message
	com.palantir.edu.profiles.resource.ProfileResource.getUserId(ProfileResource.java:36)
`,
			},
		},
		{
			name: "Stacktrace and message param substitution - if stack trace element matches safe and unsafe params, prefer safe",
			input: []string{
				`{"type":"service.1","time":"2017-04-12T17:41:07.744Z","level":"ERROR","message":"Error handling request {}, safe: {}","origin":"com.palantir.remoting2.servers.jersey.JsonExceptionMapper","thread":"qtp1360518503-16","params":{"request": "/foo","throwableMessage":"Safe message"},"uid":null,"sid":null,"tokenId":null,"orgId":null,"traceId":"fa4f6a37ac662fbd","stacktrace":"java.lang.NullPointerException: {throwableMessage}\n\tcom.palantir.edu.profiles.resource.ProfileResource.getUserId(ProfileResource.java:36)\n","unsafeParams":{"0":"8df8ace6-a068-4094-a7ff-0273469302f5","throwableMessage":"Unsafe message"}}`,
			},
			output: []string{
				`ERROR [2017-04-12T17:41:07.744Z] com.palantir.remoting2.servers.jersey.JsonExceptionMapper: Error handling request 8df8ace6-a068-4094-a7ff-0273469302f5, safe: {} (request: /foo, throwableMessage: Safe message) (0: 8df8ace6-a068-4094-a7ff-0273469302f5, throwableMessage: Unsafe message)
java.lang.NullPointerException: Safe message
	com.palantir.edu.profiles.resource.ProfileResource.getUserId(ProfileResource.java:36)
`,
			},
		},
		{
			name: "Stacktrace and message param substitution - nested param",
			input: []string{
				`{"type":"service.1","time":"2017-04-12T17:41:07.744Z","level":"ERROR","message":"Error handling request {}, safe: {}","origin":"com.palantir.remoting2.servers.jersey.JsonExceptionMapper","thread":"qtp1360518503-16","params":{"request": "/foo","throwableMessage":"Message"},"uid":null,"sid":null,"tokenId":null,"orgId":null,"traceId":"fa4f6a37ac662fbd","stacktrace":"java.lang.NullPointerException: {{{throwableMessage}}\n\tcom.palantir.edu.profiles.resource.ProfileResource.getUserId(ProfileResource.java:36)\n","unsafeParams":{"0":"8df8ace6-a068-4094-a7ff-0273469302f5"}}`,
			},
			output: []string{
				`ERROR [2017-04-12T17:41:07.744Z] com.palantir.remoting2.servers.jersey.JsonExceptionMapper: Error handling request 8df8ace6-a068-4094-a7ff-0273469302f5, safe: {} (request: /foo, throwableMessage: Message) (0: 8df8ace6-a068-4094-a7ff-0273469302f5)
java.lang.NullPointerException: {{Message}
	com.palantir.edu.profiles.resource.ProfileResource.getUserId(ProfileResource.java:36)
`,
			},
		},
		{
			name: "Align output lines",
			input: []string{
				`{"type":"service.1","time":"2017-05-25T18:49:10.6Z","level":"INFO","message":"Special node for '{node}' already exists","origin":"com.palantir.example.NodeCreator","params":{"node":"my-special-node"},"unsafeParams":{}}`,
				`{"type":"service.1","time":"2017-05-25T18:49:10.652Z","level":"DEBUG","message":"Special node for '{node}' already exists","origin":"com.palantir.example.NodeCreator","params":{"node":"my-special-node"},"unsafeParams":{}}`,
				`{"type":"service.1","time":"2017-05-25T18:49:10.652Z","level":"INFO","message":"Special node for '{node}' already exists","origin":"com.palantir.example.NodeCreator","params":{"node":"my-special-node"},"unsafeParams":{}}`,
				`{"type":"service.1","time":"2017-05-25T18:49:10.652Z","level":"WARN","message":"Special node for '{node}' already exists","origin":"com.palantir.example.NodeCreator","params":{"node":"my-special-node"},"unsafeParams":{}}`,
				`{"type":"service.1","time":"2017-05-25T18:49:10.652Z","level":"ERROR","message":"Special node for '{node}' already exists","origin":"com.palantir.example.NodeCreator","params":{"node":"my-special-node"},"unsafeParams":{}}`,
				`{"type":"service.1","time":"2017-05-25T18:49:10.652Z","level":"FATAL","message":"Special node for '{node}' already exists","origin":"com.palantir.example.NodeCreator","params":{"node":"my-special-node"},"unsafeParams":{}}`,
			},
			output: []string{
				`INFO  [2017-05-25T18:49:10.6Z]   com.palantir.example.NodeCreator: Special node for '{node}' already exists (node: my-special-node)`,
				`DEBUG [2017-05-25T18:49:10.652Z] com.palantir.example.NodeCreator: Special node for '{node}' already exists (node: my-special-node)`,
				`INFO  [2017-05-25T18:49:10.652Z] com.palantir.example.NodeCreator: Special node for '{node}' already exists (node: my-special-node)`,
				`WARN  [2017-05-25T18:49:10.652Z] com.palantir.example.NodeCreator: Special node for '{node}' already exists (node: my-special-node)`,
				`ERROR [2017-05-25T18:49:10.652Z] com.palantir.example.NodeCreator: Special node for '{node}' already exists (node: my-special-node)`,
				`FATAL [2017-05-25T18:49:10.652Z] com.palantir.example.NodeCreator: Special node for '{node}' already exists (node: my-special-node)`,
			},
		},
		{
			name: "Test formatting of large values",
			input: []string{
				`{"type":"service.1","time":"2017-05-25T18:49:10.6Z","level":"INFO","message":"Special node for '{node}' already exists","origin":"com.palantir.example.NodeCreator","params":{"node":"my-special-node"},"unsafeParams":{"0":38711478389}}`,
			},
			output: []string{
				`INFO  [2017-05-25T18:49:10.6Z]   com.palantir.example.NodeCreator: Special node for '{node}' already exists (node: my-special-node) (0: 38711478389)`,
			},
		},
		{
			name: "Test formatting of floating point values",
			input: []string{
				`{"type":"service.1","time":"2017-05-25T18:49:10.6Z","level":"INFO","message":"Special node for '{node}' already exists","origin":"com.palantir.example.NodeCreator","params":{"node":"my-special-node"},"unsafeParams":{"0":12345.8711478389}}`,
			},
			output: []string{
				`INFO  [2017-05-25T18:49:10.6Z]   com.palantir.example.NodeCreator: Special node for '{node}' already exists (node: my-special-node) (0: 12345.8711478389)`,
			},
		},
		{
			name: "Test wrapped service log",
			input: []string{
				`{"type":"wrapped.1","payload":{"type":"serviceLogV1","serviceLogV1":{"type":"service.1","time":"2017-04-12T17:41:07.744Z","level":"ERROR","message":"Error handling request {}","origin":"com.palantir.remoting2.servers.jersey.JsonExceptionMapper","thread":"qtp1360518503-16","params":{},"uid":null,"sid":null,"tokenId":null,"orgId":null,"traceId":"fa4f6a37ac662fbd","unsafeParams":{"0":"8df8ace6-a068-4094-a7ff-0273469302f5","throwableMessage":null}}},"entityName":"codex-hub","entityVersion":"v2.1.0"}`,
			},
			output: []string{
				`ERROR [2017-04-12T17:41:07.744Z] com.palantir.remoting2.servers.jersey.JsonExceptionMapper: Error handling request 8df8ace6-a068-4094-a7ff-0273469302f5 (0: 8df8ace6-a068-4094-a7ff-0273469302f5, throwableMessage: <nil>)`,
			},
		},
		{
			name: "Test service log with empty origin",
			input: []string{
				`{"type":"service.1","time":"2017-04-12T17:41:07.744Z","level":"ERROR","message":"Error handling request {}","params":{},"uid":null,"sid":null,"tokenId":null,"orgId":null,"traceId":"fa4f6a37ac662fbd","unsafeParams":{}}`,
			},
			output: []string{
				`ERROR [2017-04-12T17:41:07.744Z] Error handling request {}`,
			},
		},
	})
}
