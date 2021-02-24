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

package wrapped1logtests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/palantir/pkg/objmatcher"
	"github.com/palantir/pkg/safejson"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/svclog/svc1log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Svc1LogJSONTestSuite(t *testing.T, entityName string, entityVersion string, loggerProvider func(w io.Writer, level wlog.LogLevel, origin string) svc1log.Logger) {
	jsonOutputTests(t, entityName, entityVersion, loggerProvider)
	jsonParamsOnlyMarshaledIfLoggedTest(t, loggerProvider)
	paramIsntOverwrittenByParams(t, entityName, entityVersion, loggerProvider)
	extraParamsDoNotAppearTest(t, entityName, entityVersion, loggerProvider)
	jsonLoggerUpdateTest(t, entityName, entityVersion, loggerProvider)
}

type testStruct struct {
	NumVal            int `json:"num-val"`
	ExportedStringVal string
	privateStrVal     string
}

type testParamStorerObject struct {
	safeParams   map[string]interface{}
	unsafeParams map[string]interface{}
}

func (t testParamStorerObject) SafeParams() map[string]interface{} {
	return t.safeParams
}

func (t testParamStorerObject) UnsafeParams() map[string]interface{} {
	return t.unsafeParams
}

type testError struct {
	message      string
	stacktrace   string
	safeParams   map[string]interface{}
	unsafeParams map[string]interface{}
}

func (t testError) Error() string {
	return t.message
}

func (t testError) Format(state fmt.State, c rune) {
	if state.Flag('+') && c == 'v' {
		_, _ = fmt.Fprint(state, t.stacktrace)
	}
}

func (t testError) SafeParams() map[string]interface{} {
	return t.safeParams
}

func (t testError) UnsafeParams() map[string]interface{} {
	return t.unsafeParams
}

type TestCase struct {
	Name        string
	Message     string
	Origin      string
	LogParams   []svc1log.Param
	JSONMatcher objmatcher.MapMatcher
}

