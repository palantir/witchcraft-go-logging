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

package req2logtests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/palantir/pkg/objmatcher"
	"github.com/palantir/witchcraft-go-logging/wlog/logreader"
	"github.com/palantir/witchcraft-go-logging/wlog/reqlog/req2log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	Name                  string
	ExtraHeaderParams     map[string]string
	ExtraQueryParams      []string
	SafeHeaderParams      []string
	SafeQueryParams       []string
	ForbiddenHeaderParams []string
	JSONMatcher           objmatcher.MapMatcher
}

func TestCases() []TestCase {
	return []TestCase{
		{
			Name: "request.2 log entry with no whitelisted content",
			JSONMatcher: map[string]objmatcher.Matcher{
				"type":         objmatcher.NewEqualsMatcher("request.2"),
				"time":         objmatcher.NewRegExpMatcher(".+"),
				"method":       objmatcher.NewEqualsMatcher("GET"),
				"protocol":     objmatcher.NewEqualsMatcher("HTTP/1.1"),
				"path":         objmatcher.NewEqualsMatcher("/some/path/here"),
				"status":       objmatcher.NewEqualsMatcher(json.Number("200")),
				"requestSize":  objmatcher.NewEqualsMatcher(json.Number("0")),
				"responseSize": objmatcher.NewEqualsMatcher(json.Number("100")),
				"duration":     objmatcher.NewAnyMatcher(),
				"uid":          objmatcher.NewEqualsMatcher("be9f645d-52e0-49e9-ba31-db32927615db"),
				"sid":          objmatcher.NewEqualsMatcher("ad4d4ae6-65af-4e2a-91a3-cf401acb1d4c"),
				"tokenId":      objmatcher.NewEqualsMatcher("9277f9af-8d99-408a-94ef-f51e82be2ff8"),
				"orgId":        objmatcher.NewEqualsMatcher("0998e573-31d7-4999-8bf9-0bc5f4592db9"),
				"unsafeParams": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"Fooheaderparamname": objmatcher.NewEqualsMatcher("fooHeaderParamVal"),
					"fooQueryVarName":    objmatcher.NewEqualsMatcher("fooQueryVarVal"),
					"barQueryVarName":    objmatcher.NewEqualsMatcher("barQueryVarVal"),
				}),
			},
		},
		{
			Name:                  "request.2 log entry with forbidden header",
			ForbiddenHeaderParams: []string{"FooHeaderParamName"},
			JSONMatcher: map[string]objmatcher.Matcher{
				"type":         objmatcher.NewEqualsMatcher("request.2"),
				"time":         objmatcher.NewRegExpMatcher(".+"),
				"method":       objmatcher.NewEqualsMatcher("GET"),
				"protocol":     objmatcher.NewEqualsMatcher("HTTP/1.1"),
				"path":         objmatcher.NewEqualsMatcher("/some/path/here"),
				"status":       objmatcher.NewEqualsMatcher(json.Number("200")),
				"requestSize":  objmatcher.NewEqualsMatcher(json.Number("0")),
				"responseSize": objmatcher.NewEqualsMatcher(json.Number("100")),
				"duration":     objmatcher.NewAnyMatcher(),
				"uid":          objmatcher.NewEqualsMatcher("be9f645d-52e0-49e9-ba31-db32927615db"),
				"sid":          objmatcher.NewEqualsMatcher("ad4d4ae6-65af-4e2a-91a3-cf401acb1d4c"),
				"tokenId":      objmatcher.NewEqualsMatcher("9277f9af-8d99-408a-94ef-f51e82be2ff8"),
				"orgId":        objmatcher.NewEqualsMatcher("0998e573-31d7-4999-8bf9-0bc5f4592db9"),
				"unsafeParams": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"fooQueryVarName": objmatcher.NewEqualsMatcher("fooQueryVarVal"),
					"barQueryVarName": objmatcher.NewEqualsMatcher("barQueryVarVal"),
				}),
			},
		},
		{
			Name:             "request.2 log entry with whitelisted header",
			SafeHeaderParams: []string{"Fooheaderparamname"},
			JSONMatcher: map[string]objmatcher.Matcher{
				"type":     objmatcher.NewEqualsMatcher("request.2"),
				"time":     objmatcher.NewRegExpMatcher(".+"),
				"method":   objmatcher.NewEqualsMatcher("GET"),
				"protocol": objmatcher.NewEqualsMatcher("HTTP/1.1"),
				"path":     objmatcher.NewEqualsMatcher("/some/path/here"),
				"params": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"Fooheaderparamname": objmatcher.NewEqualsMatcher("fooHeaderParamVal"),
				}),
				"status":       objmatcher.NewEqualsMatcher(json.Number("200")),
				"requestSize":  objmatcher.NewEqualsMatcher(json.Number("0")),
				"responseSize": objmatcher.NewEqualsMatcher(json.Number("100")),
				"duration":     objmatcher.NewAnyMatcher(),
				"uid":          objmatcher.NewEqualsMatcher("be9f645d-52e0-49e9-ba31-db32927615db"),
				"sid":          objmatcher.NewEqualsMatcher("ad4d4ae6-65af-4e2a-91a3-cf401acb1d4c"),
				"tokenId":      objmatcher.NewEqualsMatcher("9277f9af-8d99-408a-94ef-f51e82be2ff8"),
				"orgId":        objmatcher.NewEqualsMatcher("0998e573-31d7-4999-8bf9-0bc5f4592db9"),
				"unsafeParams": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"fooQueryVarName": objmatcher.NewEqualsMatcher("fooQueryVarVal"),
					"barQueryVarName": objmatcher.NewEqualsMatcher("barQueryVarVal"),
				}),
			},
		},
		{
			Name:            "request.2 log entry with one whitelisted query parameter",
			SafeQueryParams: []string{"fooQueryVarName"},
			JSONMatcher: map[string]objmatcher.Matcher{
				"type":     objmatcher.NewEqualsMatcher("request.2"),
				"time":     objmatcher.NewRegExpMatcher(".+"),
				"method":   objmatcher.NewEqualsMatcher("GET"),
				"protocol": objmatcher.NewEqualsMatcher("HTTP/1.1"),
				"path":     objmatcher.NewEqualsMatcher("/some/path/here"),
				"params": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"fooQueryVarName": objmatcher.NewEqualsMatcher("fooQueryVarVal"),
				}),
				"status":       objmatcher.NewEqualsMatcher(json.Number("200")),
				"requestSize":  objmatcher.NewEqualsMatcher(json.Number("0")),
				"responseSize": objmatcher.NewEqualsMatcher(json.Number("100")),
				"duration":     objmatcher.NewAnyMatcher(),
				"uid":          objmatcher.NewEqualsMatcher("be9f645d-52e0-49e9-ba31-db32927615db"),
				"sid":          objmatcher.NewEqualsMatcher("ad4d4ae6-65af-4e2a-91a3-cf401acb1d4c"),
				"tokenId":      objmatcher.NewEqualsMatcher("9277f9af-8d99-408a-94ef-f51e82be2ff8"),
				"orgId":        objmatcher.NewEqualsMatcher("0998e573-31d7-4999-8bf9-0bc5f4592db9"),
				"unsafeParams": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"Fooheaderparamname": objmatcher.NewEqualsMatcher("fooHeaderParamVal"),
					"barQueryVarName":    objmatcher.NewEqualsMatcher("barQueryVarVal"),
				}),
			},
		},
		{
			Name: "request.2 log entry with multi-value parameters",
			ExtraQueryParams: []string{
				"fooQueryVarName=extra-val-foo",
				"fooqueryvarname=case-sensitive-so-val-alone",
				"barQueryVarName=extra-val-bar",
			},
			SafeQueryParams: []string{"fooQueryVarName"},
			JSONMatcher: map[string]objmatcher.Matcher{
				"type":     objmatcher.NewEqualsMatcher("request.2"),
				"time":     objmatcher.NewRegExpMatcher(".+"),
				"method":   objmatcher.NewEqualsMatcher("GET"),
				"protocol": objmatcher.NewEqualsMatcher("HTTP/1.1"),
				"path":     objmatcher.NewEqualsMatcher("/some/path/here"),
				"params": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"fooQueryVarName": objmatcher.NewEqualsMatcher([]interface{}{"fooQueryVarVal", "extra-val-foo"}),
					"fooqueryvarname": objmatcher.NewEqualsMatcher("case-sensitive-so-val-alone"),
				}),
				"status":       objmatcher.NewEqualsMatcher(json.Number("200")),
				"requestSize":  objmatcher.NewEqualsMatcher(json.Number("0")),
				"responseSize": objmatcher.NewEqualsMatcher(json.Number("100")),
				"duration":     objmatcher.NewAnyMatcher(),
				"uid":          objmatcher.NewEqualsMatcher("be9f645d-52e0-49e9-ba31-db32927615db"),
				"sid":          objmatcher.NewEqualsMatcher("ad4d4ae6-65af-4e2a-91a3-cf401acb1d4c"),
				"tokenId":      objmatcher.NewEqualsMatcher("9277f9af-8d99-408a-94ef-f51e82be2ff8"),
				"orgId":        objmatcher.NewEqualsMatcher("0998e573-31d7-4999-8bf9-0bc5f4592db9"),
				"unsafeParams": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"Fooheaderparamname": objmatcher.NewEqualsMatcher("fooHeaderParamVal"),
					"barQueryVarName":    objmatcher.NewEqualsMatcher([]interface{}{"barQueryVarVal", "extra-val-bar"}),
				}),
			},
		},
		{
			Name:              "request.2 log entry with default safe header User-Agent",
			ExtraHeaderParams: map[string]string{"User-Agent": "userAgentVal"},
			JSONMatcher: map[string]objmatcher.Matcher{
				"type":     objmatcher.NewEqualsMatcher("request.2"),
				"time":     objmatcher.NewRegExpMatcher(".+"),
				"method":   objmatcher.NewEqualsMatcher("GET"),
				"protocol": objmatcher.NewEqualsMatcher("HTTP/1.1"),
				"path":     objmatcher.NewEqualsMatcher("/some/path/here"),
				"params": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"User-Agent": objmatcher.NewEqualsMatcher("userAgentVal"),
				}),
				"status":       objmatcher.NewEqualsMatcher(json.Number("200")),
				"requestSize":  objmatcher.NewEqualsMatcher(json.Number("0")),
				"responseSize": objmatcher.NewEqualsMatcher(json.Number("100")),
				"duration":     objmatcher.NewAnyMatcher(),
				"uid":          objmatcher.NewEqualsMatcher("be9f645d-52e0-49e9-ba31-db32927615db"),
				"sid":          objmatcher.NewEqualsMatcher("ad4d4ae6-65af-4e2a-91a3-cf401acb1d4c"),
				"tokenId":      objmatcher.NewEqualsMatcher("9277f9af-8d99-408a-94ef-f51e82be2ff8"),
				"orgId":        objmatcher.NewEqualsMatcher("0998e573-31d7-4999-8bf9-0bc5f4592db9"),
				"unsafeParams": objmatcher.MapMatcher(map[string]objmatcher.Matcher{
					"fooQueryVarName":    objmatcher.NewEqualsMatcher("fooQueryVarVal"),
					"barQueryVarName":    objmatcher.NewEqualsMatcher("barQueryVarVal"),
					"Fooheaderparamname": objmatcher.NewEqualsMatcher("fooHeaderParamVal"),
				}),
			},
		},
	}
}

