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

package diag1log_test

import (
	"testing"

	"github.com/palantir/witchcraft-go-logging/conjure/sls/spec/logging"
	"github.com/palantir/witchcraft-go-logging/internal/conjuretype"
	"github.com/palantir/witchcraft-go-logging/wlog/diaglog/diag1log"
	"github.com/stretchr/testify/require"
)

func TestThreadDumpV1FromGoroutines(t *testing.T) {
	for _, test := range []struct {
		Name     string
		Input    string
		Expected logging.ThreadDumpV1
	}{
		{
			Name: "single goroutine",
			Input: `goroutine 14 [select]:
net/http.(*persistConn).writeLoop(0xc0000bd0e0)
	/usr/local/Cellar/go/1.11.2/libexec/src/net/http/transport.go:1885 +0x113
created by net/http.(*Transport).dialConn
	/usr/local/Cellar/go/1.11.2/libexec/src/net/http/transport.go:1339 +0x966
`,
			Expected: logging.ThreadDumpV1{
				Threads: []logging.ThreadInfoV1{
					{
						Name:   strPtr("goroutine 14 [select]"),
						Id:     safelongPtr(14),
						Params: map[string]interface{}{"status": "select"},
						StackTrace: []logging.StackFrameV1{
							{
								Address:   strPtr("0x113"),
								Procedure: strPtr("net/http.(*persistConn).writeLoop"),
								File:      strPtr("net/http/transport.go"),
								Line:      intPtr(1885),
								Params:    map[string]interface{}{},
							},
							{
								Address:   strPtr("0x966"),
								Procedure: strPtr("net/http.(*Transport).dialConn"),
								File:      strPtr("net/http/transport.go"),
								Line:      intPtr(1339),
								Params: map[string]interface{}{
									"goroutineCreator": true,
								},
							},
						},
					},
				},
			},
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			dump := diag1log.ThreadDumpV1FromGoroutines([]byte(test.Input))
			require.Equal(t, test.Expected, dump)
		})
	}
}

func strPtr(s string) *string { return &s }

func intPtr(i int) *int { return &i }

func safelongPtr(i int64) *conjuretype.SafeLong {
	s, _ := conjuretype.NewSafeLong(i)
	return &s
}