func TestCases(entityName string, entityVersion string) []TestCase {
	return []TestCase{
		{
			Name:    "basic service log entry",
			Message: "this is a test",
			LogParams: []svc1log.Param{
				svc1log.UID("user-1"),
				svc1log.SID("session-1"),
				svc1log.TraceID("X-Y-Z"),
				svc1log.SafeParams(map[string]interface{}{
					"key": "value",
					"int": 10,
				}),
				svc1log.UnsafeParams(map[string]interface{}{
					"Password": "HelloWorld!",
				}),
				svc1log.Tags(map[string]string{
					"key1": "value1",
					"key2": "value2",
				}),
			},
			JSONMatcher: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"type":          objmatcher.NewEqualsMatcher("wrapped.1"),
				"entityName":    objmatcher.NewEqualsMatcher(entityName),
				"entityVersion": objmatcher.NewEqualsMatcher(entityVersion),
				"payload": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"type": objmatcher.NewEqualsMatcher("serviceLogV1"),
					"serviceLogV1": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"level":   objmatcher.NewEqualsMatcher("INFO"),
						"type":    objmatcher.NewEqualsMatcher("service.1"),
						"time":    objmatcher.NewRegExpMatcher(".+"),
						"message": objmatcher.NewEqualsMatcher("this is a test"),
						"params": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"key": objmatcher.NewEqualsMatcher("value"),
							"int": objmatcher.NewEqualsMatcher(json.Number("10")),
						}),
						"uid":     objmatcher.NewEqualsMatcher("user-1"),
						"sid":     objmatcher.NewEqualsMatcher("session-1"),
						"traceId": objmatcher.NewEqualsMatcher("X-Y-Z"),
						"unsafeParams": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"Password": objmatcher.NewEqualsMatcher("HelloWorld!"),
						}),
						"tags": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"key1": objmatcher.NewEqualsMatcher("value1"),
							"key2": objmatcher.NewEqualsMatcher("value2"),
						}),
					}),
				}),
			}),
		},
		{
			Name:    "service log entry with non-primitive objects in params map",
			Message: "this is a test",
			LogParams: []svc1log.Param{
				svc1log.UID("user-1"),
				svc1log.SID("session-1"),
				svc1log.TraceID("X-Y-Z"),
				svc1log.SafeParams(map[string]interface{}{
					"structKey": testStruct{
						NumVal:            13,
						ExportedStringVal: "exportedFoo",
						privateStrVal:     "privateFoo",
					},
					"mapKey": map[string]interface{}{
						"mapKey1": "map-val-1",
					},
					"sliceKey":  []string{"one", "two", "three"},
					"stringKey": "stringVal",
				}),
				svc1log.UnsafeParams(map[string]interface{}{
					"structKey": testStruct{
						NumVal:            13,
						ExportedStringVal: "exportedFoo",
						privateStrVal:     "privateFoo",
					},
					"mapKey": map[string]interface{}{
						"mapKey1": "map-val-1",
					},
					"sliceKey":  []string{"one", "two", "three"},
					"stringKey": "stringVal",
				}),
			},
			JSONMatcher: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"type":          objmatcher.NewEqualsMatcher("wrapped.1"),
				"entityName":    objmatcher.NewEqualsMatcher(entityName),
				"entityVersion": objmatcher.NewEqualsMatcher(entityVersion),
				"payload": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"type": objmatcher.NewEqualsMatcher("serviceLogV1"),
					"serviceLogV1": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"level":   objmatcher.NewEqualsMatcher("INFO"),
						"time":    objmatcher.NewRegExpMatcher(".+"),
						"type":    objmatcher.NewEqualsMatcher("service.1"),
						"message": objmatcher.NewEqualsMatcher("this is a test"),
						"params": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"structKey": objmatcher.NewEqualsMatcher(map[string]interface{}{
								"num-val":           json.Number("13"),
								"ExportedStringVal": "exportedFoo",
								// note: "privateStrVal" not expected to be included because it is not an exported field
							}),
							"mapKey": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
								"mapKey1": objmatcher.NewEqualsMatcher("map-val-1"),
							}),
							"sliceKey":  objmatcher.NewEqualsMatcher([]interface{}{"one", "two", "three"}),
							"stringKey": objmatcher.NewEqualsMatcher("stringVal"),
						}),
						"uid":     objmatcher.NewEqualsMatcher("user-1"),
						"sid":     objmatcher.NewEqualsMatcher("session-1"),
						"traceId": objmatcher.NewEqualsMatcher("X-Y-Z"),
						"unsafeParams": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"structKey": objmatcher.NewEqualsMatcher(map[string]interface{}{
								"num-val":           json.Number("13"),
								"ExportedStringVal": "exportedFoo",
								// note: "privateStrVal" not expected to be included because it is not an exported field
							}),
							"mapKey": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
								"mapKey1": objmatcher.NewEqualsMatcher("map-val-1"),
							}),
							"sliceKey":  objmatcher.NewEqualsMatcher([]interface{}{"one", "two", "three"}),
							"stringKey": objmatcher.NewEqualsMatcher("stringVal"),
						}),
					}),
				}),
			}),
		},
		{
			Name:    "service log entry with origin set on base logger",
			Message: "this is a test",
			Origin:  "github.com/palantir/witchcraft-go-logging",
			LogParams: []svc1log.Param{
				svc1log.SafeParams(map[string]interface{}{
					"key": "value",
					"int": 10,
				}),
				svc1log.UnsafeParams(map[string]interface{}{
					"Password": "HelloWorld!",
				}),
			},
			JSONMatcher: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"type":          objmatcher.NewEqualsMatcher("wrapped.1"),
				"entityName":    objmatcher.NewEqualsMatcher(entityName),
				"entityVersion": objmatcher.NewEqualsMatcher(entityVersion),
				"payload": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"type": objmatcher.NewEqualsMatcher("serviceLogV1"),
					"serviceLogV1": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"level":   objmatcher.NewEqualsMatcher("INFO"),
						"time":    objmatcher.NewRegExpMatcher(".+"),
						"origin":  objmatcher.NewEqualsMatcher("github.com/palantir/witchcraft-go-logging"),
						"type":    objmatcher.NewEqualsMatcher("service.1"),
						"message": objmatcher.NewEqualsMatcher("this is a test"),
						"params": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"key": objmatcher.NewEqualsMatcher("value"),
							"int": objmatcher.NewEqualsMatcher(json.Number("10")),
						}),
						"unsafeParams": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"Password": objmatcher.NewEqualsMatcher("HelloWorld!"),
						}),
					}),
				}),
			}),
		},
		{
			Name:      "parameter that is set manually overrides base value",
			Message:   "this is a test",
			Origin:    "github.com/palantir/witchcraft-go-logging",
			LogParams: []svc1log.Param{svc1log.Origin("custom-origin")},
			JSONMatcher: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"type":          objmatcher.NewEqualsMatcher("wrapped.1"),
				"entityName":    objmatcher.NewEqualsMatcher(entityName),
				"entityVersion": objmatcher.NewEqualsMatcher(entityVersion),
				"payload": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"type": objmatcher.NewEqualsMatcher("serviceLogV1"),
					"serviceLogV1": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"level":   objmatcher.NewEqualsMatcher("INFO"),
						"time":    objmatcher.NewRegExpMatcher(".+"),
						"origin":  objmatcher.NewEqualsMatcher("custom-origin"),
						"type":    objmatcher.NewEqualsMatcher("service.1"),
						"message": objmatcher.NewEqualsMatcher("this is a test"),
					}),
				}),
			}),
		},
		{
			Name:    "stacktrace includes error parameters",
			Message: "something happened",
			Origin:  "github.com/palantir/witchcraft-go-logging",
			LogParams: []svc1log.Param{
				svc1log.Stacktrace(
					testError{
						message: "some error message",
						stacktrace: `Failed to open file
something/something:123`,
						safeParams: map[string]interface{}{
							"safeKey": "safeVal",
						},
						unsafeParams: map[string]interface{}{
							"unsafeKey": "unsafeVal",
						},
					},
				),
			},
			JSONMatcher: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"type":          objmatcher.NewEqualsMatcher("wrapped.1"),
				"entityName":    objmatcher.NewEqualsMatcher(entityName),
				"entityVersion": objmatcher.NewEqualsMatcher(entityVersion),
				"payload": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"type": objmatcher.NewEqualsMatcher("serviceLogV1"),
					"serviceLogV1": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"level":   objmatcher.NewEqualsMatcher("INFO"),
						"time":    objmatcher.NewRegExpMatcher(".+"),
						"origin":  objmatcher.NewEqualsMatcher("github.com/palantir/witchcraft-go-logging"),
						"type":    objmatcher.NewEqualsMatcher("service.1"),
						"message": objmatcher.NewEqualsMatcher("something happened"),
						"params": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"safeKey": objmatcher.NewEqualsMatcher("safeVal"),
						}),
						"unsafeParams": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"unsafeKey": objmatcher.NewEqualsMatcher("unsafeVal"),
						}),
						"stacktrace": objmatcher.NewRegExpMatcher("(?s)Failed to open file.+"),
					}),
				}),
			}),
		},
		{
			Name:    "parameters included from ParamStorer parameter",
			Message: "something happened",
			Origin:  "github.com/palantir/witchcraft-go-logging",
			LogParams: []svc1log.Param{
				svc1log.Params(testParamStorerObject{
					safeParams: map[string]interface{}{
						"safeObjectParamKey": "safeObjectParamValue",
					},
					unsafeParams: map[string]interface{}{
						"unsafeObjectParamKey": "unsafeObjectParamValue",
					},
				}),
			},
			JSONMatcher: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"type":          objmatcher.NewEqualsMatcher("wrapped.1"),
				"entityName":    objmatcher.NewEqualsMatcher(entityName),
				"entityVersion": objmatcher.NewEqualsMatcher(entityVersion),
				"payload": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"type": objmatcher.NewEqualsMatcher("serviceLogV1"),
					"serviceLogV1": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"level":   objmatcher.NewEqualsMatcher("INFO"),
						"time":    objmatcher.NewRegExpMatcher(".+"),
						"type":    objmatcher.NewEqualsMatcher("service.1"),
						"message": objmatcher.NewEqualsMatcher("something happened"),
						"origin":  objmatcher.NewEqualsMatcher("github.com/palantir/witchcraft-go-logging"),
						"params": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"safeObjectParamKey": objmatcher.NewEqualsMatcher("safeObjectParamValue"),
						}),
						"unsafeParams": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"unsafeObjectParamKey": objmatcher.NewEqualsMatcher("unsafeObjectParamValue"),
						}),
					}),
				}),
			}),
		},
		{
			Name:    "param isn't overwritten by params",
			Message: "msg",
			LogParams: []svc1log.Param{
				svc1log.SafeParam("param", "value"),
				svc1log.SafeParams(map[string]interface{}{"params": "values"}),
			},
			JSONMatcher: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"type":          objmatcher.NewEqualsMatcher("wrapped.1"),
				"entityName":    objmatcher.NewEqualsMatcher(entityName),
				"entityVersion": objmatcher.NewEqualsMatcher(entityVersion),
				"payload": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"type": objmatcher.NewEqualsMatcher("serviceLogV1"),
					"serviceLogV1": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"level":   objmatcher.NewEqualsMatcher("INFO"),
						"message": objmatcher.NewEqualsMatcher("msg"),
						"time":    objmatcher.NewRegExpMatcher(".+"),
						"type":    objmatcher.NewEqualsMatcher("service.1"),
						"params": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"param":  objmatcher.NewEqualsMatcher("value"),
							"params": objmatcher.NewEqualsMatcher("values"),
						}),
					}),
				}),
			}),
		},
		{
			Name:    "duplicate origin",
			Message: "this is a test",
			LogParams: []svc1log.Param{
				svc1log.Origin("origin.0"),
				svc1log.Origin("origin.1"),
				svc1log.UID("user-1"),
				svc1log.SID("session-1"),
				svc1log.TraceID("X-Y-Z"),
				svc1log.SafeParams(map[string]interface{}{
					"key": "value",
					"int": 10,
				}),
				svc1log.UnsafeParams(map[string]interface{}{
					"Password": "HelloWorld!",
				}),
				svc1log.Tags(map[string]string{
					"key1": "value1",
					"key2": "value2",
				}),
			},
			JSONMatcher: objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"type":          objmatcher.NewEqualsMatcher("wrapped.1"),
				"entityName":    objmatcher.NewEqualsMatcher(entityName),
				"entityVersion": objmatcher.NewEqualsMatcher(entityVersion),
				"payload": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"type": objmatcher.NewEqualsMatcher("serviceLogV1"),
					"serviceLogV1": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"origin":  objmatcher.NewEqualsMatcher("origin.1"),
						"level":   objmatcher.NewEqualsMatcher("INFO"),
						"time":    objmatcher.NewRegExpMatcher(".+"),
						"type":    objmatcher.NewEqualsMatcher("service.1"),
						"message": objmatcher.NewEqualsMatcher("this is a test"),
						"params": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"key": objmatcher.NewEqualsMatcher("value"),
							"int": objmatcher.NewEqualsMatcher(json.Number("10")),
						}),
						"uid":     objmatcher.NewEqualsMatcher("user-1"),
						"sid":     objmatcher.NewEqualsMatcher("session-1"),
						"traceId": objmatcher.NewEqualsMatcher("X-Y-Z"),
						"unsafeParams": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"Password": objmatcher.NewEqualsMatcher("HelloWorld!"),
						}),
						"tags": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
							"key1": objmatcher.NewEqualsMatcher("value1"),
							"key2": objmatcher.NewEqualsMatcher("value2"),
						}),
					}),
				}),
			}),
		},
	}
}

