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

package diag1logtests

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/palantir/pkg/objmatcher"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/witchcraft-go-logging/conjure/sls/spec/logging"
	"github.com/palantir/witchcraft-go-logging/internal/conjuretype"
	"github.com/palantir/witchcraft-go-logging/wlog/diaglog/diag1log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	Name         string
	Diagnostic   logging.Diagnostic
	UnsafeParams map[string]interface{}
	JSONMatcher  objmatcher.MapMatcher
}

func TestCases() []TestCase {
	return []TestCase{
		{
			Name: "generic diagnostic log entry",
			Diagnostic: logging.NewDiagnosticFromGeneric(logging.GenericDiagnostic{
				DiagnosticType: "DIAG_TYPE",
				Value: map[string]string{
					"testKey": "test_value",
				},
			}),
			UnsafeParams: map[string]interface{}{
				"Password": "HelloWorld!",
			},
			JSONMatcher: map[string]objmatcher.Matcher{
				"type": objmatcher.NewEqualsMatcher("diagnostic.1"),
				"time": objmatcher.NewRegExpMatcher(".+"),
				"diagnostic": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"type": objmatcher.NewEqualsMatcher("generic"),
					"generic": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"diagnosticType": objmatcher.NewEqualsMatcher("DIAG_TYPE"),
						"value": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"testKey": objmatcher.NewEqualsMatcher("test_value"),
						}),
					}),
				}),
				"unsafeParams": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"Password": objmatcher.NewEqualsMatcher("HelloWorld!"),
				}),
			},
		},
		{
			Name: "thread dump diagnostic log entry",
			Diagnostic: logging.NewDiagnosticFromThreadDump(logging.ThreadDumpV1{
				Threads: []logging.ThreadInfoV1{
					{
						Id:   safeLongVal(13),
						Name: stringVal("testName"),
						StackTrace: []logging.StackFrameV1{
							{
								Address:   stringVal("address_val"),
								Procedure: stringVal("procedure_val"),
								File:      stringVal("file_val"),
								Line:      intVal(99),
								Params: map[string]interface{}{
									"stackFrameParam": 33,
								},
							},
						},
						Params: map[string]interface{}{
							"threadParam": 77,
						},
					},
				},
			}),
			UnsafeParams: map[string]interface{}{
				"Password": "HelloWorld!",
			},
			JSONMatcher: map[string]objmatcher.Matcher{
				"type": objmatcher.NewEqualsMatcher("diagnostic.1"),
				"time": objmatcher.NewRegExpMatcher(".+"),
				"diagnostic": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"type": objmatcher.NewEqualsMatcher("threadDump"),
					"threadDump": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"threads": objmatcher.SliceMatcher([]objmatcher.Matcher{
							objmatcher.MapMatcher(map[string]objmatcher.Matcher{
								"id":   objmatcher.NewEqualsMatcher(json.Number("13")),
								"name": objmatcher.NewEqualsMatcher("testName"),
								"stackTrace": objmatcher.SliceMatcher([]objmatcher.Matcher{
									objmatcher.MapMatcher(map[string]objmatcher.Matcher{
										"address":   objmatcher.NewEqualsMatcher("address_val"),
										"procedure": objmatcher.NewEqualsMatcher("procedure_val"),
										"file":      objmatcher.NewEqualsMatcher("file_val"),
										"line":      objmatcher.NewEqualsMatcher(json.Number("99")),
										"params": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
											"stackFrameParam": objmatcher.NewEqualsMatcher(json.Number("33")),
										}),
									}),
								}),
								"params": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
									"threadParam": objmatcher.NewEqualsMatcher(json.Number("77")),
								}),
							}),
						}),
					}),
				}),
				"unsafeParams": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"Password": objmatcher.NewEqualsMatcher("HelloWorld!"),
				}),
			},
		},
	}
}

func safeLongVal(in int64) *conjuretype.SafeLong {
	val, err := conjuretype.NewSafeLong(in)
	if err != nil {
		panic(err)
	}
	return &val
}

func intVal(in int) *int {
	return &in
}

func stringVal(in string) *string {
	return &in
}

func JSONTestSuite(t *testing.T, loggerProvider func(w io.Writer) diag1log.Logger) {
	jsonOutputTests(t, loggerProvider)
}

func jsonOutputTests(t *testing.T, loggerProvider func(w io.Writer) diag1log.Logger) {
	for i, tc := range TestCases() {
		t.Run(tc.Name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := loggerProvider(buf)

			logger.Diagnostic(
				tc.Diagnostic,
				diag1log.UnsafeParams(tc.UnsafeParams),
			)

			gotEventLog := map[string]interface{}{}
			logEntry := buf.Bytes()
			err := safejson.Unmarshal(logEntry, &gotEventLog)

			require.NoError(t, err, "Case %d: %s\nEvent log line is not a valid map: %v", i, tc.Name, string(logEntry))

			assert.NoError(t, tc.JSONMatcher.Matches(gotEventLog), "Case %d: %s", i, tc.Name)
		})
	}
}
