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

func TestEvent2Logs(t *testing.T) {
	RunLogTests(t, []LogTest{
		{
			name: "Event log entry",
			input: []string{
				`{"type":"event.2","time":"2017-05-08T21:34:03.571Z","eventName":"com.palantir.lime.searchResult","values":{"resultIndex":2},"uid":null,"sid":null,"tokenId":null,"orgId":null,"traceId":null,"unsafeParams":{"searchTerm":"unsafeVal"}}`,
			},
			output: []string{
				`[2017-05-08T21:34:03.571Z] com.palantir.lime.searchResult (resultIndex: 2) (searchTerm: unsafeVal)`,
			},
		},
		{
			name: "Align output lines",
			input: []string{
				`{"type":"event.2","time":"2017-05-08T21:34:03.5Z","eventName":"com.palantir.lime.searchResult","values":{"resultIndex":2},"uid":null,"sid":null,"tokenId":null,"orgId":null,"traceId":null,"unsafeParams":{"searchTerm":"unsafeVal"}}`,
				`{"type":"event.2","time":"2017-05-08T21:34:03.57Z","eventName":"com.palantir.lime.searchResult","values":{"resultIndex":2},"uid":null,"sid":null,"tokenId":null,"orgId":null,"traceId":null,"unsafeParams":{"searchTerm":"unsafeVal"}}`,
				`{"type":"event.2","time":"2017-05-08T21:34:03.571Z","eventName":"com.palantir.lime.searchResult","values":{"resultIndex":2},"uid":null,"sid":null,"tokenId":null,"orgId":null,"traceId":null,"unsafeParams":{"searchTerm":"unsafeVal"}}`,
			},
			output: []string{
				`[2017-05-08T21:34:03.5Z]   com.palantir.lime.searchResult (resultIndex: 2) (searchTerm: unsafeVal)`,
				`[2017-05-08T21:34:03.57Z]  com.palantir.lime.searchResult (resultIndex: 2) (searchTerm: unsafeVal)`,
				`[2017-05-08T21:34:03.571Z] com.palantir.lime.searchResult (resultIndex: 2) (searchTerm: unsafeVal)`,
			},
		},
	})
}