func jsonOutputTests(t *testing.T, entityName string, entityVersion string, loggerProvider func(w io.Writer, level wlog.LogLevel, origin string) svc1log.Logger) {
	for i, tc := range TestCases(entityName, entityVersion) {
		t.Run(tc.Name, func(t *testing.T) {
			buf := bytes.Buffer{}
			logger := loggerProvider(&buf, wlog.DebugLevel, tc.Origin)

			logger.Info(tc.Message, tc.LogParams...)

			gotServiceLog := map[string]interface{}{}
			logEntry := buf.Bytes()
			err := safejson.Unmarshal(logEntry, &gotServiceLog)
			require.NoError(t, err, "Case %d: %s\nService log line is not a valid map: %v", i, tc.Name, string(logEntry))
			logEntryRewrite, err := safejson.Marshal(gotServiceLog)
			require.NoError(t, err, "Not able to marshal log line")
			assert.Equal(t, len(strings.TrimRight(string(logEntry), "\n")), len(string(logEntryRewrite)), "Differing length on remarshal, possibly due to duplicate keys in the original payload.")

			assert.NoError(t, tc.JSONMatcher.Matches(gotServiceLog), "Case %d: %s", i, tc.Name)
		})
	}
}
func jsonParamsOnlyMarshaledIfLoggedTest(t *testing.T, loggerProvider func(w io.Writer, level wlog.LogLevel, origin string) svc1log.Logger) {
	t.Run("Params only marshaled if logged", func(t *testing.T) {
		logger := loggerProvider(&bytes.Buffer{}, wlog.InfoLevel, "")
		// demonstrates that writing to a log at a level that is lower than the logger's level will not marshal the
		// parameters (if marshal occurred, this would panic).
		assert.NotPanics(t, func() {
			logger.Debug("Test Message", svc1log.SafeParam("testType", jsonMarshalPanicType{}))
		})
	})
}

