// Copyright (c) 2021 Palantir Technologies. All rights reserved.
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

package wrapped1logtests

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/palantir/pkg/objmatcher"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/pkg/safelong"
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog/diaglog/diag1log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Diag1TestCase struct {
	Name         string
	Diagnostic   logging.Diagnostic
	UnsafeParams map[string]any
	JSONMatcher  objmatcher.MapMatcher
}

func Diag1TestCases(entityName, entityVersion string) []Diag1TestCase {
	return []Diag1TestCase{
		{
			Name: "generic diagnostic log entry",
			Diagnostic: logging.NewDiagnosticFromGeneric(logging.GenericDiagnostic{
				DiagnosticType: "DIAG_TYPE",
				Value: map[string]string{
					"testKey": "test_value",
				},
			}),
			UnsafeParams: map[string]any{
				"Password": "HelloWorld!",
			},
			JSONMatcher: map[string]objmatcher.Matcher{
				"type":          objmatcher.NewEqualsMatcher("wrapped.1"),
				"entityName":    objmatcher.NewEqualsMatcher(entityName),
				"entityVersion": objmatcher.NewEqualsMatcher(entityVersion),
				"payload": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"type": objmatcher.NewEqualsMatcher("diagnosticLogV1"),
					"diagnosticLogV1": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
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
					}),
				}),
			},
		},
		{
			Name: "thread dump diagnostic log entry",
			Diagnostic: logging.NewDiagnosticFromThreadDump(logging.ThreadDumpV1{
				Threads: []logging.ThreadInfoV1{
					{
						Id:   safeLongVal(13),
						Name: new("testName"),
						StackTrace: []logging.StackFrameV1{
							{
								Address:   new("address_val"),
								Procedure: new("procedure_val"),
								File:      new("file_val"),
								Line:      new(99),
								Params: map[string]any{
									"stackFrameParam": 33,
								},
							},
						},
						Params: map[string]any{
							"threadParam": 77,
						},
					},
				},
			}),
			UnsafeParams: map[string]any{
				"Password": "HelloWorld!",
			},
			JSONMatcher: map[string]objmatcher.Matcher{
				"type":          objmatcher.NewEqualsMatcher("wrapped.1"),
				"entityName":    objmatcher.NewEqualsMatcher(entityName),
				"entityVersion": objmatcher.NewEqualsMatcher(entityVersion),
				"payload": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"type": objmatcher.NewEqualsMatcher("diagnosticLogV1"),
					"diagnosticLogV1": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
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
					}),
				}),
			},
		},
	}
}

func safeLongVal(in int64) *safelong.SafeLong {
	val, err := safelong.NewSafeLong(in)
	if err != nil {
		panic(err)
	}
	return &val
}

//go:fix inline
func intVal(in int) *int {
	return new(in)
}

//go:fix inline
func stringVal(in string) *string {
	return new(in)
}

func Diag1LogJSONTestSuite(t *testing.T, entityName, entityVersion string, loggerProvider func(w io.Writer) diag1log.Logger) {
	diag1LogJSONOutputTests(t, entityName, entityVersion, loggerProvider)
}

func diag1LogJSONOutputTests(t *testing.T, entityName, entityVersion string, loggerProvider func(w io.Writer) diag1log.Logger) {
	for i, tc := range Diag1TestCases(entityName, entityVersion) {
		t.Run(tc.Name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := loggerProvider(buf)

			logger.Diagnostic(
				tc.Diagnostic,
				diag1log.UnsafeParams(tc.UnsafeParams),
			)

			gotEventLog := map[string]any{}
			logEntry := buf.Bytes()
			err := safejson.Unmarshal(logEntry, &gotEventLog)

			require.NoError(t, err, "Case %d: %s\nEvent log line is not a valid map: %v", i, tc.Name, string(logEntry))

			assert.NoError(t, tc.JSONMatcher.Matches(gotEventLog), "Case %d: %s", i, tc.Name)
		})
	}
}