func JSONTestSuite(t *testing.T, loggerProvider func(w io.Writer, params ...req2log.LoggerCreatorParam) req2log.Logger) {
	jsonOutputTests(t, loggerProvider)
}

func jsonOutputTests(t *testing.T, loggerProvider func(w io.Writer, params ...req2log.LoggerCreatorParam) req2log.Logger) {
	for i, tc := range TestCases() {
		t.Run(tc.Name, func(t *testing.T) {

			req := GenerateRequest(map[string]string{
				"Authorization":      "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ2cDlrWFZMZ1NlbTZNZHN5a25ZVjJ3PT0iLCJzaWQiOiJyVTFLNW1XdlRpcVJvODlBR3NzZFRBPT0iLCJqdGkiOiJrbmY1cjQyWlFJcVU3L1VlZ3I0ditBPT0iLCJvcmciOiJDWmpsY3pIWFNabUwrUXZGOUZrdHVRPT0ifQ.GMqKu_zrkgNR5I-jAWdR6x0G2gObVYRbqw7iJJatI4A",
				"FooHeaderParamName": "fooHeaderParamVal",
			}, tc.ExtraQueryParams, tc.ExtraHeaderParams)

			buf := &bytes.Buffer{}
			logger := loggerProvider(
				buf,
				req2log.SafeQueryParams(tc.SafeQueryParams...),
				req2log.SafeHeaderParams(tc.SafeHeaderParams...),
				req2log.ForbiddenHeaderParams(tc.ForbiddenHeaderParams...),
			)
			logger.Request(req2log.Request{
				Request:        req,
				RouteInfo:      req2log.RouteInfo{},
				ResponseStatus: http.StatusOK,
				ResponseSize:   int64(100),
				Duration:       1 * time.Second,
			})

			entries, err := logreader.EntriesFromContent(buf.Bytes())
			require.NoError(t, err)
			require.Equal(t, 1, len(entries), "request log should have exactly 1 entry")

			err = tc.JSONMatcher.Matches(map[string]interface{}(entries[0]))
			assert.NoError(t, err, "Case %d: %s\n%v", i, tc.Name, err)
		})
	}
}

func GenerateRequest(requestHeaders map[string]string, extraQueryParams []string, extraHeaderParams map[string]string) *http.Request {
	uri, err := url.Parse(strings.Join(append([]string{"https://localhost:8443/some/path/here?fooQueryVarName=fooQueryVarVal&barQueryVarName=barQueryVarVal"}, extraQueryParams...), "&"))
	if err != nil {
		panic(err)
	}
	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Path:     uri.Path,
			RawQuery: uri.RawQuery,
			Scheme:   uri.Scheme,
			Host:     uri.Host,
		},
		Body:          nil,
		ContentLength: int64(0),
		Form:          uri.Query(),
		RequestURI:    uri.String(),
		Header:        http.Header{},
		Proto:         "HTTP/1.1",
	}
	for k, v := range requestHeaders {
		req.Header.Set(k, v)
	}
	for k, v := range extraHeaderParams {
		req.Header.Set(k, v)
	}
	return req
}