// Verifies that if different parameters are specified using SafeParam and SafeParams params, all of the values are
// present in the final output (that is, these parameters should be additive).
func paramIsntOverwrittenByParams(t *testing.T, entityName string, entityVersion string, loggerProvider func(w io.Writer, level wlog.LogLevel, origin string) svc1log.Logger) {
	t.Run("SafeParam and SafeParams params are additive", func(t *testing.T) {
		var buf bytes.Buffer
		logger := loggerProvider(&buf, wlog.InfoLevel, "")

		logger.Info("msg", svc1log.SafeParam("param", "value"), svc1log.SafeParams(map[string]interface{}{"params": "values"}))

		gotServiceLog := map[string]interface{}{}
		logEntry := buf.Bytes()
		err := safejson.Unmarshal(logEntry, &gotServiceLog)
		require.NoError(t, err, "Service log line is not a valid map: %v", string(logEntry))

		assert.NoError(t, objmatcher.MapMatcher(map[string]objmatcher.Matcher{
			"type":          objmatcher.NewEqualsMatcher("wrapped.1"),
			"entityName":    objmatcher.NewEqualsMatcher(entityName),
			"entityVersion": objmatcher.NewEqualsMatcher(entityVersion),
			"payload": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"type": objmatcher.NewEqualsMatcher("serviceLogV1"),
				"serviceLogV1": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"level":   objmatcher.NewEqualsMatcher("INFO"),
					"message": objmatcher.NewEqualsMatcher("msg"),
					"time":    objmatcher.NewRegExpMatcher(".+"),
					"type":    objmatcher.NewEqualsMatcher("service.1"),
					"params": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"param":  objmatcher.NewEqualsMatcher("value"),
						"params": objmatcher.NewEqualsMatcher("values"),
					}),
				}),
			}),
		}).Matches(gotServiceLog))
	})
}

// Verifies that parameters remain separate between different logger calls (ensures there is not a bug where parameters
// are modified by making a logger call).
func extraParamsDoNotAppearTest(t *testing.T, entityName string, entityVersion string, loggerProvider func(w io.Writer, level wlog.LogLevel, origin string) svc1log.Logger) {
	t.Run("SafeParam and SafeParams params stay separate across logger calls", func(t *testing.T) {
		var buf bytes.Buffer
		logger := loggerProvider(&buf, wlog.DebugLevel, "")

		reusedParams := svc1log.SafeParams(map[string]interface{}{"params": "values"})
		logger.Info("msg", reusedParams, svc1log.SafeParam("param", "value"))
		gotServiceLog := map[string]interface{}{}
		logEntry := buf.Bytes()
		err := safejson.Unmarshal(logEntry, &gotServiceLog)
		require.NoError(t, err, "Service log line is not a valid map: %v", string(logEntry))
		assert.NoError(t, objmatcher.MapMatcher(map[string]objmatcher.Matcher{
			"type":          objmatcher.NewEqualsMatcher("wrapped.1"),
			"entityName":    objmatcher.NewEqualsMatcher(entityName),
			"entityVersion": objmatcher.NewEqualsMatcher(entityVersion),
			"payload": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"type": objmatcher.NewEqualsMatcher("serviceLogV1"),
				"serviceLogV1": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"level":   objmatcher.NewEqualsMatcher("INFO"),
					"message": objmatcher.NewEqualsMatcher("msg"),
					"time":    objmatcher.NewRegExpMatcher(".+"),
					"type":    objmatcher.NewEqualsMatcher("service.1"),
					"params": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"param":  objmatcher.NewEqualsMatcher("value"),
						"params": objmatcher.NewEqualsMatcher("values"),
					}),
				}),
			}),
		}).Matches(gotServiceLog))

		buf.Reset()
		logger.Info("msg", reusedParams)

		gotServiceLog = map[string]interface{}{}
		logEntry = buf.Bytes()
		err = safejson.Unmarshal(logEntry, &gotServiceLog)
		require.NoError(t, err, "Service log line is not a valid map: %v", string(logEntry))
		assert.NoError(t, objmatcher.MapMatcher(map[string]objmatcher.Matcher{
			"type":          objmatcher.NewEqualsMatcher("wrapped.1"),
			"entityName":    objmatcher.NewEqualsMatcher(entityName),
			"entityVersion": objmatcher.NewEqualsMatcher(entityVersion),
			"payload": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
				"type": objmatcher.NewEqualsMatcher("serviceLogV1"),
				"serviceLogV1": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"level":   objmatcher.NewEqualsMatcher("INFO"),
					"message": objmatcher.NewEqualsMatcher("msg"),
					"time":    objmatcher.NewRegExpMatcher(".+"),
					"type":    objmatcher.NewEqualsMatcher("service.1"),
					"params": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
						"params": objmatcher.NewEqualsMatcher("values"),
					}),
				}),
			}),
		}).Matches(gotServiceLog))
	})
}

func jsonLoggerUpdateTest(t *testing.T, entityName string, entityVersion string, loggerProvider func(w io.Writer, level wlog.LogLevel, origin string) svc1log.Logger) {
	t.Run("update JSON logger", func(t *testing.T) {
		currCase := TestCases(entityName, entityVersion)[0]

		buf := bytes.Buffer{}
		logger := loggerProvider(&buf, wlog.ErrorLevel, currCase.Origin)

		// log at info level
		logger.Info(currCase.Message, currCase.LogParams...)

		// output should be empty
		assert.Equal(t, "", buf.String())

		// update configuration to log at info level
		logger.SetLevel(wlog.InfoLevel)

		// log at info level
		logger.Info(currCase.Message, currCase.LogParams...)

		// output should exist and match
		gotServiceLog := map[string]interface{}{}
		logEntry := buf.Bytes()
		err := safejson.Unmarshal(logEntry, &gotServiceLog)
		require.NoError(t, err, "Service log line is not a valid map: %v", string(logEntry))

		assert.NoError(t, currCase.JSONMatcher.Matches(gotServiceLog), "No match")
	})
}

// panics when marshaled as JSON
type jsonMarshalPanicType struct{}

func (t jsonMarshalPanicType) MarshalJSON() ([]byte, error) {
	panic("jsonMarshalPanicType panics on MarshalJSON")
}
